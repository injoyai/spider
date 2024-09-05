package main

import (
	"github.com/injoyai/logs"
	"github.com/injoyai/spider/app"
	_ "github.com/injoyai/spider/rule/demo"
)

func main() {

	r := app.App.Get(0)
	logs.Err(r.Run())
}
