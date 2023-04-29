package formatter

import (
	"encoding/json"
	"time"

	"github.com/dcilke/gu"
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

var _ Formatter = (*Timestamp)(nil)

type Timestamp struct {
	color      bool
	formatKey  Stringer
	timeFormat string
	keys       []string
}

func NewTimestamp(color bool, formatKeys Stringer, timeFormat string) Formatter {
	return &Timestamp{
		color:      color,
		formatKey:  formatKeys,
		timeFormat: timeFormat,
		keys:       []string{KeyTimestamp, KeyAtTimestamp, KeyTime},
	}
}

func (f *Timestamp) Format(m map[string]any) string {
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

	if ok, value := gu.SameOrZero(timestamp, attimestamp, time); ok {
		if value == "" {
			return Colorize(DefaultTimeValue, ColorDarkGray, f.color)
		}
		return Colorize(value, ColorDarkGray, f.color)
	}
	return kvJoin(
		f.formatKey(KeyTimestamp), Colorize(timestamp, ColorDarkGray, f.color),
		f.formatKey(KeyAtTimestamp), Colorize(attimestamp, ColorDarkGray, f.color),
		f.formatKey(KeyTime), Colorize(time, ColorDarkGray, f.color),
	)
}

func (f *Timestamp) ExcludeKeys() []string {
	return f.keys
}

func (f *Timestamp) getTime(i any) string {
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
