package neteasy

import (
	"github.com/gozelle/fs"
	"github.com/gozelle/logger"
	"github.com/krakenspider/kraken"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var log = logger.NewLogger("neteasy")

func TestScroll(t *testing.T) {
	Handler := func(c *kraken.Context) error {
		return c.AutoScrollBottom(kraken.AutoScrollBottomParams{
			RenderInterval: 3 * time.Second,
			WaitInterval:   2 * time.Second,
			Handler:        nil,
		})
	}
	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	kraken.DriverPath = dp

	router := kraken.NewRouter(Handler)

	err = kraken.Request(
		"https://163.com",
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

func TestWheelScroll(t *testing.T) {
	Handler := func(c *kraken.Context) error {
		e := c.AutoWheelScrollBottom(kraken.AutoWheelScrollBottomParams{
			RenderInterval: 3 * time.Second,
		})
		if e != nil {
			log.Errorf("error: %v", e)
		}
		return e
	}
	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	kraken.DriverPath = dp

	router := kraken.NewRouter(Handler)

	err = kraken.Request(
		"https://163.com",
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
