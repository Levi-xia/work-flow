package utils

import "time"

func ParseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, timeStr)
}

func FormatTime(time time.Time, layout string) string {
	if time.IsZero() {
		return ""
	}
	return time.Format(layout)
}
