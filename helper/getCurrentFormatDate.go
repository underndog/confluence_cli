package helper

import "time"

func getCurrentFormattedDate(format string) string {
	return time.Now().Format(format)
}
