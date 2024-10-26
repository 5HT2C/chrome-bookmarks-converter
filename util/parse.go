package util

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
			Log(LogInfo, "MapAppend() created", c)
			m[c.K] = c.V
		}
	}

	return m
}

func StringToTime(s string) time.Time {
	if sec, err := strconv.ParseInt(s, 10, 64); err != nil {
		return time.Time{}
	} else {
		return time.Unix(sec, 0)
	}
}

func StringEmpty(s string) bool {
	return len(s) == 0 || s == "0"
}

func StringEmptyScore(s string) int64 {
	if StringEmpty(s) {
		return 0
	}

	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}

	return int64(len(s))
}

func StringConditional(s, d string, c bool) string {
	if !c {
		Log(LogInfo, "StringConditional defaulted", c, s, d)
		return d
	}

	return s
}

func StringOrDefault(s, d string) string {
	return StringConditional(s, d, !StringEmpty(s))
}
