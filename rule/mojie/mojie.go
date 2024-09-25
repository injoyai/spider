package mojie

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	"strings"
)

var _ = app.Rule{
	Name: "魔戒",
	Desc: "",
	Root: app.Request{
		By:  "find",
		Url: "https://mojie.cyou/#/dashboard",
	},
	Proxy: "http://127.0.0.1:1081",
	Actions: map[string]app.Action{
		"find": func(ctx *app.Context) {

			doc, err := ctx.Document()
			if err != nil {
				logs.Err(err)
				return
			}

			logs.Debug(doc.Text())

			return

			//判断是否登录，没登陆则调起登录界面进行手动登录

			doc.Find("link,a").Each(func(i int, selection *goquery.Selection) {
				if u, ok := selection.Attr("href"); ok && !ctx.Exist(u) {
					if !strings.HasSuffix(u, ".css") && !strings.HasSuffix(u, ".js") && strings.HasPrefix(u, "http") {
						ctx.Set(u, struct{}{})
						ctx.Do(app.Request{
							By:  "find",
							Url: u,
						})
						ctx.Output(u)
					}
				}
			})

		},
	},
}.Register()
