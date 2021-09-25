package logger

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

/*
 Used to fetch values dynamically at runtime
*/

var (
	DefaultCaller = Caller(2)
	DefaultTimer  = Timer(time.RFC3339)
)

// Valuer is returns a log value.
type Valuer func() (key string, val interface{})

// Caller returns returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func() (string, interface{}) {
		_, file, line, _ := runtime.Caller(depth)
		if strings.LastIndex(file, "github.com/driftingboy/tsdb/logger") > 0 {
			_, file, line, _ = runtime.Caller(depth + 1)
		}
		// _ = strings.LastIndexByte(file, '/')
		return "Caller", file[:] + ":" + strconv.Itoa(line)
	}
}

func Timer(layout string) Valuer {
	return func() (key string, val interface{}) {
		return "Timer", time.Now().Format(layout)
	}
}
