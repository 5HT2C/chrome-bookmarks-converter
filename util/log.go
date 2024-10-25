package util

import (
	"log"
	"os"
)

var (
	IsSafe = true
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
)

type LogLevel int

const (
	LogOff LogLevel = iota
	LogInfo
	LogWarn
	LogFuck
)

func (l LogLevel) String() string {
	switch l {
	case LogOff:
		fallthrough
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogFuck:
		return "SEGF" // dude what do you evennn meannnnn
	default:
		return "????"
	}
}

func Log(err error, lvl LogLevel, ctx string, v ...any) {
	fmtLogger := "[%s] logger: %s %s\n"
	fmtCaught := "[%s] caught: %s %s\n%s\n"

	if IsSafe && (lvl == LogOff || lvl >= LogFuck) {
		logger.Fatalf(fmtLogger, lvl, ctx, v, err)
		return
	}

	if err == nil {
		logger.Printf(fmtLogger, lvl, ctx, v)
	} else {
		logger.Printf(fmtCaught, lvl, ctx, v, err)
	}
}
