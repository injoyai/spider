package v2

import (
	"context"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/spider/tool"
	"time"
)

type Rule struct {
	*Config                      //配置信息
	*safe.Runner                 //内部运行机制
	client       *Client         //http客户端
	limit        *tool.Limit     //限制协程数量
	queue        chan func()     //任务队列
	ctx          context.Context //上下文
}

// UserAgent 用户标识
func (this *Rule) UserAgent() string {
	if len(this.UserAgents) > 0 {
		return this.UserAgents[this.Rand.Intn(len(this.UserAgents))]
	}
	return DefaultUserAgent
}

// Pause 随机暂停,模拟人为
func (this *Rule) Pause() {
	min := int64(this.Config.Pause[0])
	max := int64(this.Config.Pause[1])
	if max <= min {
		return
	}
	t := time.Duration(this.Rand.Int63n(max-min)) + this.Config.Pause[0]
	time.Sleep(t)
}

func (this *Rule) action(by string, ctx *Context) error {
	action := this.Actions[by]
	if action == nil {
		return fmt.Errorf("动作[%s]不存在", by)
	}
	this.Pause()
	return action(ctx)
}

func (this *Rule) run(ctx context.Context) error {

	err := (&Context{
		Rule:    this,
		Context: this.ctx,
		Tag:     maps.NewSafe(),
	}).Do(this.Root)
	if err != nil {
		return err
	}

	for {
		//随机暂停,模拟人工操作
		this.Pause()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case f := <-this.queue:
			f()
		case <-this.limit.Free:
			//需要一个判断全部协程执行完成
			return nil
		}
	}
}
