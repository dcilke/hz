package writer

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	KeyLevel = "level"
	KeyLog   = "log"

	LevelTraceStr = "trace"
	LevelDebugStr = "debug"
	LevelInfoStr  = "info"
	LevelWarnStr  = "warn"
	LevelErrorStr = "error"
	LevelFatalStr = "fatal"
	LevelPanicStr = "panic"

	LevelTraceNum = 10
	LevelDebugNum = 20
	LevelInfoNum  = 30
	LevelWarnNum  = 40
	LevelErrorNum = 50
	LevelFatalNum = 60
	LevelPanicNum = 100

	DefaultLevelValue      = "???"
	DefaultLevelTraceValue = "TRC"
	DefaultLevelDebugValue = "DBG"
	DefaultLevelInfoValue  = "INF"
	DefaultLevelWarnValue  = "WRN"
	DefaultLevelErrorValue = "ERR"
	DefaultLevelFatalValue = "FTL"
	DefaultLevelPanicValue = "PNC"
)

type levelFormatter struct {
	noColor   bool
	formatKey Stringer
	keys      []string
}

func newLevelFormatter(noColor bool, formatKey Stringer) Formatter {
	return &levelFormatter{
		noColor:   noColor,
		formatKey: formatKey,
		keys:      []string{KeyLevel, KeyLog},
	}
}

func (f *levelFormatter) Format(m map[string]any, _ string) string {
	var level string
	var loglevel string
	if i, ok := m[KeyLevel]; ok {
		level = f.getLevel(i)
	}
	if log, ok := m[KeyLog]; ok {
		if obj, ok := log.(map[string]any); ok {
			if l, ok := obj[KeyLevel]; ok {
				loglevel = f.getLevel(l)
			}
		}
	}
	if ok, value := sameOrEmpty(level, loglevel); ok {
		if value == "" {
			return ""
		}
		return value
	}

	return kvJoin(
		f.formatKey(KeyLevel), level,
		f.formatKey(fmt.Sprintf("%s.%s", KeyLog, KeyLevel)), loglevel,
	)
}

func (f *levelFormatter) ExcludeKeys() []string {
	return f.keys
}

func (f *levelFormatter) getLevel(i any) string {
	if i == nil {
		return ""
	}
	if l, ok := i.(string); ok {
		return f.formatStrAsLevelString(l)
	}
	if n, ok := i.(json.Number); ok {
		if l, err := n.Int64(); err == nil {
			return f.formatNumAsLevelString(l)
		}
	}
	l := strings.ToUpper(fmt.Sprintf("%s", i))
	if len(l) > 3 {
		l = l[0:3]
	}
	return l
}

func (f *levelFormatter) formatStrAsLevelString(l string) string {
	switch l {
	case LevelPanicStr:
		return boldrize(DefaultLevelPanicValue, ColorRed, f.noColor)
	case LevelFatalStr:
		return boldrize(DefaultLevelFatalValue, ColorRed, f.noColor)
	case LevelErrorStr:
		return boldrize(DefaultLevelErrorValue, ColorRed, f.noColor)
	case LevelWarnStr:
		return colorize(DefaultLevelWarnValue, ColorRed, f.noColor)
	case LevelInfoStr:
		return colorize(DefaultLevelInfoValue, ColorGreen, f.noColor)
	case LevelDebugStr:
		return colorize(DefaultLevelDebugValue, ColorYellow, f.noColor)
	case LevelTraceStr:
		return colorize(DefaultLevelTraceValue, ColorMagenta, f.noColor)
	default:
		return colorize(DefaultLevelValue, ColorBold, f.noColor)
	}
}

func (f levelFormatter) formatNumAsLevelString(l int64) string {
	if l >= LevelPanicNum {
		return boldrize(DefaultLevelPanicValue, ColorRed, f.noColor)
	}
	if l >= LevelFatalNum {
		return boldrize(DefaultLevelFatalValue, ColorRed, f.noColor)
	}
	if l >= LevelErrorNum {
		return boldrize(DefaultLevelErrorValue, ColorRed, f.noColor)
	}
	if l >= LevelWarnNum {
		return colorize(DefaultLevelWarnValue, ColorRed, f.noColor)
	}
	if l >= LevelInfoNum {
		return colorize(DefaultLevelInfoValue, ColorGreen, f.noColor)
	}
	if l >= LevelDebugNum {
		return colorize(DefaultLevelDebugValue, ColorYellow, f.noColor)
	}
	if l >= LevelTraceNum {
		return colorize(DefaultLevelTraceValue, ColorMagenta, f.noColor)
	}
	return colorize(DefaultLevelValue, ColorBold, f.noColor)
}
