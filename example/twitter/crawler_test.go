package twitter_test

import (
	"github.com/gozelle/fs"
	"github.com/krakenspider/kraken"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrawler(t *testing.T) {
	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	kraken.DriverPath = dp

	router := kraken.NewRouter(DefaultHandler)

	err = kraken.Request("https://x.com/elonmusk",
		kraken.WithChromeArgs([]string{}),
		kraken.WithRouter(router),
	)

	require.NoError(t, err)
}

func DefaultHandler(c *kraken.Context) (err error) {
	c.JustThink()
	elems, err := c.FindElements(kraken.ByCSSSelector, `a[href="/elonmusk"]`).Valid()
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
