package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dcilke/gu"
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

var _ Formatter = (*Level)(nil)

type Level struct {
	color     bool
	formatKey Stringer
	keys      []string
}

func NewLevel(color bool, formatKey Stringer) Formatter {
	return &Level{
		color:     color,
		formatKey: formatKey,
		keys:      []string{KeyLevel, KeyLog + "." + KeyLevel},
	}
}

func (f *Level) Format(m map[string]any) string {
	levels := GetLevels(m)

	if ok, value := gu.SameOrZero(levels[KeyLevel], levels[KeyLog]); ok {
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

func (f *Level) ExcludeKeys() []string {
	return f.keys
}

func (f *Level) format(l string) string {
	switch l {
	case LevelPanicStr:
		return Boldrize(DefaultLevelPanicValue, ColorRed, f.color)
	case LevelFatalStr:
		return Boldrize(DefaultLevelFatalValue, ColorRed, f.color)
	case LevelErrorStr:
		return Boldrize(DefaultLevelErrorValue, ColorRed, f.color)
	case LevelWarnStr:
		return Colorize(DefaultLevelWarnValue, ColorRed, f.color)
	case LevelInfoStr:
		return Colorize(DefaultLevelInfoValue, ColorGreen, f.color)
	case LevelDebugStr:
		return Colorize(DefaultLevelDebugValue, ColorYellow, f.color)
	case LevelTraceStr:
		return Colorize(DefaultLevelTraceValue, ColorMagenta, f.color)
	default:
		ll := strings.ToUpper(l)
		if len(ll) > 3 {
			ll = ll[0:3]
		}
		return Colorize(ll, ColorBold, f.color)
	}
}

func GetLevels(m map[string]any) map[string]string {
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

	if n, ok := i.(json.Number); ok {
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
	} else {
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
		}
	}
	return fmt.Sprintf("%s", i)
}
