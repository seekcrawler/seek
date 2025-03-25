package seek

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

	engine.Handle("/i/flow/login", func(c *Context) {
		fmt.Println("/i/flow/login")
	})

	engine.Handle("/:username", func(c *Context) {
		fmt.Println("/:username", c.Params)
	})
	engine.Handle("/:username/following", func(c *Context) {
		fmt.Println("/:username/following", c.Params)
	})

	{
		u, _ := url.Parse("http://example.com/i/flow/login?name=123")
		err := engine.handle(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}
	{
		u, _ := url.Parse("http://example.com/elonmusk?name=123")
		err := engine.handle(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}
	{
		u, _ := url.Parse("http://example.com/elonmusk/following?name=123")
		err := engine.handle(&Context{
			URL: *u,
		})
		require.NoError(t, err)
	}
	time.Sleep(3 * time.Second)
}
