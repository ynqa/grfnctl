package util

import (
	"time"

	"github.com/go-openapi/strfmt"
)

const dateTimeSecondsLayout = "2006-01-02T15:04:05"

// FormatDateTimeInLocal kindly formats the provided DateTime in the local timezone with second precision.
func FormatDateTimeInLocal(dt strfmt.DateTime) string {
	if dt.IsZero() {
		return ""
	}

	localTime := time.Time(dt).In(time.Local)
	return localTime.Format(dateTimeSecondsLayout)
}
