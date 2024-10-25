package util

import (
	"log"
	"os"
)

var (
	LogModeSafe = true // panic on fatal errors
	LogLvlDebug = true
	logger      = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
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

func Log(lvl LogLevel, ctx string, v ...any) {
	fmtLogger := "[%s] logger: %s %s\n"
	fmtCaught := "[%s] caught: %s %s\n%s\n"
	fmtArgs := make([]any, 0)

	var err error = nil
	for _, arg := range v {
		if err != nil {
			fmtArgs = append(fmtArgs, v)
			continue
		}

		switch arg.(type) {
		case error:
			if arg != nil {
				err = arg.(error)
				continue
			}
		}

		fmtArgs = append(fmtArgs, v)
	}

	if LogModeSafe && (lvl == LogOff || lvl >= LogFuck) {
		logger.Fatalf(fmtCaught, lvl, ctx, err, append(make([]any, len(fmtArgs)), fmtArgs...))
		return
	}

	if err == nil {
		logger.Printf(fmtLogger, lvl, ctx, v)
	} else {
		logger.Printf(fmtCaught, lvl, ctx, err, append(make([]any, len(fmtArgs)), fmtArgs...))
	}
}
