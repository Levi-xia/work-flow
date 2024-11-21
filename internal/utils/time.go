package utils

import "time"

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func FormatTime(time time.Time, layout string) string {
	return time.Format(layout)
}
