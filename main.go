package main

import (
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	_ "github.com/injoyai/spider/rule/demo"
	_ "github.com/injoyai/spider/rule/selenium"
)

func main() {

	r := app.App.Get(1)
	logs.Err(r.Run())
}
