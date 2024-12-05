package kraken

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

func TestURL(t *testing.T) {
	u, _ := url.Parse("http://example.com:80/user/tom/123?name=123")
	t.Log(u.Scheme)
	t.Log(u.Hostname())
	t.Log(u.Port())
}

func TestEngine(t *testing.T) {

	engine := NewRouter(func(c *Context) {})

	engine.Handle("/hello", func(c *Context) {
		fmt.Println("Hello world")
	})

	group := engine.Group("/hello")
	group.Handle("/sub", func(context *Context) {
		fmt.Println("sub handler")
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
		err := engine.handle(&Context{
			URL: *u,
		})
		require.NoError(t, err)

	}
	{
		u, _ := url.Parse("http://example.com/user/tom/123?name=123")
		err := engine.handle(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}

	{
		u, _ := url.Parse("http://example.com/hello/sub?name=123")
		err := engine.handle(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}

	time.Sleep(3 * time.Second)
}
