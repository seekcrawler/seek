package kraken

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

func newElement(wd selenium.WebDriver, elem selenium.WebElement, extractor *Extractor) Element {
	return Element{wd: wd, elem: elem, extractor: extractor, baseScroller: &baseScroller{
		wd:   wd,
		args: []any{elem},
		wait: extractor.Wait,
		scrollTopElem: func() string {
			return "arguments[0]"
		},
		scrollBottomElem: func() (string, string) {
			return "arguments[0]", "arguments[0]"
		},
		scrollHeightElem: func() string {
			return "arguments[0]"
		},
	}}
}

type Element struct {
	*baseScroller
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
	return p.extractor.findElement(p.elem, by, selector, nil, calcTimeDuration(timeout))
}

func (p Element) FindElementWithPolling(by By, selector string, poll func(), timeout ...time.Duration) Element {
	if p.err != nil {
		return Element{
			err: p.err,
		}
	}
	return p.extractor.findElement(p.elem, by, selector, poll, calcTimeDuration(timeout))
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

func (p Element) Screenshot(scroll bool) ([]byte, error) {
	return p.elem.Screenshot(scroll)
}

func (p Element) Click() error {
	if p.err != nil {
		return p.err
	}
	return p.elem.Click()
}

func (p Element) Href() (string, error) {
	return p.elem.GetAttribute("href")
}

func (p Element) MouseOver() error {
	if p.err != nil {
		return p.err
	}
	//offset, err := p.elem.Location()
	//if err != nil {
	//	return
	//}
	//p.wd.StorePointerActions(shortuuid.New(),
	//	selenium.MousePointer,
	//	selenium.PointerMoveAction(0, *offset, selenium.FromViewport),
	//)
	//err = p.wd.PerformActions()
	//if err != nil {
	//	return
	//}
	//_ = p.wd.ReleaseActions()
	_, err := p.wd.ExecuteScript(prepareEventScript("mouseover"), []interface{}{p.elem})
	if err != nil {
		return err
	}
	return nil
}

func (p Element) MouseOut() error {
	if p.err != nil {
		return p.err
	}
	//offset := selenium.Point{
	//	X: 0,
	//	Y: 0,
	//}
	//p.wd.StorePointerActions(shortuuid.New(),
	//	selenium.MousePointer,
	//	selenium.PointerMoveAction(0, offset, selenium.FromViewport),
	//)
	//err = p.wd.PerformActions()
	//if err != nil {
	//	return
	//}
	//_ = p.wd.ReleaseActions()
	_, err := p.wd.ExecuteScript(prepareEventScript("mouseout"), []interface{}{p.elem})
	if err != nil {
		return err
	}
	return nil
}

func (s baseScroller) WheelScrollX(x int64) error {
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`arguments[0].dispatchEvent(new WheelEvent('wheel',{deltaX:%d,bubbles:true,cancelable:true}));`, x), s.args)
	return err
}

func (s baseScroller) WheelScrollY(y int64) error {
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`arguments[0].dispatchEvent(new WheelEvent('wheel',{deltaY:%d,bubbles:true,cancelable:true}));`, y), s.args)
	return err
}

func (s baseScroller) AutoWheelScrollBottom(params AutoWheelScrollBottomParams) (err error) {
	if params.RowHeight == 0 {
		params.RowHeight = 14
	}
	var y = int64(0)
	for {
		var h int64
		h, err = s.ScrollHeight()
		if err != nil {
			return
		}
		h += params.PaddingHeight
		y += params.RowHeight
		err = s.WheelScrollY(params.RowHeight)
		if err != nil {
			return
		}
		if params.Handle != nil {
			err = params.Handle()
			if err != nil {
				return
			}
		}
		if y > h {
			return
		}
		s.wait(params.RenderInterval)
	}
}
