package learning

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type singleton struct{}

var (
	once     sync.Once
	instance *singleton
)

func newSingleton() *singleton {
	once.Do(func() {
		fmt.Println("initialize instance")
		instance = &singleton{}
	})
	fmt.Println("return instance")
	return instance
}

func TestLazy(t *testing.T) {
	obj := newSingleton()

	if obj == nil {
		t.Fatalf("should not nil")
	}
}

func TestSingleton(t *testing.T) {
	for i := 0; i < 100; i++ {
		_ = newSingleton()
	}
}

func TestTime(t *testing.T) {
	y, m, d := time.Now().Date()
	t.Log(y)
	t.Log(m)
	t.Log(d)
}

func TestNow(t *testing.T) {
	a := time.Date(2023, 12, 24, 0, 0, 0, 0, time.Local)
	t.Log(a.Format("02"))
}

func TestAddOneMonth(t *testing.T) {
	a := time.Now()
	b := a.AddDate(0, 1, 0)
	t.Log(a)
	t.Log(b)
}
