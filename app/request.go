package app

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	By     string //让哪个动作处理请求结果,必选
	Method string
	Url    string
	Header http.Header
	Body   []byte
	Cookie []*http.Cookie
}

func (this *Request) GetMethod() string {
	if len(this.Method) == 0 {
		return http.MethodGet
	}
	return this.Method
}

func (this *Request) GetBody() io.Reader {
	if this.Body == nil {
		return nil
	}
	return bytes.NewReader(this.Body)
}

type Task struct {
	*Context
	Request
}

func (this *Task) Do() {
	this.rule.limit.Add()
	this.rule.Log.Printf("[信息] 开始请求[%s]\n", this.Url)
	go func() {
		defer this.rule.limit.Done()
		req, resp, err := this.do()
		if err != nil {
			this.rule.Log.Printf("[错误] 请求地址[%s]失败: %v\n", this.Url, err)
			return
		}
		//todo 这里可能会阻塞,当全部协程在执行时,并吧队列插入满,则会阻塞,
		//todo 协程停止不掉(等待插入队列),队列插不进去(等待协程释放)
		this.rule.doAction(this.Context, this.By, req, resp)
	}()
}

func (this *Task) do() (*http.Request, *http.Response, error) {
	if this.rule.Depth > 0 && this.depth > this.rule.Depth {
		return nil, nil, errors.New("到达最大深度")
	}
	req, err := http.NewRequest(this.GetMethod(), this.Url, this.GetBody())
	if err != nil {
		return nil, nil, err
	}
	for k, v := range this.rule.Header {
		req.Header.Add(k, strings.Join(v, ","))
	}
	for _, v := range this.Request.Cookie {
		req.AddCookie(v)
	}
	if this.rule.DisableCookie {
		req.Header.Del("Cookie")
	}
	resp, err := this.rule.client.Do(req)
	return req, resp, err
}
