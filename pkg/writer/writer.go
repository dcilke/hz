package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/mattn/go-colorable"
)

const (
	PinTimestamp = "timestamp"
	PinLevel     = "level"
	PinCaller    = "caller"
	PinMessage   = "message"
	PinError     = "error"

	defaultTimeFormat = "15:04:05"
)

// Ensure we are adhering to the io.Writer interface.
var _ io.Writer = (*Writer)(nil)

var (
	defaultPinOrder = []string{
		PinTimestamp,
		PinLevel,
		PinCaller,
		PinMessage,
		PinError,
	}

	bufPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

// Formatter defines a formatter for a specific pin key
type Formatter interface {
	Format(map[string]any, string) string
	ExcludeKeys() []string
}

// Extractor extracts multiple values and formats them
type Extractor func(map[string]any, string) string

// Stringer stringifies a value
type Stringer func(any) string

// Writer parses the JSON input and writes it in an
// (optionally) colorized, human-friendly format to Out.
type Writer struct {
	// out is the output destination.
	out io.Writer

	// err is the error destination.
	err io.Writer

	// noColor disables the colorized output.
	noColor bool

	// timeFormat specifies the format for timestamp in output.
	timeFormat string

	// pinOrder defines the order of set keys in output.
	pinOrder []string

	// excludeKeys defines contextual keys to not display in output.
	excludeKeys []string

	formatter map[string]Formatter

	extractor Extractor

	formatKey Stringer
}

type Option func(w *Writer)

// Override the output writer, defaults to os.Stdout.
func WithOut(out io.Writer) Option {
	return func(w *Writer) {
		w.out = out
	}
}

// Override the error writer, defaults to os.Stderr.
func WithErr(err io.Writer) Option {
	return func(w *Writer) {
		w.err = err
	}
}

// Disable colorized output.
func WithNoColor() Option {
	return func(w *Writer) {
		w.noColor = true
	}
}

// Override the time format, defaults to "15:04:05".
func WithTimeFormat(timeFormat string) Option {
	return func(w *Writer) {
		w.timeFormat = timeFormat
	}
}

func WithFormatter(key string, f Formatter) Option {
	return func(w *Writer) {
		w.formatter[key] = f
	}
}

func WithPinOrder(order []string) Option {
	return func(w *Writer) {
		w.pinOrder = order
	}
}

func WithExtractor(f Extractor) Option {
	return func(w *Writer) {
		w.extractor = f
	}
}

func WithExcludeKeys(keys []string) Option {
	return func(w *Writer) {
		w.excludeKeys = keys
	}
}

func WithKeyFormatter(f Stringer) Option {
	return func(w *Writer) {
		w.formatKey = f
	}
}

// New creates and initializes a new ConsoleWriter.
func New(options ...Option) Writer {
	w := Writer{
		out:         os.Stdout,
		err:         os.Stderr,
		timeFormat:  defaultTimeFormat,
		pinOrder:    defaultPinOrder,
		excludeKeys: make([]string, 0, 10),
		formatter:   make(map[string]Formatter, 6),
	}

	for _, opt := range options {
		opt(&w)
	}

	// Fix color on Windows
	if w.out == os.Stdout || w.out == os.Stderr {
		w.out = colorable.NewColorable(w.out.(*os.File))
	}

	if w.err == os.Stdout || w.err == os.Stderr {
		w.err = colorable.NewColorable(w.err.(*os.File))
	}

	// Set default key formatter
	if w.formatKey == nil {
		w.formatKey = formatKey(w.noColor)
	}

	// Ensure default formatters, if not specified in input
	if _, ok := w.formatter[PinTimestamp]; !ok {
		w.formatter[PinTimestamp] = newTimestampFormatter(w.noColor, w.formatKey, w.timeFormat)
	}
	if _, ok := w.formatter[PinLevel]; !ok {
		w.formatter[PinLevel] = newLevelFormatter(w.noColor, w.formatKey)
	}
	if _, ok := w.formatter[PinMessage]; !ok {
		w.formatter[PinMessage] = newMessageFormatter(w.noColor, w.formatKey)
	}
	if _, ok := w.formatter[PinCaller]; !ok {
		w.formatter[PinCaller] = newCallerFormatter(w.noColor, w.formatKey)
	}
	if _, ok := w.formatter[PinError]; !ok {
		w.formatter[PinError] = newErrorFormatter(w.noColor, w.formatKey)
	}

	// Ensure default extractor
	if w.extractor == nil {
		w.extractor = extractor(w.noColor, w.formatKey)
	}

	for _, v := range w.formatter {
		w.excludeKeys = append(w.excludeKeys, v.ExcludeKeys()...)
	}

	return w
}

func (w Writer) Print(a ...any) (int, error) {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if b, err := fmt.Fprint(buf, a...); err != nil {
		return b, err
	}
	b, err := buf.WriteTo(w.out)
	return int(b), err
}

func (w Writer) Printf(format string, a ...any) (int, error) {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if b, err := fmt.Fprintf(buf, format, a...); err != nil {
		return b, err
	}
	b, err := buf.WriteTo(w.out)
	return int(b), err
}

func (w Writer) Error(e error) (int, error) {
	return fmt.Fprint(w.err, e)
}

// WriteBytes transforms the JSON input with formatters and appends to w.Out.
func (w Writer) Write(p []byte) (int, error) {
	var msg map[string]any
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	if err := d.Decode(&msg); err != nil {
		return w.Printf("%s", string(p))
	}
	return w.writeMap(msg)
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (w Writer) WriteAny(a any) (int, error) {
	if a == nil {
		return 0, nil
	}
	if m, ok := a.(map[string]any); ok {
		return w.writeMap(m)
	}
	if b, ok := a.([]byte); ok {
		return w.Write(b)
	}
	return w.Print(a)
}

func (w Writer) writeMap(a map[string]any) (int, error) {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	for _, p := range w.pinOrder {
		w.writePinned(buf, a, p)
	}

	w.writeFields(buf, a)
	b, err := buf.WriteTo(w.out)
	return int(b), err
}

// writeFields appends formatted key-value pairs to buf.
func (w Writer) writeFields(buf *bytes.Buffer, evt map[string]any) {
	var keys = make([]string, 0, len(evt))
	for key := range evt {
		if includes(w.excludeKeys, key) {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write space only if something has already been written to the buffer, and if there are keys.
	if buf.Len() > 0 && len(keys) > 0 {
		buf.WriteByte(' ')
	}

	for i, key := range keys {
		buf.WriteString(w.extractor(evt, key))
		if i < len(keys)-1 { // Skip space for last key
			buf.WriteByte(' ')
		}
	}
}

// writePinned appends a formatted part to buf.
func (w Writer) writePinned(buf *bytes.Buffer, evt map[string]any, p string) {
	var s string
	if f, ok := w.formatter[p]; ok {
		s = f.Format(evt, p)
	} else {
		s = w.extractor(evt, p)
	}

	if len(s) > 0 {
		if buf.Len() > 0 {
			buf.WriteByte(' ') // Write space only if not the first part
		}
		buf.WriteString(s)
	}
}

func formatKey(noColor bool) Stringer {
	return func(i any) string {
		return colorize(fmt.Sprintf("%s=", i), ColorCyan, noColor)
	}
}
func extractor(noColor bool, fn Stringer) Extractor {
	return func(m map[string]any, k string) string {
		ret := fn(k)
		switch fValue := m[k].(type) {
		case string:
			if needsQuote(fValue) {
				ret += strconv.Quote(fValue)
			} else {
				ret += fValue
			}
		case json.Number:
			ret += fValue.String()
		default:
			b, err := json.Marshal(fValue)
			if err != nil {
				ret += fmt.Sprintf(colorize("[error: %v]", ColorRed, noColor), err)
			} else {
				ret += string(b)
			}
		}
		return ret
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
