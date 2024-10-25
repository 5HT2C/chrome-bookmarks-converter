package utils

import (
	"strconv"
	"time"
)

type MapCondition struct {
	Condition bool
	K         string
	V         string
}

func MapAppend(m map[string]string, e ...MapCondition) map[string]string {
	if m == nil || len(m) == 0 {
		m = make(map[string]string)
	}

	for _, c := range e {
		if c.Condition {
			m[c.K] = c.V
		}
	}

	return m
}

func StringToTime(s string) time.Time {
	if sec, err := strconv.ParseInt(s, 10, 64); err != nil {
		return time.Unix(0, 0)
	} else {
		return time.Unix(sec, 0)
	}
}
