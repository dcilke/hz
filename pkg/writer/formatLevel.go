package writer

import (
	"encoding/json"
	"fmt"
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
	levels := getLevels(m)

	if ok, value := sameOrEmpty(levels[KeyLevel], levels[KeyLog]); ok {
		if value == "" {
			return ""
		}
		return f.format(value)
	}

	return kvJoin(
		f.formatKey(KeyLevel), f.format(levels[KeyLevel]),
		f.formatKey(fmt.Sprintf("%s.%s", KeyLog, KeyLevel)), f.format(levels[KeyLog]),
	)
}

func (f *levelFormatter) ExcludeKeys() []string {
	return f.keys
}

func (f *levelFormatter) format(l string) string {
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

func getLevels(m map[string]any) map[string]string {
	levels := make(map[string]string, 2)
	if i, ok := m[KeyLevel]; ok {
		levels[KeyLevel] = getLevel(i)
	}
	if log, ok := m[KeyLog]; ok {
		if obj, ok := log.(map[string]any); ok {
			if l, ok := obj[KeyLevel]; ok {
				levels[KeyLog] = getLevel(l)
			}
		}
	}
	return levels
}

func getLevel(i any) string {
	if i == nil {
		return ""
	}

	if l, ok := i.(string); ok {
		switch l {
		case LevelPanicStr:
			return LevelPanicStr
		case LevelFatalStr:
			return LevelFatalStr
		case LevelErrorStr:
			return LevelErrorStr
		case LevelWarnStr:
			return LevelWarnStr
		case LevelInfoStr:
			return LevelInfoStr
		case LevelDebugStr:
			return LevelDebugStr
		case LevelTraceStr:
			return LevelTraceStr
		}
	} else if n, ok := i.(json.Number); ok {
		if l, err := n.Int64(); err == nil {
			if l >= LevelPanicNum {
				return LevelPanicStr
			}
			if l >= LevelFatalNum {
				return LevelFatalStr
			}
			if l >= LevelErrorNum {
				return LevelErrorStr
			}
			if l >= LevelWarnNum {
				return LevelWarnStr
			}
			if l >= LevelInfoNum {
				return LevelInfoStr
			}
			if l >= LevelDebugNum {
				return LevelDebugStr
			}
			if l >= LevelTraceNum {
				return LevelTraceStr
			}
		}
	}
	return fmt.Sprintf("%s", i)
}