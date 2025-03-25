package ant_design

import (
	"fmt"
	"github.com/gozelle/fs"
	"github.com/seekcrawler/seek"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScroll(t *testing.T) {
	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	seek.DriverPath = dp

	router := seek.NewRouter(func(c *seek.Context) (err error) {

		elem, err := c.FindElement(seek.ByCSSSelector, `.ant-table-tbody-virtual-holder`).Valid()
		require.NoError(t, err)

		fmt.Println("start to scroll")
		fmt.Println(elem.ScrollHeight())

		err = elem.AutoWheelScrollBottom(seek.AutoWheelScrollBottomParams{
			PaddingHeight: 300,
			RowHeight:     300,
		})
		fmt.Println("stop")
		if err != nil {
			t.Log(err)
		}
		return
	})

	err = seek.Request("https://ant-design.antgroup.com/components/table-cn#table-demo-virtual-list",
		seek.WithChromeArgs([]string{
			//"--no-sandbox",
			//"--headless",    // 无头模式运行
			//"--disable-gpu", // 禁用 GPU
			//"--window-size=15360,3600",    // 设置窗口大小
			//"--force-device-scale-factor=2", // 设置缩放因子为 2 (确保高分辨率)
			//"--high-dpi-support=1.0",        // 避免在Linux环境下出现错误，可选
			//"--disable-dev-shm-usage",       // 避免在Linux环境下出现错误，可选
		}),
		seek.WithRouter(router),
	)

	require.NoError(t, err)
}
