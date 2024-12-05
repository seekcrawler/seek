package kraken

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestEngine(t *testing.T) {

	engine := New()
	engine.Handle("/hello", func(c *Context) {
		fmt.Println("Hello world")
	})

	group := engine.Group("/shop", func(c *Context) {
		fmt.Println("this form shop middleware before")
		c.Next()
		fmt.Println("this form shop middleware after")
	})

	group.Handle("/list/:name", func(c *Context) {
		fmt.Println("this shop list", c.Params.ByName("name"))
	})

	engine.Handle("/user/:name", func(c *Context) {
		fmt.Println(c.Params)
		fmt.Println("Hello user")
	})
	engine.Handle("/user/:name/:id", func(c *Context) {
		fmt.Println(c.Params)
		fmt.Println("Hello user id")
	})

	{
		u, _ := url.Parse("http://example.com/hello?name=123")
		err := engine.handleHTTPRequest(&Context{
			URL: *u,
		})
		require.NoError(t, err)

	}
	{
		u, _ := url.Parse("http://example.com/user/tom/123?name=123")
		err := engine.handleHTTPRequest(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}
	{
		u, _ := url.Parse("http://example.com/shop/list/jack?name=123")
		err := engine.handleHTTPRequest(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}
}
