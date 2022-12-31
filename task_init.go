package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

// New 新建计时器(任务调度),最小周期秒
func New() *Cron {
	return &Cron{
		Cron: cron.New(cron.WithSeconds()),
		m:    make(map[string]*Task),
	}
}

// Cron 定时器(任务调度),任务起一个协程
type Cron struct {
	*cron.Cron
	m  map[string]*Task
	mu sync.RWMutex
}

// GetTasks 读取全部任务
func (this *Cron) GetTasks() []*Task {
	m := make(map[cron.EntryID]Task)
	this.mu.RLock()
	for _, v := range this.m {
		m[v.ID] = *v
	}
	this.mu.RUnlock()

	taskList := []*Task(nil)
	for _, v := range this.Cron.Entries() {
		taskList = append(taskList, &Task{
			Key:   m[v.ID].Key,
			Spec:  m[v.ID].Spec,
			Entry: v,
		})
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
	en := this.Cron.Entry(task.ID)
	if en.ID == 0 {
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
		return nil
	}
	return &Task{Key: key, Entry: en}
}

// SetTask 设置任务
func (this *Cron) SetTask(key, spec string, task func()) error {
	return this.SetJob(key, spec, cron.FuncJob(task))
}

// SetJob 设置任务
func (this *Cron) SetJob(key, spec string, job cron.Job) error {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if ok {
		//存在相同任务,则移除
		this.Cron.Remove(task.ID)
	}
	id, err := this.Cron.AddJob(spec, job)
	if err != nil {
		return err
	}
	this.mu.Lock()
	this.m[key] = &Task{
		Key:   key,
		Spec:  spec,
		Entry: cron.Entry{ID: id},
	}
	this.mu.Unlock()
	return nil
}

// DelTask 删除任务
func (this *Cron) DelTask(key string) {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if ok {
		this.Cron.Remove(task.ID)
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
	}
}

// Task 任务
type Task struct {
	Key        string //任务唯一标识
	Spec       string //定时规则
	cron.Entry        //任务
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
