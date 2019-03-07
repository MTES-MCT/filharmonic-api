package util

import "time"

var mockTime = time.Time{}

func Now() time.Time {
	if !mockTime.IsZero() {
		return mockTime
	}
	return time.Now()
}

func SetTime(t time.Time) {
	mockTime = t
}
