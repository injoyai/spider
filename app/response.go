package app

import (
	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/selenium"
	"io"
	"net/http"
	"time"
)

type Action func(ctx *Response)

func (this *Rule) newContext() *Response {
	return &Response{
		rule: this,
		Safe: maps.NewSafe(),
	}
}

type Response struct {
	rule    *Rule
	depth   int
	Request *http.Request
	*http.Response
	*maps.Safe
	context.Context
	cookies   []*http.Cookie
	UserAgent string
}

func (this *Response) DoRequest(method string, url string, body io.Reader, op ...func(r *http.Request)) (*Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for _, f := range op {
		f(req)
	}
	resp, err := this.Client().Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{
		rule:      this.rule,
		depth:     0,
		Request:   req,
		Response:  resp,
		Safe:      this.Safe,
		Context:   this.Context,
		cookies:   this.cookies,
		UserAgent: this.UserAgent,
	}, nil
}

func (this *Response) Client() *Client {
	return this.rule.client
}

/*
Document
使用方法参考
https://blog.csdn.net/qq_38334677/article/details/129225231

	//查找标签: 		doc.Find("body,div,...") 多个用,隔开
	//查找ID: 		doc.Find("#id1")
	//查找class: 	doc.Find(".class1")
	//查找属性: 		doc.Find("div[lang]") doc.Find("div[lang=zh]") doc.Find("div[id][lang=zh]")
	//查找子节点: 	doc.Find("body>div")
	//过滤数据: 		doc.Find("div:contains(xxx)")
	//过滤节点: 		dom.Find("span:has(div)")
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
*/
func (this *Response) Document() (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(this.Body)
}

// Map 把body解析成conv.Map
func (this *Response) Map(codec ...codec.Interface) *conv.Map {
	return conv.NewMap(this.Body, codec...)
}

// Do 发起新的请求
func (this *Response) Do(req Request) {
	if len(this.UserAgent) == 0 {
		if this.Request != nil && len(this.Request.Header.Get("User-Agent")) > 0 {
			this.UserAgent = this.Request.Header.Get("User-Agent")
		} else {
			this.UserAgent = this.rule.getUserAgent()
		}
	}
	if req.Header == nil {
		req.Header = http.Header{}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", this.UserAgent)
	}
	select {
	case this.rule.queue <- &Task{
		Response: this,
		Request:  req,
	}:
	default:
	}
}

// Next 下一步处理响应数据
func (this *Response) Next(by string) {
	action := this.rule.Actions[by]
	if action != nil {
		action(this)
	}
}

// Exit 退出爬取
func (this *Response) Exit() {
	this.rule.Stop()
}

// Output 输出结果
func (this *Response) Output(v any) {
	if this.rule.OnOutput != nil {
		this.rule.OnOutput(v)
	}
}

func (this *Response) Cookies() []*http.Cookie {
	if this.cookies == nil {
		this.cookies = this.Response.Cookies()
	}
	return this.cookies
}

// Chrome 调起浏览器进行操作,例扫码登录
func (this *Response) Chrome(driverPath, browserPath string, option ...selenium.Option) (*selenium.WebDriver, error) {
	//selenium.Debug(false)
	wb, err := selenium.Chrome(driverPath, browserPath, func(e *selenium.Entity) error {
		if len(this.rule.client.Proxy) > 0 {
			e.SetProxy(this.rule.client.Proxy)
			e.SetUserAgent(this.UserAgent)
		}
		for _, v := range option {
			if err := v(e); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = wb.Get(this.Request.URL.String())
	return wb, err
}

func (this *Response) GetCookiesFromChrome(wb *selenium.WebDriver) ([]*http.Cookie, error) {
	cookies, err := wb.GetCookies()
	if err != nil {
		return nil, err
	}
	httpCookies := make([]*http.Cookie, len(cookies))
	for i := range cookies {
		httpCookies[i] = &http.Cookie{
			Name:     cookies[i].Name,
			Value:    cookies[i].Value,
			Path:     cookies[i].Path,
			Domain:   cookies[i].Domain,
			Expires:  time.Unix(int64(cookies[i].Expiry)/1e3, int64(cookies[i].Expiry)%1e3*1e6),
			Secure:   cookies[i].Secure,
			HttpOnly: cookies[i].HTTPOnly,
			SameSite: http.SameSite(conv.Int(cookies[i].SameSite)),
		}
	}
	return httpCookies, nil
}

func (this *Response) GetUserAgentFromChrome(wb *selenium.WebDriver) (string, error) {
	userAgent, err := wb.ExecuteScript("return window.navigator.userAgent", nil)
	return conv.String(userAgent), err
}

// SaveCookies 保存cookie到文件
func (this *Response) SaveCookies(filename string) error {
	return oss.New(filename, this.Cookies())
}

// LoadingCookies 加载本地cookie
func (this *Response) LoadingCookies(filename string) error {
	cookies := []*http.Cookie(nil)
	bs, err := oss.ReadBytes(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bs, &cookies); err != nil {
		return err
	}
	this.cookies = cookies
	return nil
}
