package app

import "errors"

var App = &app{}

type app struct {
	Rules map[string]*Rule
}

func (this *app) Register(r *Rule) {
	if this.Rules == nil {
		this.Rules = make(map[string]*Rule)
	}
	this.Rules[r.Name] = r
}

func (this *app) Get(key string) *Rule {
	return this.Rules[key]
}

func (this *app) Run(key string, op ...func(r *Rule)) error {
	r := this.Get(key)
	if r != nil {
		for _, v := range op {
			v(r)
		}
		return r.Run()
	}
	return errors.New("规则不存在")
}
