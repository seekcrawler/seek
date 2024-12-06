package kraken

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

type iScroller interface {
	ScrollTop() error
	ScrollBottom() error
	WheelScrollX(x int64) error
	WheelScrollY(y int64) error
	ScrollHeight() (height int64, err error)
	WaitScrollHeightIncreased(previous int64, timeout ...time.Duration) error
	AutoScrollBottom(params AutoScrollBottomParams) (err error)
	AutoWheelScrollBottom(params AutoWheelScrollBottomParams) (err error)
}

var _ iScroller = (*scroller)(nil)

type scroller struct {
	elem string
	wd   selenium.WebDriver
	args []any
	wait func(d ...time.Duration)
}

func (s scroller) ScrollTop() error {
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`%s.scrollTo({top:0,left:0,behavior:"smooth"});`, s.elem), s.args)
	return err
}

func (s scroller) ScrollBottom() error {
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`%s.scrollTo({top:%s.scrollHeight,left:0,behavior:"smooth"});`, s.elem, s.elem), s.args)
	return err
}
func (s scroller) WheelScrollX(x int64) error {
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`%s.dispatchEvent(new WheelEvent('wheel',{deltaX:%d,bubbles:true,cancelable:true}));`, s.elem, x), s.args)
	return err
}
func (s scroller) WheelScrollY(y int64) error {
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`%s.dispatchEvent(new WheelEvent('wheel',{deltaY:%d,bubbles:true,cancelable:true}));`, s.elem, y), s.args)
	return err
}

func (s scroller) ScrollHeight() (height int64, err error) {
	scrollHeight, err := s.wd.ExecuteScript(fmt.Sprintf("return %s.scrollHeight;", s.elem), s.args)
	if err != nil {
		return
	}
	v, _ := scrollHeight.(float64)
	height = int64(v)
	return
}

func (s scroller) WaitScrollHeightIncreased(previous int64, timeout ...time.Duration) error {
	_timeout := fixTimeDuration(calcTimeDuration(timeout))
	start := time.Now()
	for {
		height, err := s.ScrollHeight()
		if err != nil {
			return err
		}
		if height > previous {
			return nil
		}
		if time.Since(start) > _timeout {
			return TimoutErr
		}
		s.wait(CheckElementInterval)
	}
}

type AutoScrollBottomParams struct {
	RenderInterval time.Duration
	WaitInterval   time.Duration
	Handler        func() error
}

func (s scroller) AutoScrollBottom(params AutoScrollBottomParams) (err error) {
	if params.WaitInterval == 0 {
		params.WaitInterval = 3 * time.Second
	}
	for {
		var h int64
		h, err = s.ScrollHeight()
		if err != nil {
			return
		}
		err = s.ScrollBottom()
		if err != nil {
			return
		}
		e := s.WaitScrollHeightIncreased(h, params.WaitInterval)
		if e != nil {
			return
		}
		if params.Handler != nil {
			err = params.Handler()
			if err != nil {
				return
			}
		}
		s.wait(params.RenderInterval)
	}
}

type AutoWheelScrollBottomParams struct {
	RenderInterval time.Duration
	RowHeight      int64
	PaddingHeight  int64
	Handle         func() error
}

func (s scroller) AutoWheelScrollBottom(params AutoWheelScrollBottomParams) (err error) {
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
