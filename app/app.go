package app

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"time"
)

type Rule struct {
	Name    string           //规则名称
	Desc    string           //规则描述
	Pause   [2]time.Duration //随机暂停时间范围,用于模拟用户操作 例 [10,20]
	Limit   int              //最大并发数
	Depth   int              //爬取深度,小于等于0表示无限制
	Header  http.Header      //请求头
	Timeout time.Duration    //超时时间

	OnStart func(u *url.URL) //开始爬取事件,可以过滤网站等操作
}

type Item struct {
	Method   string
	u        *url.URL
	Client   *http.Client
	GetProxy func() *url.URL
}

func (this *Item) init() {
	if len(this.Method) == 0 {
		this.Method = http.MethodGet
	}
	if this.Client == nil {
		this.Client = &http.Client{
			Transport: &http.Transport{
				//连接结束后会直接关闭,
				//否则会加到连接池复用
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					//设置可以访问HTTPS
					InsecureSkipVerify: true,
				},
			},
			//设置连接超时时间,连接成功后无效
			//连接成功后数据读取时间可以超过这个时间
			//数据读取超时等可以nginx配置
			Timeout: time.Second * 10,
		}
	}
}

// setProxy 设置代理,如果设置了的话
func (this *Item) setProxy() error {
	if this.GetProxy != nil {
		if p := this.GetProxy(); p != nil {
			if transport, ok := this.Client.Transport.(*http.Transport); ok {
				switch p.Scheme {
				case "http", "https":
					transport.Proxy = http.ProxyURL(p)
				case "socks5", "socks5h":
					dialer, err := proxy.FromURL(p, proxy.Direct)
					if err != nil {
						return err
					}
					//transport.DialContext
					transport.Dial = dialer.Dial
				default:
					transport.Proxy = http.ProxyURL(p)
				}
				return nil
			}
		}
	}
	return nil
}

func (this *Item) do() {
	http.NewRequest(this.Method, this.u.String(), nil)
}

type App struct {
	Rule []*Rule
}

func Register(i *Rule) {

}

func init() {
	//从文件中加载规则
	oss.RangeFile()
}

func loadByFile() {

}
