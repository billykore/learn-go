package templates

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	v, err := time.Parse(time.DateOnly, "2024-01-31")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v.Add(24 * time.Hour).Day())
}
