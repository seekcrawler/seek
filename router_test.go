package kraken

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestEngine(t *testing.T) {

	engine := NewRouter()
	engine.Handle("/hello", func(c *Context) {
		fmt.Println("Hello world")
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
}
