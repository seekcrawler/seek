package kraken

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

var (
	DriverPath              = ""
	DefaultExtractorTimeout = 5 * time.Minute
	CheckElementInterval    = 100 * time.Millisecond
)

const minExtractorTimeout = 0

var (
	ElementNotFoundErr  = errors.New("element not found")
	ExtractorStoppedErr = errors.New("extractor has stopped")
	TimoutErr           = errors.New("timeout")
)

var (
	handlerNotFoundErr = errors.New("handler not found")
)

type extractorStatus int

const (
	extractorStop extractorStatus = iota
)

type By string

const (
	ByID              By = "id"
	ByXPATH           By = "xpath"
	ByLinkText        By = "link text"
	ByPartialLinkText By = "partial link text"
	ByName            By = "name"
	ByTagName         By = "tag name"
	ByClassName       By = "class name"
	ByCSSSelector     By = "css selector"
)

type Runner func(ctx *Extractor)

func NewExtractor() *Extractor {
	return &Extractor{}
}

type Extractor struct {
	*scroller
	url      url.URL
	wd       selenium.WebDriver
	hasEnd   atomic.Bool
	hasClose atomic.Bool
	errC     chan error
	stopC    chan struct{}
	crawler  *crawler
}

func (p *Extractor) Wait(t ...time.Duration) {
	d := calcTimeDuration(t)
	if d > 0 {
		time.Sleep(fixTimeDuration(d))
	}
}

func (p *Extractor) CurrentURL() *url.URL {
	if p.wd != nil {
		r, _ := p.wd.CurrentURL()
		u, err := url.Parse(r)
		if err == nil {
			return u
		}
	}
	return &url.URL{}
}

func (p *Extractor) Start(ctx *Context) (err error) {
	defer func() {
		p.close()
	}()

	timeout := time.NewTimer(DefaultExtractorTimeout)

	log.Infof("run extractor, url: %s handlers: %d", p.url.String(), len(ctx.handlers))

	go ctx.Next()

	select {
	case <-p.stopC:
		return
	case <-timeout.C:
		err = fmt.Errorf("run extractor url: %s timeout", p.url.String())
		return
	}
}

func (p *Extractor) done() {
	p.crawler.sendDone(nil)
}

func (p *Extractor) stop() {
	ok := p.hasEnd.CompareAndSwap(false, true)
	if ok && !p.hasClose.Load() {
		p.stopC <- struct{}{}
	}
}

func initExtractor(c *crawler, wd selenium.WebDriver, url url.URL) {
	extractor := c.extractor
	extractor.crawler = c
	extractor.wd = wd
	extractor.url = url
	if extractor.stopC == nil {
		extractor.stopC = make(chan struct{})
	}
	if extractor.errC == nil {
		extractor.errC = make(chan error)
	}
	if extractor.scroller == nil {
		extractor.scroller = &scroller{
			elem: "document.body",
			wd:   wd,
			args: nil,
			wait: extractor.Wait,
		}
	}
}

func (p *Extractor) close() {
	log.Debugf("close extractor, url: %s", p.url.String())
	if p.hasClose.CompareAndSwap(false, true) {
		close(p.stopC)
		close(p.errC)
	}
}

func (p *Extractor) FindElements(by By, selector string, timeout ...time.Duration) Elements {
	return p.findElements(nil, by, selector, calcTimeDuration(timeout))
}

type iFindElements interface {
	FindElements(by, value string) ([]selenium.WebElement, error)
}

func (p *Extractor) findElements(parent iFindElements, by By, selector string, timeout time.Duration) Elements {
	timeout = fixTimeDuration(timeout)
	start := time.Now()
	if parent == nil {
		parent = p.wd
	}
	for {
		if p.hasEnd.Load() {
			log.Infof("cancel find elements, by: %s selector: %s", by, selector)
			return Elements{
				err: ExtractorStoppedErr,
			}
		}
		results, err := parent.FindElements(string(by), selector)
		if err == nil {
			if len(results) > 0 {
				log.Debugf("find elements success, by: %s selector: %s, count: %d", by, selector, len(results))
				var elems []Element
				for _, elem := range results {
					elems = append(elems, newElement(p.wd, elem, p))
				}
				if !p.hasEnd.Load() {
					return Elements{
						wd:    p.wd,
						elems: elems,
					}
				} else {
					log.Infof("ignore to find elemets: by: %s selector: %s", by, selector)
					return Elements{
						err: ElementNotFoundErr,
					}
				}
			}
		}
		if time.Since(start) > timeout {
			log.Warnf("find elementos failed, by: %s selector: %s", by, selector)
			return Elements{
				err: ElementNotFoundErr,
			}
		}
		time.Sleep(CheckElementInterval)
	}
}

func (p *Extractor) FindElement(by By, selector string, timeout ...time.Duration) Element {
	return p.findElement(nil, by, selector, calcTimeDuration(timeout))

}

type iFindElement interface {
	FindElement(by, value string) (selenium.WebElement, error)
}

func (p *Extractor) findElement(parent iFindElement, by By, selector string, timeout time.Duration) Element {
	timeout = fixTimeDuration(timeout)
	start := time.Now()
	if parent == nil {
		parent = p.wd
	}
	for {
		if p.hasEnd.Load() {
			log.Infof("cancel find element, by: %s selector: %s", by, selector)
			return Element{
				err: ExtractorStoppedErr,
			}
		}
		elem, err := parent.FindElement(string(by), selector)
		if err == nil {
			var isDisplayed bool
			isDisplayed, err = elem.IsDisplayed()
			if err == nil && isDisplayed {
				log.Debugf("find element success, by: %s selector: %s", by, selector)
				if !p.hasEnd.Load() {
					return newElement(p.wd, elem, p)
				} else {
					log.Infof("ignore find elemet: by: %s selector: %s", by, selector)
					return Element{
						err: ElementNotFoundErr,
					}
				}
			}
		}
		if time.Since(start) > timeout {
			log.Warnf("find element failed, by: %s selector: %s", by, selector)
			return Element{
				err: ElementNotFoundErr,
			}
		}
		time.Sleep(CheckElementInterval)
	}
}

func (p *Extractor) GetCookies() ([]selenium.Cookie, error) {
	return p.wd.GetCookies()
}

func (p *Extractor) GetCookie(name string, timeout ...time.Duration) (cookie selenium.Cookie, err error) {
	_timeout := fixTimeDuration(calcTimeDuration(timeout))
	start := time.Now()
	for {
		cookie, err = p.wd.GetCookie(name)
		if err == nil {
			return
		}
		if time.Since(start) > _timeout {
			err = TimoutErr
			return
		}
		time.Sleep(CheckElementInterval)
	}
}

func (p *Extractor) AddCoolie(cookie selenium.Cookie) error {
	return p.wd.AddCookie(&cookie)
}

func (p *Extractor) AddCookies(cookies []selenium.Cookie) error {
	for _, v := range cookies {
		err := p.wd.AddCookie(&v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Extractor) ParseCookiesJson(content []byte) (cookies []selenium.Cookie, err error) {
	err = json.Unmarshal(content, &cookies)
	return
}

func (p *Extractor) Refresh() error {
	return p.wd.Refresh()
}

func (p *Extractor) Redirect(path string) error {
	if strings.HasPrefix(path, "http") {
		return p.wd.Get(path)
	}
	port := p.url.Port()
	if port != "" {
		port = fmt.Sprintf(":%s", port)
	}
	return p.wd.Get(fmt.Sprintf("%s://%s%s%s", p.url.Scheme, p.url.Hostname(), port, path))
}
