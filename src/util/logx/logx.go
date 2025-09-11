package logx

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type level string

const (
	levelDebug level = "DEBUG"
	levelInfo  level = "INFO"
	levelWarn  level = "WARN"
	levelError level = "ERROR"
)

var (
	asJSON    = false
	minLevel  = levelInfo
	writeLock sync.Mutex
)

func SetJSON(enabled bool) { asJSON = enabled }

func SetLevel(l string) {
	switch l {
	case "debug":
		minLevel = levelDebug
	case "info":
		minLevel = levelInfo
	case "warn":
		minLevel = levelWarn
	case "error":
		minLevel = levelError
	}
}

func shouldLog(l level) bool {
	order := map[level]int{levelDebug: 1, levelInfo: 2, levelWarn: 3, levelError: 4}
	return order[l] >= order[minLevel]
}

func logf(l level, format string, a ...any) {
	if !shouldLog(l) {
		return
	}
	writeLock.Lock()
	defer writeLock.Unlock()
	if asJSON {
		rec := map[string]any{
			"ts":    time.Now().Format(time.RFC3339Nano),
			"level": string(l),
			"msg":   fmt.Sprintf(format, a...),
		}
		enc := json.NewEncoder(os.Stderr)
		_ = enc.Encode(rec)
		return
	}
	fmt.Fprintf(os.Stderr, "%s %-5s %s\n", time.Now().Format(time.RFC3339), string(l), fmt.Sprintf(format, a...))
}

func Debugf(format string, a ...any) { logf(levelDebug, format, a...) }
func Infof(format string, a ...any)  { logf(levelInfo, format, a...) }
func Warnf(format string, a ...any)  { logf(levelWarn, format, a...) }
func Errorf(format string, a ...any) { logf(levelError, format, a...) }
