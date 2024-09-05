package lib

type Limit struct {
	C    chan struct{}
	Free chan struct{}
}

// NewLimit 同时执行最大数
func NewLimit(limit int) *Limit {
	return &Limit{
		C:    make(chan struct{}, limit),
		Free: make(chan struct{}),
	}
}

// Try 尝试,返回是否到达最大限制
func (this *Limit) Try() bool {
	select {
	case this.C <- struct{}{}:
	default:
		return false
	}
	return true
}

// Add 等待加入成功
func (this *Limit) Add() {
	this.C <- struct{}{}
}

// Done 释放,执行完成
func (this *Limit) Done() {
	select {
	case <-this.C:
		if len(this.C) == 0 {
			select {
			case this.Free <- struct{}{}:
			default:
			}
		}
	default:
	}
}
