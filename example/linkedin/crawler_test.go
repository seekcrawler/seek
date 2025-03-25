package twitter_test

import (
	"github.com/gozelle/fs"
	"github.com/seekcrawler/seek"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrawler(t *testing.T) {
	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	seek.DriverPath = dp

	router := seek.NewRouter(DefaultHandler)

	err = seek.Request("https://x.com/elonmusk",
		seek.WithChromeArgs([]string{}),
		seek.WithRouter(router),
	)

	require.NoError(t, err)
}

func DefaultHandler(c *seek.Context) (err error) {
	c.JustThink()
	elems, err := c.FindElements(seek.ByCSSSelector, `a[href="/elonmusk"]`).Valid()
	if err != nil {
		c.Errorf("a is not found")
		return
	}
	c.Debugf("elems: %d", elems.Len())
	for _, v := range elems.Elements() {
		v.MouseOver()
		c.JustWait()
		v.MouseOut()
	}
	return
}
