package kraken

import (
	"net/url"
)

type Context struct {
	*Extractor
	URL      url.URL
	Params   Params
	handlers HandlersChain
	index    int8
	abort    func() bool
}

func (c *Context) Abort(fn func() bool) {
	c.abort = fn
}

func (c *Context) Done() {
	c.Extractor.done()
}

func (c *Context) reset() {
	c.Params = c.Params[:0]
	c.handlers = nil
	c.index = -1
}

func (c *Context) stop() {
	if c.Extractor != nil {
		c.Extractor.stop()
	}
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.handlers[c.index] != nil {
			index := c.index
			abort := false
			if c.abort != nil {
				abort = c.abort()
			}
			if abort {
				c.stop()
			} else {
				c.handlers[index](c)
			}
		}
		c.index++
	}
	c.stop()
}

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc
