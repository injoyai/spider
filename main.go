package main

import (
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	_ "github.com/injoyai/spider/rule/demo"
	_ "github.com/injoyai/spider/rule/mojie"
	_ "github.com/injoyai/spider/rule/selenium"
)

func main() {
	logs.Err(app.App.Run("魔戒", func(r *app.Rule) {
		r.Proxy = "http://127.0.0.1:1081"
	}))
}
