package mojie

import (
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	"net/http"
)

var _ = app.Rule{
	Name: "魔戒",
	Desc: "",
	Root: app.Request{
		By:  "get_info",
		Url: "https://mojie.cyou/api/v1/user/getSubscribe",
		Header: http.Header{
			"Authorization": cache.NewFile("./auth", "auth").GetStrings("auth"),
		},
	},
	Depth: 3,
	Proxy: "http://127.0.0.1:1081",
	Actions: map[string]app.Action{
		"get_info": func(ctx *app.Response) {

			m := ctx.Map()
			msg := m.GetVar("message")
			if msg != nil {

				//登录状态过期，调起浏览器进行登录
				logs.Debug("登录过期，进行登录")
				ctx.Do(app.Request{By: "login"})

				resp, err := ctx.DoRequest("POST", "https://mojie.cyou/api/v1/passport/auth/login", nil)
				if err != nil {
					logs.Err(err)
					return
				}
				auth := resp.Map().GetString("data.token")
				cache.NewFile("./auth", "auth").Set("auth", auth).Save()

				ctx.Do(app.Request{By: "get_info", Url: "https://mojie.cyou/api/v1/user/getSubscribe", Header: http.Header{
					"Authorization": cache.NewFile("./auth", "auth").GetStrings("auth"),
				}})
				return

			}

			logs.Debugf(
				"已用%s/总共%s",
				oss.SizeString(m.GetInt64("data.d")),
				oss.SizeString(m.GetInt64("data.transfer_enable")),
			)

		},
	},
}.Register()
