package kraken

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

func newElement(wd selenium.WebDriver, elem selenium.WebElement, extractor *Extractor) Element {
	return Element{wd: wd, elem: elem, extractor: extractor}
}

type Element struct {
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

func (p Element) ScrollTop() error {
	_, err := p.wd.ExecuteScript(`arguments[0].scrollTo({top:0,left:0,behavior:"smooth"});`, []interface{}{p.elem})
	return err
}

func (p Element) ScrollBottom() error {
	_, err := p.wd.ExecuteScript(`arguments[0].scrollTo({top:arguments[0].scrollHeight,left:0,behavior:"smooth"});`, []interface{}{p.elem})
	return err
}

func (p Element) WheelScrollY(y int64) error {
	_, err := p.wd.ExecuteScript(fmt.Sprintf(`arguments[0].dispatchEvent(new WheelEvent('wheel',{deltaY:%d,bubbles:true,cancelable:true}));`, y), []interface{}{p.elem})
	return err
}

func (p Element) ScrollHeight() (height int64, err error) {
	scrollHeight, err := p.wd.ExecuteScript("return arguments[0].scrollHeight;", []interface{}{p.elem})
	if err != nil {
		return
	}
	v, _ := scrollHeight.(float64)
	height = int64(v)
	return
}

func (p Element) WaitScrollHeightIncreased(previous int64, timeout ...time.Duration) error {
	_timeout := fixTimeDuration(calcTimeDuration(timeout))
	start := time.Now()
	for {
		height, err := p.ScrollHeight()
		if err != nil {
			return err
		}
		if height > previous {
			return nil
		}
		if time.Since(start) > _timeout {
			return TimoutErr
		}
		time.Sleep(CheckElementInterval)
	}
}

func (p Element) AutoScrollBottom(renderInterval time.Duration, handle func() error) (err error) {
	for {
		var h int64
		h, err = p.ScrollHeight()
		if err != nil {
			return
		}
		err = p.ScrollBottom()
		if err != nil {
			return
		}
		e := p.WaitScrollHeightIncreased(h, renderInterval)
		if e != nil {
			return
		}
		if handle != nil {
			err = handle()
			if err != nil {
				return
			}
		}
		p.extractor.Wait(renderInterval)
	}
}

func (p Element) AutoWheelScrollBottom(renderInterval time.Duration, rowHeight int64, handle func() error) (err error) {
	var y = int64(0)
	for {
		var h int64
		h, err = p.ScrollHeight()
		if err != nil {
			return
		}
		y += rowHeight
		if y > h {
			return
		}
		err = p.WheelScrollY(rowHeight)
		if err != nil {
			return
		}
		if handle != nil {
			err = handle()
			if err != nil {
				return
			}
		}
		if y > h {
			return
		}
		p.extractor.Wait(renderInterval)
	}
}
