package utils

import (
	"fmt"
	"math"
)

const SQL_DATETIME_FORMAT = "2006-01-02 15:04:05"
const SQL_OUTPUT_DATETIMEFORMAT = "2006-01-02T15:04:05Z"

func Suffix(cnt int, singular string, plural string) string {
	ret := plural
	if cnt == 1 || cnt == -1 {
		ret = singular
	}

	return ret
}

func SecToTimePrint(secondCount float64) string {
	var mins = math.Floor(secondCount / 60)
	var hrs = math.Floor(mins / 60)

	var secs = secondCount - (mins * 60)
	mins -= hrs * 60
	return fmt.Sprintf("%02d:%02d:%02d", int(hrs), int(mins), int(secs))
}
