package utils

import (
	"time"
)

func UnixMili() int64 {
	return time.Now().UnixNano() / 1000000
}
