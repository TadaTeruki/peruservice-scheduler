package timeconv

import "time"

// UTC
const layout = "2006-01-02T15:04:05.000Z"

func TimeToString(t time.Time) string {
	return t.Format(layout)
}

func StringToTime(s string) (time.Time, error) {
	return time.Parse(layout, s)
}
