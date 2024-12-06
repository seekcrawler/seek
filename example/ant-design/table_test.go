package ant_design

import (
	"fmt"
	"github.com/gozelle/fs"
	"github.com/krakenspider/kraken"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScroll(t *testing.T) {
	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	kraken.DriverPath = dp

	router := kraken.NewRouter(func(c *kraken.Context) {

		elem, err := c.FindElement(kraken.ByCSSSelector, `.ant-table-tbody-virtual-holder`).Valid()
		require.NoError(t, err)

		fmt.Println("start to scroll")
		fmt.Println(elem.ScrollHeight())

		err = elem.AutoWheelScrollBottom(0, 300, func() error {
			return nil
		})
		fmt.Println("stop")
		if err != nil {
			t.Log(err)
		}
	})

	err = kraken.Request("https://ant-design.antgroup.com/components/table-cn#table-demo-virtual-list",
		kraken.WithChromeArgs([]string{
			//"--no-sandbox",
			//"--headless",    // 无头模式运行
			//"--disable-gpu", // 禁用 GPU
			//"--window-size=15360,3600",    // 设置窗口大小
			//"--force-device-scale-factor=2", // 设置缩放因子为 2 (确保高分辨率)
			//"--high-dpi-support=1.0",        // 避免在Linux环境下出现错误，可选
			//"--disable-dev-shm-usage",       // 避免在Linux环境下出现错误，可选
		}),
		kraken.WithRouter(router),
	)

	require.NoError(t, err)
}
