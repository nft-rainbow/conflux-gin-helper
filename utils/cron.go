package utils

import (
	"fmt"
	"time"
)

func NextScheduleTime(spec string) (time.Time, error) {
	switch spec {
	case "@monthly":
		return BeginnigOfNextMonth(time.Now()), nil
	case "@daily":
		return TomorrowBegin(), nil
	}
	return time.Time{}, fmt.Errorf("unsupport")
}
