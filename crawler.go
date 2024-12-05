package kraken

import (
	"fmt"
	"github.com/gozelle/logger"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"net/url"
	"os"
	"time"
)

var log = logger.NewLogger("crawler")

type RawUrl string

type Conf struct {
	debug                bool
	pageExtractorFactory PageExtractorFactory
	chromeArgs           []string
}

type PageExtractorFactory func(u url.URL) *Extractor

type Option func(c *Conf)

func WithChromeArgs(args []string) Option {
	return func(c *Conf) {
		c.chromeArgs = args
	}
}

func WithPageExtractorFactory(f PageExtractorFactory) Option {
	return func(c *Conf) {
		c.pageExtractorFactory = f
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
	}
}

type crawler struct {
	visitUrl  chan RawUrl
	done      chan error
	extractor *Extractor
}

func (c *crawler) close() {
	log.Debugf("close crawler")
	close(c.visitUrl)
	close(c.done)
}

func (c *crawler) prepareExtractor(conf *Conf, u url.URL) (extractor *Extractor, err error) {

	if conf.pageExtractorFactory == nil {
		err = fmt.Errorf("page extractor fatory is not initialized")
		return
	}
	extractor = conf.pageExtractorFactory(u)
	if extractor == nil {
		err = fmt.Errorf("new url: %s extractor is nil", u.String())
		return
	}
	return
}

func (c *crawler) defaultConf() *Conf {
	return &Conf{
		chromeArgs: []string{
			"--no-sandbox",
			"--headless",    // 无头模式运行
			"--disable-gpu", // 禁用 GPU
			//"--window-size=15360,3600",    // 设置窗口大小
			"--force-device-scale-factor=2", // 设置缩放因子为 2 (确保高分辨率)
			"--high-dpi-support=1.0",        // 避免在Linux环境下出现错误，可选
			"--disable-dev-shm-usage",       // 避免在Linux环境下出现错误，可选
		},
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

	go c.exec(conf, wd)
	go c.watchUrlChange(wd)
	go func() {

	}()

	err = wd.Get(rawUrl)
	if err != nil {
		err = fmt.Errorf("request url: %s error: %s", rawUrl, err)
		return
	}

	select {
	case err = <-c.done:
		if err != nil {
			log.Errorf("exec crawer error: %s", err)
		}
		return
	}
}

func (c *crawler) exec(conf *Conf, wd selenium.WebDriver) {
	for {
		select {
		case visitUrl, ok := <-c.visitUrl:
			if !ok {
				log.Infof("exist crawler exec, url: %s", visitUrl)
				return
			}
			u, err := url.Parse(string(visitUrl))
			if err != nil {
				c.done <- fmt.Errorf("parse url error: %w", err)
				return
			}
			c.extractor, err = c.prepareExtractor(conf, *u)
			if err != nil {
				c.done <- fmt.Errorf("prepare url: %s extractor error: %w", visitUrl, err)
				return
			}
			initExtractor(c.extractor, wd, *u)
			err = wd.Get(u.String())
			if err != nil {
				log.Errorf("wd get url: %s error: %s", u, err)
				c.done <- fmt.Errorf("visit url %s error: %w", visitUrl, err)
				return
			}

			var status ExtractorStatus
			status, err = c.extractor.Start()
			if err != nil {
				c.done <- err
				return
			}
			if status == ExtractorDone {
				log.Infof("visit url: %s done", visitUrl)
				c.done <- nil
				return
			}
		}
	}

}

func (c *crawler) watchUrlChange(wd selenium.WebDriver) {
	var currentUrl string
	for {
		select {
		default:
			newUrl, _ := wd.CurrentURL()
			if newUrl == "" {
				log.Debugf("close url watcher")
				return
			}
			if newUrl != currentUrl {
				if currentUrl == "" || currentUrl == "data:" {
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
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (c *crawler) stop() {
	c.done <- nil
}
