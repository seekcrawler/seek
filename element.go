package kraken

import (
	"github.com/tebeka/selenium"
	"time"
)

func newElement(wd selenium.WebDriver, elem selenium.WebElement, extractor *Extractor) Element {
	return Element{wd: wd, elem: elem, extractor: extractor, scroller: &scroller{
		elem: "arguments[0]",
		wd:   wd,
		args: []any{elem},
		wait: extractor.Wait,
	}}
}

type Element struct {
	*scroller
	err       error
	wd        selenium.WebDriver
	elem      selenium.WebElement
	extractor *Extractor
}

func (p Element) Error() error {
	return p.err
}

func (p Element) FindElement(by By, selector string, timeout ...time.Duration) Element {
	if p.err != nil {
		return Element{
			err: p.err,
		}
	}
	return p.extractor.findElement(p.elem, by, selector, calcTimeDuration(timeout))
}

func (p Element) FindElements(by By, selector string, timeout ...time.Duration) Elements {
	if p.err != nil {
		return Elements{
			err: p.err,
		}
	}
	return p.extractor.findElements(p.elem, by, selector, calcTimeDuration(timeout))
}

func (p Element) Input(val string) error {
	if p.err != nil {
		return p.err
	}
	return p.elem.SendKeys(val)
}

func (p Element) Valid() (Element, error) {
	return p, p.err
}

func (p Element) Text() (test string, err error) {
	if p.err != nil {
		return "", p.err
	}
	return p.elem.Text()
}

func (p Element) Click() error {
	if p.err != nil {
		return p.err
	}
	return p.elem.Click()
}

func (p Element) MouseHover() (err error) {
	_, err = p.wd.ExecuteScript(prepareEventScript("mouseover"), []interface{}{p.elem})
	if err != nil {
		return
	}
	return
}

func (p Element) MouseOut() (err error) {
	_, err = p.wd.ExecuteScript(prepareEventScript("mouseout"), []interface{}{p.elem})
	if err != nil {
		return
	}
	return
}
