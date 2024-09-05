package app

import "errors"

var App = &app{}

type app struct {
	Rules []*Rule
}

func (this *app) Register(r *Rule) {
	this.Rules = append(this.Rules, r)
}

func (this *app) Get(index int) *Rule {
	if index >= 0 && index < len(this.Rules) {
		return this.Rules[index]
	}
	return nil
}

func (this *app) Run(index int, op ...func(r *Rule)) error {
	if index >= 0 && index < len(this.Rules) {
		return this.Rules[index].Run()
	}
	return errors.New("无效")
}
