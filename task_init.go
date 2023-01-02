package task

import (
	"fmt"
	"github.com/injoyai/conv"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CheckSpec(spec string) error {
	_, err := cron.ParseStandard(spec)
	return err
}

// New 新建计时器(任务调度),最小周期秒
func New() *Cron {
	return &Cron{
		c: cron.New(cron.WithSeconds()),
		m: make(map[string]*Task),
	}
}

// Cron 定时器(任务调度),任务起一个协程
type Cron struct {
	c  *cron.Cron
	m  map[string]*Task
	mu sync.RWMutex
}

// GetTaskAll 读取全部任务
func (this *Cron) GetTaskAll() []*Task {
	m := make(map[cron.EntryID]*Task)
	this.mu.RLock()
	for _, v := range this.m {
		m[v.ID] = v
	}
	this.mu.RUnlock()

	taskList := []*Task(nil)
	for _, v := range this.c.Entries() {
		task, ok := m[v.ID]
		if !ok {
			//删除不在档的任务
			this.c.Remove(v.ID)
			continue
		}
		taskList = append(taskList, task.SetEntry(v))
	}
	return taskList
}

// GetTask 读取任务
func (this *Cron) GetTask(key string) *Task {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if !ok {
		return nil
	}
	en := this.c.Entry(task.ID)
	if en.ID == 0 {
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
		return nil
	}
	return task.SetEntry(en)
}

// SetTask 设置任务
func (this *Cron) SetTask(key, spec string, handler func(), data ...interface{}) error {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if ok {
		//存在相同任务,则移除
		this.c.Remove(task.ID)
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
	}
	id, err := this.c.AddFunc(spec, handler)
	if err != nil {
		return err
	}
	this.mu.Lock()
	this.m[key] = newTask(key, spec, cron.Entry{ID: id}, data...)
	this.mu.Unlock()
	return nil
}

// DelTask 删除任务
func (this *Cron) DelTask(key string) {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if ok {
		this.c.Remove(task.ID)
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
	}
}

// Task 任务
type Task struct {
	Key        string      //任务唯一标识
	Spec       string      //定时规则
	Data       interface{} //自定义数据
	cron.Entry             //任务
}

func newTask(key, spec string, e cron.Entry, v ...interface{}) *Task {
	return &Task{
		Key:   key,
		Spec:  spec,
		Data:  conv.GetDefaultInterface(nil, v...),
		Entry: e,
	}
}

func (this *Task) SetEntry(e cron.Entry) *Task {
	this.Entry = e
	return this
}

func (this *Task) String() string {
	return fmt.Sprintf("名称(%s),生效(%v),规则(%s),上次执行时间(%v),下次执行时间(%v)",
		this.Key, this.Valid(), this.Spec, this.timeStr(this.Prev), this.timeStr(this.Next))
}

func (this *Task) timeStr(t time.Time) string {
	if t.IsZero() {
		return "无"
	}
	return t.String()
}

// Interval 间隔时间
type Interval time.Duration

func (this Interval) String() string {
	return fmt.Sprintf("@every %s", time.Duration(this))
}

// NewIntervalSpec 新建间隔任务
func NewIntervalSpec(t time.Duration) string {
	return Interval(t).String()
}

// Date 按日志执行
type Date struct {
	Month  []int //月 1, 12
	Week   []int //周 0, 6
	Day    []int //天 1, 31
	Hour   []int //时 0, 23
	Minute []int //分 0, 59
	Second []int //秒 0, 59
}

func (this Date) spec(ints []int) string {
	if len(ints) > 0 {
		list := make([]string, len(ints))
		for i, v := range ints {
			list[i] = strconv.Itoa(v)
		}
		return strings.Join(list, ",")
	}
	return "*"
}

func (this Date) String() string {
	return strings.Join([]string{
		this.spec(this.Second),
		this.spec(this.Minute),
		this.spec(this.Hour),
		this.spec(this.Day),
		this.spec(this.Month),
		this.spec(this.Week),
	}, " ")
}
