package e

import "time"

const (
	DEFAULT_TIME = "1000-01-01 00:00:00"
)

func DefaultTime() time.Time {
	ti, err := time.Parse("2006-01-02 03:04:05", DEFAULT_TIME)
	if err != nil {
		return time.Now()
	}
	return ti
}
