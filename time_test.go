package learning

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	now := time.Now()
	expTime := now.Add(5 * time.Minute).Format("02-01-2006 15:04:05.999")
	fmt.Println(expTime)
}
