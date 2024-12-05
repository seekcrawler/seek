package kraken

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestEngine(t *testing.T) {

	engine := New()
	engine.Handle("/hello", func(c *Context) error {
		fmt.Println("Hello world")
		return nil
	})
	engine.Handle("/user/:name", func(c *Context) error {
		fmt.Println(c.Params)
		fmt.Println("Hello user")
		return nil
	})
	engine.Handle("/user/:name/:id", func(c *Context) error {
		fmt.Println(c.Params)
		fmt.Println("Hello user id")
		return nil
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

}
