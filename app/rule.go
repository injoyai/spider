package app

import (
	"context"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/spider/tool"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
)

type Rule struct {
	Name          string            //规则名称,展示用
	Desc          string            //规则描述,展示用
	Pause         [2]time.Duration  //随机暂停时间范围,用于模拟用户操作 例1~5秒随机暂停 [time.Second,time.Second*5]
	Limit         int               //最大并发数
	Depth         int               //爬取深度,小于等于0表示无限制
	Header        http.Header       //请求头
	UserAgents    []string          //如果填写则使用取其中一个数据,否则使用默认值
	DisableCookie bool              //全局禁用cookie
	Timeout       time.Duration     //HTTP请求超时时间,默认不超时
	Proxy         string            //代理地址,空为不代理
	Log           *log.Logger       //日志
	Root          Request           //根请求
	Actions       map[string]Action //不同的动作(解析/再次请求)
	OnOutput      func(any)         //输出

	//内部字段
	*safe.Runner                 //内部运行机制
	client       *Client         //http客户端
	rand         *rand.Rand      //用于生成随机数
	limit        *tool.Limit     //限制协程数量
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
	r.client = newClient(r.Timeout, r.Proxy)
	r.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	r.limit = tool.NewLimit(r.Limit)
	r.queue = make(chan *Task, 1000)
	r.ctx = context.Background()
	App.Register(r)
	return r
}

func (this *Rule) getUserAgent() string {
	if len(this.UserAgents) > 0 {
		return this.UserAgents[this.rand.Intn(len(this.UserAgents))]
	}
	return DefaultUserAgent
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

func (this *Rule) doAction(ctx *Response, by string, req *http.Request, resp *http.Response) {
	action := this.Actions[by]
	if action == nil {
		this.Log.Printf("[警告] 动作[%s]不存在\n", by)
		return
	}
	action(&Response{
		rule:     this,
		Request:  req,
		Response: resp,
		Context:  ctx.Context,
		Safe:     ctx.Safe,
		depth:    ctx.depth + 1,
	})
}

func (this *Rule) run(ctx context.Context) error {

	(&Response{
		rule:    this,
		Context: this.ctx,
		Safe:    maps.NewSafe(),
	}).Do(this.Root)

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
