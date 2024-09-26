package selenium

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/injoyai/conv"
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	"io"
	"time"
)

var _ = app.Rule{
	Name: "Selenium",
	Desc: "",
	Root: app.Request{
		By:  "login",
		Url: "https://mojie.app/#/login",
	},
	Actions: map[string]app.Action{
		"login": func(ctx *app.Response) {
			wb, err := ctx.Chrome("./bin/chrome/chromedriver.exe", "./bin/chrome/chrome.exe")
			if err != nil {
				logs.Err(err)
				return
			}
			defer wb.Close()

			<-time.After(time.Second * 2)

			logs.Debug(ctx.GetUserAgentFromChrome(wb))

			logs.Debug(ctx.GetCookiesFromChrome(wb))

			ctx.Do(app.Request{
				By:  "order",
				Url: "https://mojie.app/#/order",
			})
		},
		"order": func(ctx *app.Response) {

			logs.Debug("order")

			bs := conv.Bytes(ctx.Body)
			logs.Debug(string(bs))
			ctx.Body = io.NopCloser(bytes.NewReader(bs))

			doc, err := ctx.Document()
			if err != nil {
				logs.Err(err)
				return
			}

			logs.Debug(doc.Text())

			logs.Debug("done")

			doc.Find("span").Each(func(i int, selection *goquery.Selection) {
				logs.Debug(selection.Text())
			})

		},
	},
}.Register()
