package selenium

import (
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	"github.com/injoyai/spider/tool"
	"net/http"
	"time"
)

var _ = app.Rule{
	Name:          "Selenium",
	Desc:          "",
	Pause:         [2]time.Duration{},
	Limit:         0,
	Depth:         0,
	Header:        nil,
	DisableCookie: false,
	Cookie: func() []*http.Cookie {
		cookies, err := tool.LoadingCookies("./data/cookie/selenium.json")
		logs.PrintErr(err)
		return cookies
	}(),
	Timeout:  0,
	Proxy:    nil,
	Log:      nil,
	Root:     nil,
	Actions:  nil,
	OnOutput: nil,
}.Register()
