package kraken_test

import (
	"errors"
	"fmt"
	"github.com/gozelle/fs"
	"github.com/gozelle/logger"
	"github.com/krakenspider/kraken"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var log = logger.NewLogger("test")

//func GenerateGoogleCode(secret string) (code string, err error) {
//	now := time.Now()
//	code, err = totp.GenerateCode(secret, now)
//	if err != nil {
//		return
//	}
//	return
//}
//
//func getCooliesFilePath() string {
//	binPath, err := fs.Lookupwd("/bin")
//	if err != nil {
//		panic("bin dir not found")
//	}
//	return filepath.Join(binPath, "/x-cookies.json")
//}
//
//func NewUserPageExtractor(u url.URL, options ...kraken.Option) *kraken.Extractor {
//
//	return kraken.NewExtractor().
//		Run(func(extractor *kraken.Extractor) {
//			defer func() {
//				extractor.Done()
//			}()
//			image, err := extractor.FindElement(kraken.ByCSSSelector, `a[href="#summary"]`).Text()
//			if err != nil {
//				return
//			}
//			fmt.Println("summary", image)
//
//		}).
//		Run(func(extractor *kraken.Extractor) {
//			defer func() {
//				extractor.Done()
//			}()
//			text, err := extractor.FindElement(kraken.ByCSSSelector, `a[href="#specification2"]`).Text()
//			if err != nil {
//				return
//			}
//			fmt.Println("specification", text)
//		})
//}
//
//func NewTwitterExtractor(u url.URL, options ...kraken.Option) *kraken.Extractor {
//	if u.Path == "/i/flow/login" {
//		// 找到表单，登录，输入 2 次认证
//		return kraken.NewExtractor().Run(func(extractor *kraken.Extractor) {
//			defer func() {
//				extractor.Done()
//			}()
//			err := extractor.FindElement(kraken.ByCSSSelector, `input[autocomplete="username"]`).
//				Input("Jason0751489313")
//			if err != nil {
//				log.Errorf("input username error: %s", err)
//				return
//			}
//
//			elem, err := extractor.WebDriver().FindElement(string(kraken.ByCSSSelector), `div[role="dialog"]`)
//			if err != nil {
//				log.Errorf("find dialog wrapper error: %s", err)
//				return
//			}
//
//			elems, err := elem.FindElements(string(kraken.ByTagName), `button`)
//			if err != nil {
//				log.Errorf("find diaglog buttons error: %s", err)
//				return
//			}
//			for _, btn := range elems {
//				text, _ := btn.Text()
//				if text == "下一步" {
//					_ = btn.Click()
//					break
//				}
//			}
//
//			time.Sleep(10 * time.Second)
//
//			err = extractor.FindElement(kraken.ByCSSSelector, `input[autocomplete="current-password"]`).Input("CXPYMQOMhA")
//			if err != nil {
//				log.Errorf("input password error: %s", err)
//				return
//			}
//
//			elem, err = extractor.WebDriver().FindElement(string(kraken.ByCSSSelector), `div[role="dialog"]`)
//			if err != nil {
//				log.Errorf("find dialog wrapper error: %s", err)
//				return
//			}
//
//			elems, err = elem.FindElements(string(kraken.ByTagName), `button`)
//			if err != nil {
//				log.Errorf("find diaglog buttons error: %s", err)
//				return
//			}
//			for _, btn := range elems {
//				text, _ := btn.Text()
//				if text == "登录" {
//					_ = btn.Click()
//					break
//				}
//			}
//
//			time.Sleep(15 * time.Second)
//
//			code, err := GenerateGoogleCode("X3UN6JCI5ULH7YUQ")
//			if err != nil {
//				log.Errorf("generate code error: %s", err)
//				return
//			}
//			err = extractor.FindElement(kraken.ByCSSSelector, `input[inputmode="numeric"]`).Input(code)
//			if err != nil {
//				log.Errorf("input code error: %s", err)
//				return
//			}
//
//			elem, err = extractor.WebDriver().FindElement(string(kraken.ByCSSSelector), `div[role="dialog"]`)
//			if err != nil {
//				log.Errorf("find dialog wrapper error: %s", err)
//				return
//			}
//
//			elems, err = elem.FindElements(string(kraken.ByTagName), `button`)
//			if err != nil {
//				log.Errorf("find diaglog buttons error: %s", err)
//				return
//			}
//			for _, btn := range elems {
//				text, _ := btn.Text()
//				if text == "下一步" {
//					_ = btn.Click()
//					break
//				}
//			}
//
//			time.Sleep(5 * time.Second)
//
//			cookies, err := extractor.WebDriver().GetCookies()
//			if err != nil {
//				log.Errorf("get cookies error: %s", err)
//				return
//			}
//
//			cookieFile := getCooliesFilePath()
//			d, _ := json.Marshal(cookies)
//			_ = fs.Write(cookieFile, d)
//
//			time.Sleep(3 * time.Minute)
//		})
//	} else {
//		return kraken.NewExtractor().
//			Run(func(extractor *kraken.Extractor) {
//
//				cookieFile := getCooliesFilePath()
//				login := true
//				if fs.Exists(cookieFile) {
//					data, err := fs.Read(cookieFile)
//					if err != nil {
//						log.Errorf("read cookie file error: %s", err)
//						return
//					}
//					var cookies []*selenium.Cookie
//					err = json.Unmarshal(data, &cookies)
//					if err != nil {
//						log.Errorf("unmarshal cookie file error: %s", err)
//						return
//					}
//					active := false
//					for _, cookie := range cookies {
//						if cookie.Name == "auth_token" {
//							if cookie.Expiry > uint(time.Now().Unix()+86400) {
//								active = true
//								break
//							}
//						}
//					}
//					if active {
//						for _, v := range cookies {
//							err = extractor.WebDriver().AddCookie(v)
//							if err != nil {
//								log.Errorf("add cookie error: %s", err)
//								return
//							}
//						}
//						login = false
//						err = extractor.WebDriver().Refresh()
//						if err != nil {
//							log.Errorf("refresh page error: %s", err)
//							return
//						}
//					}
//				}
//
//				if login {
//					err := extractor.FindElement(kraken.ByCSSSelector, `a[href="/login"]`).Click()
//					if err != nil {
//						log.Errorf("click login button error: %s", err)
//						return
//					}
//				}
//
//				time.Sleep(5 * time.Second)
//
//				err := extractor.FindElement(kraken.ByCSSSelector, `a[href="/elonmusk/following"]`).Click()
//				if err != nil {
//					log.Errorf("click following link error: %s", err)
//					return
//				}
//
//				wrapper := extractor.FindElement(kraken.ByCSSSelector, `div[aria-label="Timeline: Following"]`)
//				if err = wrapper.Error(); err != nil {
//					return
//				}
//
//				links := wrapper.FindElements(kraken.ByTagName, "a")
//				if err = links.Error(); err != nil {
//					log.Errorf("find flowing links error: %s", err)
//					return
//				}
//
//				for _, v := range links.Elements() {
//					text, _ := v.Text()
//					fmt.Println(text)
//
//					script := "const event = new MouseEvent('mouseover', { bubbles: true }); arguments[0].dispatchEvent(event);"
//					_, err = extractor.WebDriver().ExecuteScript(script, []interface{}{v})
//					if err != nil {
//						log.Errorf("mouse hover flowing link error: %s", err)
//						return
//					}
//
//					var card kraken.Element
//					card, err = extractor.FindElement(kraken.ByCSSSelector, `div[data-testid="HoverCard"]`).Valid()
//					if err != nil {
//						log.Errorf("find hover card error: %s", err)
//						return
//					}
//
//					var hoverLinks kraken.Elements
//					hoverLinks, err = card.FindElements(kraken.ByTagName, `a`).Valid()
//					if err != nil {
//						return
//					}
//					for i, vv := range hoverLinks.Elements() {
//						tt, _ := vv.Text()
//						if strings.Contains(tt, "Following") || strings.Contains(tt, "Followers") {
//							fmt.Println(hoverLinks.Len(), i, text, tt)
//						}
//					}
//
//					script = "const event = new MouseEvent('mouseout', { bubbles: true }); arguments[0].dispatchEvent(event);"
//					_, err = extractor.WebDriver().ExecuteScript(script, []interface{}{v})
//					if err != nil {
//						log.Errorf("mouse out flowing link error: %s", err)
//						return
//					}
//					time.Sleep(1 * time.Second)
//				}
//
//				time.Sleep(5 * time.Minute)
//			})
//	}
//}

func TestCrawler(t *testing.T) {

	dp, err := fs.Lookupwd("./drivers/chromedriver_130_arm64")
	require.NoError(t, err)

	kraken.DriverPath = dp

	router := kraken.NewRouter(DefaultHandler)

	//router.Handle("/i/flow/login", Login)
	group := router.Group("/i", func(c *kraken.Context) {
		fmt.Println("this is i")
		c.Next()
		fmt.Println("this is ii")
	})
	group.Use(func(c *kraken.Context) {
		fmt.Println("mid1")
		c.Next()
		fmt.Println("mid2")
	})
	group.Handle("/flow/login", Login)

	//router.Handle("/:username", UserHome)
	//router.Handle("/:username/following", UserFollowing)

	err = kraken.Request("https://x.com/elonmusk",
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

func DefaultHandler(c *kraken.Context) {
	fmt.Println("this is default")

	time.Sleep(time.Hour)
}

func Login(c *kraken.Context) {
	fmt.Println("this is login")
}

func UserHome(c *kraken.Context) {
	fmt.Println("this is user home")
	c.Abort(func() bool {
		time.Sleep(5 * time.Second)
		if c.URL.String() != c.Extractor.CurrentURL().String() {
			fmt.Println("this is user home, url redirect")
			return false
		}
		return true
	})

	body, err := c.Extractor.FindElement(kraken.ByTagName, "body").Elem()
	if err != nil {
		return
	}

	height, err := body.ScrollHeight()
	if err != nil {
		return
	}

	for {
		err = body.WaitScrollHeightIncreased(height, 5*time.Second)
		if err != nil {
			if errors.Is(err, kraken.TimoutErr) {
				break
			} else {
				return
			}
		}
		height, err = body.ScrollHeight()
		if err != nil {
			return
		}
		time.Sleep(3 * time.Second)
		err = c.Extractor.ScrollBodyBottom()
		if err != nil {
			return
		}
	}

	time.Sleep(5 * time.Second)
}

func UserFollowing(c *kraken.Context) {
	fmt.Println("this is user following")
}
