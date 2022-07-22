package helper

import (
	"fmt"
	"time"
)

const (
	// "2 January 2006"
	TIME_LAYOUT_DATE = "2 January 2006"

	// "24Hr:Mn:Sc"
	TIME_LAYOUT_HOUR = "15:04:05"
	
	// "2 January 2006, 24Hr:Mn:Sc"
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

// turn any time string to WIB_TIME
// it is recommended to use time layout from this package
func ParseTimeToWIB(timeString string, layout string) (string, error) {
	output, err := time.Parse(layout, timeString)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(output), nil
}

func GetCurrentZuluTime(timeString string, layout string) (string) {
	return fmt.Sprint(time.Now().UTC())
}