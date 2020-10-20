package utils

import "strconv"

func Int32ToString(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}
