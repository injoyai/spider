package app

import (
	"context"
	"crypto/tls"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/spider/lib"
	"golang.org/x/net/proxy"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Rule struct {
	Name          string            //规则名称,展示用
	Desc          string            //规则描述,展示用
	Pause         [2]time.Duration  //随机暂停时间范围,用于模拟用户操作 例1~5秒随机暂停 [time.Second,time.Second*5]
	Limit         int               //最大并发数
	Depth         int               //爬取深度,小于等于0表示无限制
	Header        http.Header       //请求头
	DisableCookie bool              //禁用cookie
	Cookie        []*http.Cookie    //自定义cookie
	Timeout       time.Duration     //HTTP请求超时时间,默认不超时
	Proxy         []string          //代理地址
	Log           *log.Logger       //日志
	Root          *Request          //根请求
	Actions       map[string]Action //不同的动作(解析/再次请求)
	OnOutput      func(any)         //输出

	//内部字段
	*safe.Runner                 //内部运行机制
	client       *http.Client    //http客户端
	rand         *rand.Rand      //用于生成随机数
	limit        *lib.Limit      //限制协程数量
	queue        chan *Task      //任务队列
	ctx          context.Context //
}

func (this Rule) Register() *Rule {
	r := &this
	if r.Log == nil {
		r.Log = log.Default()
	}
	if r.Limit <= 0 {
		r.Limit = 20
	}

	r.Runner = safe.NewRunner(nil)
	r.Runner.SetFunc(r.run)
	r.initClient()
	r.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	r.limit = lib.NewLimit(r.Limit)
	r.queue = make(chan *Task, 1000)
	r.ctx = context.Background()
	App.Register(r)
	return r
}

func (this *Rule) pause() {
	min := int64(this.Pause[0])
	max := int64(this.Pause[1])
	if max <= min {
		return
	}
	t := time.Duration(this.rand.Int63n(max-min)) + this.Pause[0]
	time.Sleep(t)
}

func (this *Rule) doAction(ctx *Context, by string, req *http.Request, resp *http.Response) {
	action := this.Actions[by]
	if action == nil {
		this.Log.Printf("[警告] 动作[%s]不存在\n", by)
		return
	}
	action(&Context{
		rule:     this,
		Request:  req,
		Response: resp,
		Context:  ctx.Context,
		Safe:     ctx.Safe,
		depth:    ctx.depth + 1,
	})
}

func (this *Rule) initClient() {
	transport := &http.Transport{
		//连接结束后会直接关闭,不复用
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			//设置可以访问HTTPS
			InsecureSkipVerify: true,
		},
	}
	for _, v := range this.Proxy {
		p, err := url.Parse(v)
		if err == nil {
			switch p.Scheme {
			case "http", "https":
				transport.Proxy = http.ProxyURL(p)
			case "socks5", "socks5h":
				dialer, err := proxy.FromURL(p, proxy.Direct)
				if err != nil {
					continue
				}
				transport.Dial = dialer.Dial
			default:
				transport.Proxy = http.ProxyURL(p)
			}
			break
		}
	}
	this.client = &http.Client{
		Transport: transport,
		Timeout:   this.Timeout,
	}
}

func (this *Rule) run(ctx context.Context) error {
	c := &Context{
		rule:    this,
		Context: this.ctx,
		Safe:    maps.NewSafe(),
	}
	c.Do(this.Root)

	for {
		//随机暂停,模拟人工操作
		this.pause()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task := <-this.queue:
			task.Do()
		case <-this.limit.Free:
			//需要一个判断全部协程执行完成
			return nil
		}
	}
}
