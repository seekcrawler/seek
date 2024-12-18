package kraken

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

type iBaseScroller interface {
	ScrollTop() error
	ScrollBottom() error
	ScrollHeight() (height int64, err error)
	WaitScrollHeightIncreased(previous int64, timeout ...time.Duration) error
	AutoScrollBottom(params AutoScrollBottomParams) (err error)
}

var _ iBaseScroller = (*baseScroller)(nil)

type baseScroller struct {
	wd               selenium.WebDriver
	args             []any
	wait             func(d ...time.Duration)
	scrollTopElem    func() string
	scrollBottomElem func() (string, string)
	scrollHeightElem func() string
	ctx              *Context
}

func (s baseScroller) ScrollTop() error {
	elem := s.scrollTopElem()
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`%s.scrollTo({top:0,left:0,behavior:"smooth"});`, elem), s.args)
	return err
}

func (s baseScroller) ScrollBottom() error {
	elem1, elem2 := s.scrollBottomElem()
	_, err := s.wd.ExecuteScript(fmt.Sprintf(`%s.scrollTo({top:%s.scrollHeight,left:0,behavior:"smooth"});`, elem1, elem2), s.args)
	return err
}

func (s baseScroller) ScrollHeight() (height int64, err error) {
	elem := s.scrollHeightElem()
	scrollHeight, err := s.wd.ExecuteScript(fmt.Sprintf("return %s.scrollHeight;", elem), s.args)
	if err != nil {
		return
	}
	v, _ := scrollHeight.(float64)
	height = int64(v)
	return
}

func (s baseScroller) WaitScrollHeightIncreased(previous int64, timeout ...time.Duration) error {
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
		s.wait(DefaultCheckElementInterval)
	}
}

type AutoScrollBottomParams struct {
	RenderInterval time.Duration
	WaitInterval   time.Duration
	Handler        func() error
}

func (s baseScroller) AutoScrollBottom(params AutoScrollBottomParams) (err error) {
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
