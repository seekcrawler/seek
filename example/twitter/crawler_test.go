package twitter_test

import (
	"fmt"
	"github.com/gozelle/fs"
	"github.com/seekcrawler/seek"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrawler(t *testing.T) {
	dp, err := fs.Lookupwd("./drivers/chromedriver_133_arm64")
	require.NoError(t, err)

	seek.DriverPath = dp

	router := seek.NewRouter(DefaultHandler)

	err = seek.Request("https://x.com/0xPrismatic/article/1872624976882512171",
		seek.WithChromeArgs([]string{}),
		seek.WithRouter(router),
	)

	require.NoError(t, err)
}

func DefaultHandler(c *seek.Context) (err error) {
	elem, err := c.FindElement(seek.ByCSSSelector, `div[data-testid="twitterArticleRichTextView"]`).Valid()
	if err != nil {
		c.Errorf("a is not found")
		return
	}
	fmt.Println(elem.Text())

	return
}
