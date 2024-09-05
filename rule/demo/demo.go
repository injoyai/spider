package demo

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	"strings"
)

var _ = app.Rule{
	Name: "演示",
	Desc: "",
	Root: &app.Request{
		By:  "find",
		Url: "https://blog.csdn.net/qq_38334677/article/details/129225231",
	},
	Actions: map[string]app.Action{
		"find": func(ctx *app.Context) {

			doc, err := ctx.Document()
			if err != nil {
				logs.Err(err)
			}

			doc.Find("link,a").Each(func(i int, selection *goquery.Selection) {
				if u, ok := selection.Attr("href"); ok && !ctx.Exist(u) {
					if !strings.HasSuffix(u, ".css") && !strings.HasSuffix(u, ".js") && strings.HasPrefix(u, "http") {
						ctx.Set(u, struct{}{})
						ctx.Do(&app.Request{
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
