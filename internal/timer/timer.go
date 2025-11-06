package timer

import (
	"fmt"
	"time"
)

func Now() time.Time {
	vnLocation, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	return time.Now().In(vnLocation)
}

func NowString() string {
	vnLocation, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	t := time.Now().In(vnLocation)

	return fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
}
