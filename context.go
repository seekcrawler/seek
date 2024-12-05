package kraken

import (
	"net/url"
)

type Context struct {
	URL       url.URL
	Params    Params
	Extractor *Extractor
	handlers  HandlersChain
	index     int8
	check     func() bool
}

func (c *Context) Done() {
	c.Extractor.done()
}

func (c *Context) Check(check func() bool) {
	c.check = check
}

func (c *Context) reset() {
	c.Params = c.Params[:0]
	c.handlers = nil
	c.index = -1
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.handlers[c.index] != nil {
			index := c.index
			go func() {
				exec := true
				if c.check != nil {
					exec = c.check()
				}
				if !exec {
					c.Done()
				} else {
					c.handlers[index](c)
				}
			}()
		}
		c.index++
	}
}

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc
