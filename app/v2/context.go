package v2

import (
	"context"
	"github.com/injoyai/base/maps"
	"net/http"
)

type Context struct {
	*Rule
	*http.Client
	context.Context
	*http.Request
	*http.Response
	Tag   *maps.Safe
	depth int
}

func (this *Context) Cookies() []*http.Cookie {
	return this.Response.Cookies()
}

func (this *Context) GetCookie(name string) *http.Cookie {
	for _, v := range this.Response.Cookies() {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (this *Context) UserAgent() string {
	return this.Request.UserAgent()
}

// Depth 执行深度
func (this *Context) Depth() int {
	return this.depth
}

func (this *Context) Do(r Request) error {
	req, err := r.Request()
	if err != nil {
		return err
	}
	r.Header.Set("User-Agent", this.UserAgent())
	resp, err := this.Client.Do(req)
	if err != nil {
		return err
	}
	return this.Rule.action(r.By, &Context{
		Rule:     this.Rule,
		Client:   this.Client,
		Context:  this.Context,
		Request:  req,
		Response: resp,
		depth:    this.depth + 1,
		Tag:      this.Tag,
	})
}
