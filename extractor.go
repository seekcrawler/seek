package kraken

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"net/url"
	"sync/atomic"
	"time"
)

var (
	DriverPath              = ""
	DefaultExtractorTimeout = 5 * time.Minute
	CheckElementInterval    = 100 * time.Millisecond
)

const minExtractorTimeout = 10 * time.Second

var (
	ElementNotFoundErr              = errors.New("element not found")
	OtherElementHasBeenProcessedErr = errors.New("other elem has been processed")
)

type ExtractorStatus int

const (
	ExtractorDone ExtractorStatus = iota
	ExtractorClose
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
	url     url.URL
	err     chan error
	timeout time.Duration
	done    chan ExtractorStatus
	wd      selenium.WebDriver
	runners []Runner

	hasDone atomic.Bool
}

func (p *Extractor) WebDriver() selenium.WebDriver {
	return p.wd
}

func (p *Extractor) Start() (status ExtractorStatus, err error) {
	defer func() {
		p.close()
	}()

	timeout := time.NewTimer(DefaultExtractorTimeout)

	log.Infof("run extractor, url: %s runners: %d", p.url.String(), len(p.runners))
	for _, v := range p.runners {
		fn := v
		go fn(p)
	}

	select {
	case status = <-p.done:
		return
	case <-timeout.C:
		err = fmt.Errorf("run extractor url: %s timeout", p.url.String())
		return
	}
}

func (p *Extractor) Done() {
	ok := p.hasDone.CompareAndSwap(false, true)
	if ok {
		p.done <- ExtractorDone
	}
}

func (p *Extractor) stop() {
	ok := p.hasDone.CompareAndSwap(false, true)
	if ok {
		p.done <- ExtractorClose
	}
}

func initExtractor(extractor *Extractor, wd selenium.WebDriver, url url.URL) {
	extractor.wd = wd
	extractor.url = url
	if extractor.done == nil {
		extractor.done = make(chan ExtractorStatus)
	}
	if extractor.err == nil {
		extractor.err = make(chan error)
	}
}

func (p *Extractor) close() {
	log.Debugf("close extractor, url: %s", p.url.String())
	close(p.done)
	close(p.err)
}

func (p *Extractor) Run(runner Runner) *Extractor {
	p.runners = append(p.runners, runner)
	return p
}
func (p *Extractor) FindElements(by By, selector string) Elements {
	return p.WithTimoutFindElements(by, selector, DefaultExtractorTimeout)
}

func (p *Extractor) WithTimoutFindElements(by By, selector string, timeout time.Duration) Elements {

	return p.findElements(nil, by, selector, timeout)
}

type iFindElements interface {
	FindElements(by, value string) ([]selenium.WebElement, error)
}

func (p *Extractor) findElements(parent iFindElements, by By, selector string, timeout time.Duration) Elements {
	if timeout <= 0 {
		timeout = minExtractorTimeout
	}
	start := time.Now()
	if parent == nil {
		parent = p.wd
	}
	for {
		if p.hasDone.Load() {
			log.Infof("cancel find elements, selector: %s %s", by, selector)
			return Elements{
				err: OtherElementHasBeenProcessedErr,
			}
		}
		results, err := parent.FindElements(string(by), selector)
		if err == nil {
			if len(results) > 0 {
				log.Debugf("find elements success, selector: %s %s, count: %d", by, selector, len(results))
				var elems []Element
				for _, elem := range results {
					elems = append(elems, Element{
						wd:   p.wd,
						elem: elem,
					})
				}
				if !p.hasDone.Load() {
					return Elements{
						wd:    p.wd,
						elems: elems,
					}
				} else {
					log.Infof("ignore to find elemets: selector: %s %s", by, selector)
					return Elements{
						err: ElementNotFoundErr,
					}
				}
			}
		}
		if time.Since(start) > timeout {
			log.Warnf("find elementos failed, selector: %s %s", by, selector)
			return Elements{
				err: ElementNotFoundErr,
			}
		}
		time.Sleep(CheckElementInterval)
	}
}

func (p *Extractor) FindElement(by By, selector string) Element {
	return p.WithTimoutFindElement(by, selector, DefaultExtractorTimeout)
}

func (p *Extractor) WithTimoutFindElement(by By, selector string, timeout time.Duration) Element {
	return p.findElement(nil, by, selector, timeout)
}

type iFindElement interface {
	FindElement(by, value string) (selenium.WebElement, error)
}

func (p *Extractor) findElement(parent iFindElement, by By, selector string, timeout time.Duration) Element {
	if timeout <= minExtractorTimeout {
		timeout = DefaultExtractorTimeout
	}
	start := time.Now()
	if parent == nil {
		parent = p.wd
	}
	for {
		if p.hasDone.Load() {
			log.Infof("cancel find element, selector: %s %s", by, selector)
			return Element{
				err: OtherElementHasBeenProcessedErr,
			}
		}
		elem, err := parent.FindElement(string(by), selector)
		if err == nil {
			var isDisplayed bool
			isDisplayed, err = elem.IsDisplayed()
			if err == nil && isDisplayed {
				log.Debugf("find element success, selector: %s %s", by, selector)
				if !p.hasDone.Load() {
					return Element{
						elem: elem,
						wd:   p.wd,
					}
				} else {
					log.Infof("ignore find elemet: selector: %s %s", by, selector)
					return Element{
						err: ElementNotFoundErr,
					}
				}
			}
		}
		if time.Since(start) > timeout {
			log.Warnf("find element failed, selector: %s %s", by, selector)
			return Element{
				err: ElementNotFoundErr,
			}
		}
		time.Sleep(CheckElementInterval)
	}
}

func (p *Extractor) ScrollBodyTop() error {
	_, err := p.wd.ExecuteScript(`window.scrollTo({top:0,left:0,behavior:"smooth"});`, nil)
	return err
}

func (p *Extractor) ScrollBodyBottom() error {
	_, err := p.wd.ExecuteScript(`window.scrollTo({top:document.body.scrollHeight,left:0,behavior:"smooth"});`, nil)
	return err
}

func (p *Extractor) WaitBodyScrollHeightChange(timeout ...time.Duration) (changed bool, err error) {
	return
}
