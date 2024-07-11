package learning

import "testing"

func TestName(t *testing.T) {
	a := 8
	a++
	s := "hello world"
	for _, c := range s {
		t.Log(c)
	}
}
