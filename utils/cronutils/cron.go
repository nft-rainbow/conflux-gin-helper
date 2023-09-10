package cronutils

import (
	"fmt"
	"time"

	"github.com/nft-rainbow/conflux-gin-helper/utils"
)

func NextScheduleTime(spec string) (time.Time, error) {
	switch spec {
	case "@daily":
		return utils.TomorrowBegin(), nil
	case "@monthly":
		return utils.BeginnigOfNextMonth(time.Now()), nil
	}
	return time.Time{}, fmt.Errorf("unsupport")
}
