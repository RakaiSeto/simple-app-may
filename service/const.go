package service

import "time"

const (
	TIME_LAYOUT_DATE = "2 January 2006"
	TIME_LAYOUT_HOUR = "15:04:05"
	TIME_LAYOUT_ALL  = "2 January 2006, 15:04:05"
)

var WIB_TIME *time.Location
var err error

func init() {
	WIB_TIME, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
}