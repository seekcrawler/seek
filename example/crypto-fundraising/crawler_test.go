package crypto_fundraising

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

	err = seek.Request("https://crypto-fundraising.info/deal-flow/page/2/",
		seek.WithChromeArgs([]string{}),
		seek.WithRouter(router),
	)

	require.NoError(t, err)
}

func DefaultHandler(c *seek.Context) (err error) {
	elems, err := c.FindElements(seek.ByCSSSelector, `div.dealflow-table > div`).Valid()
	if err != nil {
		c.Errorf("a is not found")
		return
	}

	l := -1
	for _, elem := range elems.Elements() {
		l++
		t, _ := elem.Text()
		fmt.Printf("\n\n=======\n")
		fmt.Println(l, t)
	}

	return
}
