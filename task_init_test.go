package task

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New()
	t.Log(x.SetTask("x", "0,1,2,3 * * * * *", func() {
		t.Log(1)
	}))
	x.Start()
	t.Log(x.SetTask("x2", Date{
		Second: []int{0, 1, 2, 3, 4, 5, 6},
	}.String(), func() {
		t.Log(2)
	}))
	t.Log(x.SetTask("x3", NewIntervalSpec(time.Second*30), func() {
		for _, v := range x.GetTasks() {
			t.Log(v)
		}
	}))
	select {}
}
