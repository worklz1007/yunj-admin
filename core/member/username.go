package member

import (
	"strconv"
	"time"
)

func DeletedUsername(username string, at time.Time) string {
	return username + "_del_" + strconv.FormatInt(at.UnixMilli(), 10)
}
