package kraken

import (
	"context"
	"errors"
	"fmt"
	"github.com/gozelle/logger"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

var log = logger.NewLogger("crawler")

type RawUrl string

type Conf struct {
	ctx         context.Context
	debug       bool
	chromeArgs  []string
	router      *Router
	dataHandler func(dataC chan any)
	timeout     time.Duration
	urlMode     URLMode
}

type Option func(c *Conf)

func WithContext(ctx context.Context) Option {
	return func(c *Conf) {
		c.ctx = ctx
	}
}

func WithChromeArgs(args []string) Option {
	return func(c *Conf) {
		c.chromeArgs = args
	}
}

func WithRouter(router *Router) Option {
	return func(c *Conf) {
		c.router = router
	}
}

func WithDataHandler(handler func(dataC chan any)) Option {
	return func(c *Conf) {
		c.dataHandler = handler
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Conf) {
		c.timeout = timeout
	}
}

func WithUrlMode(mode URLMode) Option {
	return func(c *Conf) {
		c.urlMode = mode
	}
}

func Request(rawUrl string, options ...Option) error {
	c := newCrawler()
	defer func() {
		c.close()
	}()
	return c.Run(rawUrl, options...)
}

func newCrawler() *crawler {
	return &crawler{
		visitUrl: make(chan RawUrl),
		done:     make(chan error),
		data:     make(chan any),
	}
}

type crawler struct {
	visitUrl  chan RawUrl
	done      chan error
	extractor *Extractor
	data      chan any
	hasDone   atomic.Bool
}

func (c *crawler) close() {
	if c.hasDone.CompareAndSwap(false, true) {
		close(c.visitUrl)
		close(c.done)
		close(c.data)
	}
	log.Infof("close crawler")
}

func (c *crawler) sendDone(err error) {
	if !c.hasDone.Load() {
		c.done <- err
	}
}

func (c *crawler) defaultConf() *Conf {
	return &Conf{
		ctx:   context.Background(),
		debug: false,
		chromeArgs: []string{
			"--no-sandbox",
			"--headless",    // 无头模式运行
			"--disable-gpu", // 禁用 GPU
			//"--window-size=15360,3600",    // 设置窗口大小
			"--force-device-scale-factor=2", // 设置缩放因子为 2 (确保高分辨率)
			"--high-dpi-support=1.0",        // 避免在Linux环境下出现错误，可选
			"--disable-dev-shm-usage",       // 避免在Linux环境下出现错误，可选
		},
		router:      nil,
		dataHandler: nil,
		timeout:     DefaultExtractorTimeout,
	}
}

func (c *crawler) Run(rawUrl string, options ...Option) (err error) {
	conf := c.defaultConf()
	for _, option := range options {
		option(conf)
	}

	var opts []selenium.ServiceOption
	if conf.debug {
		opts = append(opts, selenium.Output(os.Stdout))
		selenium.SetDebug(true)
	}

	port, err := getActivePort()
	if err != nil {
		return
	}

	service, err := selenium.NewChromeDriverService(DriverPath, port, opts...)
	if err != nil {
		err = fmt.Errorf("failed to start ChromeDriverService: %w", err)
		return
	}
	defer func() {
		_ = service.Stop()
	}()

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: conf.chromeArgs,
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return
	}
	defer func() {
		_ = wd.Quit()
	}()

	if conf.dataHandler != nil {
		go func() {
			conf.dataHandler(c.data)
		}()
	}

	go c.exec(conf.ctx, conf, wd)
	go c.watchUrlChange(wd)
	err = wd.Get(rawUrl)
	if err != nil {
		err = fmt.Errorf("request url: %s error: %s", rawUrl, err)
		return
	}

	var ok bool
	select {
	case err, ok = <-c.done:
		if !ok {
			return
		}
		if err != nil {
			log.Errorf("exec crawer error: %s", err)
		}
		return
	}
}

func (c *crawler) exec(ctx context.Context, conf *Conf, wd selenium.WebDriver) {
	for {
		select {
		case visitUrl, ok := <-c.visitUrl:
			if !ok {
				log.Infof("exist crawler exector")
				return
			}
			u, err := url.Parse(string(visitUrl))
			if err != nil {
				c.sendDone(fmt.Errorf("parse url error: %w", err))
				return
			}
			if c.extractor != nil && c.extractor.ctx != nil && c.extractor.ctx.Context != nil {
				ctx = c.extractor.ctx.Context
			}
			c.extractor = NewExtractor()
			initExtractor(c, wd, *u, conf.timeout)
			err = wd.Get(u.String())
			if err != nil {
				log.Errorf("wd get url: %s error: %s", u, err)
				c.sendDone(fmt.Errorf("visit url %s error: %w", visitUrl, err))
				return
			}

			kCtx := &Context{
				URL:     *u,
				Context: ctx,
			}
			if conf.router == nil {
				c.sendDone(fmt.Errorf("router is nil"))
				return
			}
			err = conf.router.prepareContext(kCtx, c.extractor)
			if err != nil && !errors.Is(err, handlerNotFoundErr) {
				c.sendDone(fmt.Errorf("prepare router error: %w", err))
				return
			}
			if len(kCtx.handlers) == 0 {
				if conf.router.defaultHandler == nil {
					c.sendDone(fmt.Errorf("no default handler"))
					return
				} else {
					kCtx.Extractor = c.extractor
					kCtx.handlers = append(kCtx.handlers, conf.router.defaultHandler)
				}
			}
			err = c.extractor.Start(kCtx)
			if err != nil {
				c.sendDone(err)
				return
			}
		}
	}

}

func (c *crawler) currentUrl(wd selenium.WebDriver) string {
	defer func() {
		recover()
	}()
	v, _ := wd.CurrentURL()
	return v
}

func (c *crawler) watchUrlChange(wd selenium.WebDriver) {
	var currentUrl string
	pollInterval := 5 * time.Millisecond
	const emptyUrl = "data:,"
	for {
		select {
		default:
			newUrl := c.currentUrl(wd)
			if newUrl == "" {
				//log.Debugf("close url watcher")
				c.sendDone(fmt.Errorf("browser has closed"))
				return
			}
			if newUrl == emptyUrl {
				time.Sleep(pollInterval)
				continue
			}
			if newUrl != currentUrl {
				if currentUrl == "" || currentUrl == emptyUrl {
					log.Infof("load url: %s", newUrl)
				} else {
					log.Infof("url changeed, previous: %s, now: %s", currentUrl, newUrl)
				}
				if c.extractor != nil {
					c.extractor.stop()
				}
				c.visitUrl <- RawUrl(newUrl)
				currentUrl = newUrl
			}
			time.Sleep(pollInterval)
		}
	}
}

func (c *crawler) stop() {
	c.sendDone(nil)
}
