package writer

import (
	"encoding/json"
	"time"
)

const (
	KeyTime        = "time"
	KeyTimestamp   = "timestamp"
	KeyAtTimestamp = "@timestamp"

	TimeFormat          = time.RFC3339
	TimeFormatUnixMs    = "UNIXMS"
	TimeFormatUnixMicro = "UNIXMICRO"

	DefaultTimeValue = "<nil>"
)

type timestampFormatter struct {
	noColor    bool
	formatKey  Stringer
	timeFormat string
	keys       []string
}

func newTimestampFormatter(noColor bool, formatKeys Stringer, timeFormat string) Formatter {
	return &timestampFormatter{
		noColor:    noColor,
		formatKey:  formatKeys,
		timeFormat: timeFormat,
		keys:       []string{KeyTimestamp, KeyAtTimestamp, KeyTime},
	}
}

func (f *timestampFormatter) Format(m map[string]any, _ string) string {
	var timestamp string
	var attimestamp string
	var time string
	if i, ok := m[KeyTimestamp]; ok {
		timestamp = f.getTime(i)
	}
	if i, ok := m[KeyTime]; ok {
		time = f.getTime(i)
	}
	if i, ok := m[KeyAtTimestamp]; ok {
		attimestamp = f.getTime(i)
	}

	if ok, value := sameOrEmpty(timestamp, attimestamp, time); ok {
		if value == "" {
			return colorize(DefaultTimeValue, ColorDarkGray, f.noColor)
		}
		return colorize(value, ColorDarkGray, f.noColor)
	}
	return kvJoin(
		f.formatKey(KeyTimestamp), colorize(timestamp, ColorDarkGray, f.noColor),
		f.formatKey(KeyAtTimestamp), colorize(attimestamp, ColorDarkGray, f.noColor),
		f.formatKey(KeyTime), colorize(time, ColorDarkGray, f.noColor),
	)
}

func (f *timestampFormatter) ExcludeKeys() []string {
	return f.keys
}

func (f *timestampFormatter) getTime(i any) string {
	t := ""
	switch tt := i.(type) {
	case string:
		ts, err := time.Parse(TimeFormat, tt)
		if err != nil {
			t = tt
		} else {
			t = ts.Format(f.timeFormat)
		}
	case json.Number:
		i, err := tt.Int64()
		if err != nil {
			t = tt.String()
		} else {
			var sec, nsec int64 = i, 0
			switch TimeFormat {
			case TimeFormatUnixMs:
				nsec = int64(time.Duration(i) * time.Millisecond)
				sec = 0
			case TimeFormatUnixMicro:
				nsec = int64(time.Duration(i) * time.Microsecond)
				sec = 0
			}
			ts := time.Unix(sec, nsec).UTC()
			t = ts.Format(f.timeFormat)
		}
	}
	return t
}
