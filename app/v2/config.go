package v2

import (
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

type Config struct {
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
	Rand          *rand.Rand        //随机数生成器
}

func (this Config) Register() *Rule {
	c := &this
	if c.Log == nil {
		c.Log = log.Default()
	}
	if c.Limit <= 0 {
		c.Limit = 20
	}
	if c.Rand == nil {
		c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	e := &Rule{
		Config: c,
		Runner: safe.NewRunner(nil),
		client: newClient(c.Timeout, c.Proxy),
		limit:  tool.NewLimit(c.Limit),
		queue:  make(chan func(), 1000),
	}
	App.Register(e)
	return e
}

type Action func(ctx *Context) error
