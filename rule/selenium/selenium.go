package selenium

import (
	"github.com/injoyai/spider/app"
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
	Cookie:        nil,
	Timeout:       0,
	Proxy:         nil,
	Log:           nil,
	Root:          nil,
	Actions:       nil,
	OnOutput:      nil,
	Runner:        nil,
}.Register()
