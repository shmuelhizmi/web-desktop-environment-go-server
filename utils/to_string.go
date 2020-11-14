package utils

import (
	"fmt"
	"strconv"
)

func Int32ToString(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func UInt64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func FloatToString(value float64) string {
	return fmt.Sprintf("%f", value)
}
