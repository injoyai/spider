package rule

import (
	"github.com/injoyai/spider/app"
	"net/http"
	"time"
)

type Rule struct {
	Name   string           //规则名称
	Desc   string           //规则描述
	Pause  [2]time.Duration //随机暂停时间范围,用于模拟用户操作 例 [10,20]
	Limit  int              //最大并发数
	Header http.Header      //请求头
}

// Register 从代码中加载规则
func (this *Rule) Register() {
	app.Register(this)
}
