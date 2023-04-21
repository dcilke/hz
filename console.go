// forked from https://github.com/rs/zerolog/blob/4099072c03f2f4e61fa08f70adf9a25983f0cd8e/console.go
// Added support for numeric log levels and legecs log levels
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
	colorBold     = 1
	colorDarkGray = 90

	LevelTraceValue = "trace"
	LevelDebugValue = "debug"
	LevelInfoValue  = "info"
	LevelWarnValue  = "warn"
	LevelErrorValue = "error"
	LevelFatalValue = "fatal"
	LevelPanicValue = "panic"

	LevelTraceNum = 10
	LevelDebugNum = 20
	LevelInfoNum  = 30
	LevelWarnNum  = 40
	LevelErrorNum = 50
	LevelFatalNum = 60
	LevelPanicNum = 100

	FieldNameTime        = "time"
	FieldNameTimestamp   = "timestamp"
	FieldNameAtTimestamp = "@timestamp"
	FieldNameLevel       = "level"
	FieldNameMessage     = "message"
	FieldNameError       = "error"
	FieldNameErr         = "err"
	FieldNameCaller      = "caller"
	FieldNameStack       = "stack"
	FieldNameLog         = "log"

	TimeFieldFormat     = time.RFC3339
	TimeFormatUnixMs    = "UNIXMS"
	TimeFormatUnixMicro = "UNIXMICRO"
)

var (
	consoleBufPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

const (
	consoleDefaultTimeFormat = "15:04:05"
)

// Formatter transforms the input into a formatted string.
type Formatter func(any) string

// ConsoleWriter parses the JSON input and writes it in an
// (optionally) colorized, human-friendly format to Out.
type ConsoleWriter struct {
	// out is the output destination.
	out io.Writer

	// noColor disables the colorized output.
	noColor bool

	// timeFormat specifies the format for timestamp in output.
	timeFormat string

	// partsOrder defines the order of parts in output.
	partsOrder []string

	// partsExclude defines parts to not display in output.
	partsExclude []string

	// fieldsExclude defines contextual fields to not display in output.
	fieldsExclude []string

	formatTimestamp     Formatter
	formatLevel         Formatter
	formatLog           Formatter
	formatCaller        Formatter
	formatMessage       Formatter
	formatFieldName     Formatter
	formatFieldValue    Formatter
	formatErrFieldName  Formatter
	formatErrFieldValue Formatter
}

// NewConsoleWriter creates and initializes a new ConsoleWriter.
func NewConsoleWriter(options ...func(w *ConsoleWriter)) ConsoleWriter {
	w := ConsoleWriter{
		out:        os.Stdout,
		timeFormat: consoleDefaultTimeFormat,
		partsOrder: consoleDefaultPartsOrder(),
	}

	for _, opt := range options {
		opt(&w)
	}

	// Fix color on Windows
	if w.out == os.Stdout || w.out == os.Stderr {
		w.out = colorable.NewColorable(w.out.(*os.File))
	}

	return w
}

func (w ConsoleWriter) Print(a ...any) (err error) {
	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	_, err = fmt.Fprint(buf, a...)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(w.out)
	return err
}

func (w ConsoleWriter) Printf(format string, a ...any) (err error) {
	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	_, err = fmt.Fprintf(buf, format, a...)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(w.out)
	return err
}

func (w ConsoleWriter) Println(a ...any) (err error) {
	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	_, err = fmt.Fprintln(buf, a...)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(w.out)
	return err
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (w ConsoleWriter) Write(p []byte) (err error) {
	var msg map[string]any
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&msg)
	if err != nil {
		return w.Printf("%s", string(p))
	}
	return w.writeMap(msg)
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (w ConsoleWriter) WriteAny(a any) (err error) {
	if a == nil {
		return
	}
	if m, ok := a.(map[string]any); ok {
		return w.writeMap(m)
	}
	return w.Print(a)
}

func (w ConsoleWriter) writeMap(a map[string]any) (err error) {
	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	for _, p := range w.partsOrder {
		w.writePart(buf, a, p)
	}

	w.writeFields(a, buf)
	_, err = buf.WriteTo(w.out)
	return err
}

// writeFields appends formatted key-value pairs to buf.
func (w ConsoleWriter) writeFields(evt map[string]any, buf *bytes.Buffer) {
	var fields = make([]string, 0, len(evt))
	for field := range evt {
		var isExcluded bool
		for _, excluded := range w.fieldsExclude {
			if field == excluded {
				isExcluded = true
				break
			}
		}
		if isExcluded {
			continue
		}

		switch field {
		case FieldNameLog,
			FieldNameLevel,
			FieldNameTime,
			FieldNameTimestamp,
			FieldNameAtTimestamp,
			FieldNameMessage,
			FieldNameCaller:
			continue
		}
		fields = append(fields, field)
	}
	sort.Strings(fields)

	// Write space only if something has already been written to the buffer, and if there are fields.
	if buf.Len() > 0 && len(fields) > 0 {
		buf.WriteByte(' ')
	}

	// Move the "error" field to the front
	ei := sort.Search(len(fields), func(i int) bool {
		return (fields[i] >= FieldNameError || fields[i] == FieldNameErr)
	})
	if ei < len(fields) && (fields[ei] == FieldNameError || fields[ei] == FieldNameErr) {
		field := fields[ei]
		fields[ei] = ""
		fields = append([]string{field}, fields...)
		var xfields = make([]string, 0, len(fields))
		for _, field := range fields {
			if field == "" { // Skip empty fields
				continue
			}
			xfields = append(xfields, field)
		}
		fields = xfields
	}

	for i, field := range fields {
		var fn Formatter
		var fv Formatter

		if field == FieldNameError ||
			field == FieldNameErr ||
			field == FieldNameStack {
			if w.formatErrFieldName == nil {
				fn = consoleDefaultFormatErrFieldName(w.noColor)
			} else {
				fn = w.formatErrFieldName
			}

			if w.formatErrFieldValue == nil {
				fv = consoleDefaultFormatErrFieldValue(w.noColor)
			} else {
				fv = w.formatErrFieldValue
			}
		} else {
			if w.formatFieldName == nil {
				fn = consoleDefaultFormatFieldName(w.noColor)
			} else {
				fn = w.formatFieldName
			}

			if w.formatFieldValue == nil {
				fv = consoleDefaultFormatFieldValue
			} else {
				fv = w.formatFieldValue
			}
		}

		buf.WriteString(fn(field))

		switch fValue := evt[field].(type) {
		case string:
			if needsQuote(fValue) {
				buf.WriteString(fv(strconv.Quote(fValue)))
			} else {
				buf.WriteString(fv(fValue))
			}
		case json.Number:
			buf.WriteString(fv(fValue))
		default:
			b, err := json.Marshal(fValue)
			if err != nil {
				fmt.Fprintf(buf, colorize("[error: %v]", colorRed, w.noColor), err)
			} else {
				fmt.Fprint(buf, fv(b))
			}
		}

		if i < len(fields)-1 { // Skip space for last field
			buf.WriteByte(' ')
		}
	}
}

// writePart appends a formatted part to buf.
func (w ConsoleWriter) writePart(buf *bytes.Buffer, evt map[string]any, p string) {
	var f Formatter

	if w.partsExclude != nil && len(w.partsExclude) > 0 {
		for _, exclude := range w.partsExclude {
			if exclude == p {
				return
			}
		}
	}

	switch p {
	case FieldNameLevel:
		if w.formatLevel == nil {
			f = consoleDefaultFormatLevel(w.noColor)
		} else {
			f = w.formatLevel
		}
	case FieldNameLog:
		if w.formatLog == nil {
			f = consoleDefaultFormatLog(w.noColor)
		} else {
			f = w.formatLog
		}
	case FieldNameTime, FieldNameTimestamp, FieldNameAtTimestamp:
		if w.formatTimestamp == nil {
			f = consoleDefaultFormatTimestamp(w.timeFormat, w.noColor)
		} else {
			f = w.formatTimestamp
		}
	case FieldNameMessage:
		if w.formatMessage == nil {
			f = consoleDefaultFormatMessage
		} else {
			f = w.formatMessage
		}
	case FieldNameCaller:
		if w.formatCaller == nil {
			f = consoleDefaultFormatCaller(w.noColor)
		} else {
			f = w.formatCaller
		}
	default:
		if w.formatFieldValue == nil {
			f = consoleDefaultFormatFieldValue
		} else {
			f = w.formatFieldValue
		}
	}

	v, ok := evt[p]
	if !ok {
		return
	}
	var s = f(v)

	if len(s) > 0 {
		if buf.Len() > 0 {
			buf.WriteByte(' ') // Write space only if not the first part
		}
		buf.WriteString(s)
	}
}

// needsQuote returns true when the string s should be quoted in output.
func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s any, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// ----- DEFAULT FORMATTERS ---------------------------------------------------

func consoleDefaultPartsOrder() []string {
	return []string{
		FieldNameTime,
		FieldNameTimestamp,
		FieldNameAtTimestamp,
		FieldNameLevel,
		FieldNameLog,
		FieldNameCaller,
		FieldNameMessage,
	}
}

func consoleDefaultFormatTimestamp(timeFormat string, noColor bool) Formatter {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	return func(i any) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.Parse(TimeFieldFormat, tt)
			if err != nil {
				t = tt
			} else {
				t = ts.Format(timeFormat)
			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64 = i, 0
				switch TimeFieldFormat {
				case TimeFormatUnixMs:
					nsec = int64(time.Duration(i) * time.Millisecond)
					sec = 0
				case TimeFormatUnixMicro:
					nsec = int64(time.Duration(i) * time.Microsecond)
					sec = 0
				}
				ts := time.Unix(sec, nsec).UTC()
				t = ts.Format(timeFormat)
			}
		}
		return colorize(t, colorDarkGray, noColor)
	}
}

func consoleDefaultFormatLevel(noColor bool) Formatter {
	return func(i any) string {
		if i == nil {
			return ""
		}
		if ll, ok := formatStrAsLevelString(i, noColor); ok {
			return ll
		}
		if ll, ok := formatNumAsLevelString(i, noColor); ok {
			return ll
		}
		ll := strings.ToUpper(fmt.Sprintf("%s", i))
		if len(ll) > 3 {
			ll = ll[0:3]
		}

		return ll
	}
}

func consoleDefaultFormatLog(noColor bool) Formatter {
	return func(i any) string {
		if i == nil {
			return consoleDefaultFormatMessage(i)
		}
		obj, ok := i.(map[string]any)
		if !ok {
			return consoleDefaultFormatMessage(i)
		}
		l, ok := obj["level"]
		if !ok {
			return consoleDefaultFormatMessage(i)
		}
		if ll, ok := formatStrAsLevelString(l, noColor); ok {
			return ll
		}
		if ll, ok := formatNumAsLevelString(l, noColor); ok {
			return ll
		}
		ll := strings.ToUpper(fmt.Sprintf("%s", l))
		if len(ll) > 3 {
			ll = ll[0:3]
		}

		return ll
	}
}

func formatStrAsLevelString(i any, noColor bool) (string, bool) {
	ll, ok := i.(string)
	if !ok {
		return "", false
	}
	switch ll {
	case LevelPanicValue:
		return colorize(colorize("PNC", colorRed, noColor), colorBold, noColor), true
	case LevelFatalValue:
		return colorize(colorize("FTL", colorRed, noColor), colorBold, noColor), true
	case LevelErrorValue:
		return colorize(colorize("ERR", colorRed, noColor), colorBold, noColor), true
	case LevelWarnValue:
		return colorize("WRN", colorRed, noColor), true
	case LevelInfoValue:
		return colorize("INF", colorGreen, noColor), true
	case LevelDebugValue:
		return colorize("DBG", colorYellow, noColor), true
	case LevelTraceValue:
		return colorize("TRC", colorMagenta, noColor), true
	default:
		return colorize("???", colorBold, noColor), true
	}
}

func formatNumAsLevelString(i any, noColor bool) (string, bool) {
	num, ok := i.(json.Number)
	if !ok {
		return "", false
	}
	ll, err := num.Int64()
	if err != nil {
		return "", false
	}
	if ll >= LevelPanicNum {
		return colorize(colorize("PNC", colorRed, noColor), colorBold, noColor), true
	}
	if ll >= LevelFatalNum {
		return colorize(colorize("FTL", colorRed, noColor), colorBold, noColor), true
	}
	if ll >= LevelErrorNum {
		return colorize(colorize("ERR", colorRed, noColor), colorBold, noColor), true
	}
	if ll >= LevelWarnNum {
		return colorize("WRN", colorRed, noColor), true
	}
	if ll >= LevelInfoNum {
		return colorize("INF", colorGreen, noColor), true
	}
	if ll >= LevelDebugNum {
		return colorize("DBG", colorYellow, noColor), true
	}
	return colorize("TRC", colorMagenta, noColor), true
}

func consoleDefaultFormatCaller(noColor bool) Formatter {
	return func(i any) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
			c = colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
		}
		return c
	}
}

func consoleDefaultFormatMessage(i any) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%s", i)
}

func consoleDefaultFormatFieldName(noColor bool) Formatter {
	return func(i any) string {
		return colorize(fmt.Sprintf("%s=", i), colorCyan, noColor)
	}
}

func consoleDefaultFormatFieldValue(i any) string {
	return fmt.Sprintf("%s", i)
}

func consoleDefaultFormatErrFieldName(noColor bool) Formatter {
	return func(i any) string {
		return colorize(fmt.Sprintf("%s=", i), colorCyan, noColor)
	}
}

func consoleDefaultFormatErrFieldValue(noColor bool) Formatter {
	return func(i any) string {
		str, err := strconv.Unquote(fmt.Sprintf("%s", i))
		if err != nil {
			return colorize(fmt.Sprintf("%s", i), colorRed, noColor)
		}
		return colorize(str, colorRed, noColor)
	}
}
