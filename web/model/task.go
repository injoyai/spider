package model

import "github.com/injoyai/task"

type Task struct {
	ID        int64            `json:"id"`     //id
	InDate    int64            `json:"inDate"` //创建使劲
	*TaskInfo `xorm:"extends"` //信息

	LastDate int64 `json:"lastDate" xorm:"-"` //上次执行时间
	NextDate int64 `json:"nextDate" xorm:"-"` //下次执行时间
	Enable   bool  `json:"enable" xorm:"-"`   //是否有效(启用)
}

func (this *Task) Response() *Task {

	return this
}

type TaskInfo struct {
	Name string `json:"name"`
	Memo string `json:"memo"`
	Spec string `json:"spec"` //表达式
}

type TaskCreateReq struct {
	*TaskInfo
}

func (this *TaskCreateReq) New() (*Task, string, error) {
	if err := task.CheckSpec(this.Spec); err != nil {
		return nil, "", err
	}
	return &Task{
		TaskInfo: this.TaskInfo,
	}, "Name,Memo,Spec", nil
}

type TaskUpdateReq struct {
	ID int64 `json:"id"`
	*TaskCreateReq
}

// TaskLog 任务执行记录
type TaskLog struct {
}
