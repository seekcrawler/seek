package kraken

import "net/url"

type Context struct {
	URL      url.URL
	Params   Params
	handlers HandlersChain
	index    int8
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
			c.handlers[c.index](c)
		}
		c.index++
	}
}

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc
