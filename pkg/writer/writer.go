package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"sync"

	"github.com/dcilke/gu"
	"github.com/dcilke/hz/pkg/formatter"
	"github.com/mattn/go-colorable"
)

const (
	PinTimestamp = "timestamp"
	PinLevel     = "level"
	PinCaller    = "caller"
	PinMessage   = "message"
	PinError     = "error"

	defaultTimeFormat = "15:04:05"
	defaultSep        = ' '
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

	// includeLevels defines the log levels to include in output.
	includeLevels []string

	// excludeKeys defines contextual keys to not display in output.
	excludeKeys []string

	// formatter defines a map of formatters for pins.
	formatter map[string]formatter.Formatter

	// fielder defines the default field formatter.
	fielder formatter.Fielder

	// formatKey defines the default key formatter.
	formatKey formatter.Stringer

	// flatten enables flattening of JSON objects.
	flatten bool
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

func WithFormatter(key string, f formatter.Formatter) Option {
	return func(w *Writer) {
		w.formatter[key] = f
	}
}

func WithPinOrder(order []string) Option {
	return func(w *Writer) {
		w.pinOrder = order
	}
}

func WithFielder(f formatter.Fielder) Option {
	return func(w *Writer) {
		w.fielder = f
	}
}

func WithExcludeKeys(keys []string) Option {
	return func(w *Writer) {
		w.excludeKeys = keys
	}
}

func WithKeyFormatter(f formatter.Stringer) Option {
	return func(w *Writer) {
		w.formatKey = f
	}
}

func WithLevelFilter(s string) Option {
	return func(w *Writer) {
		w.includeLevels = append(w.includeLevels, s)
	}
}

func WithLevelFilters(s []string) Option {
	return func(w *Writer) {
		w.includeLevels = append(w.includeLevels, s...)
	}
}

func WithFlatten(b bool) Option {
	return func(w *Writer) {
		w.flatten = b
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
		formatter:   make(map[string]formatter.Formatter, 6),
		flatten:     false,
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
		w.formatKey = formatter.Key(w.noColor)
	}

	// Ensure default formatters, if not specified in input
	if _, ok := w.formatter[PinTimestamp]; !ok {
		w.formatter[PinTimestamp] = formatter.NewTimestamp(w.noColor, w.formatKey, w.timeFormat)
	}
	if _, ok := w.formatter[PinLevel]; !ok {
		w.formatter[PinLevel] = formatter.NewLevel(w.noColor, w.formatKey)
	}
	if _, ok := w.formatter[PinMessage]; !ok {
		w.formatter[PinMessage] = formatter.NewMessage(w.noColor, w.formatKey)
	}
	if _, ok := w.formatter[PinCaller]; !ok {
		w.formatter[PinCaller] = formatter.NewCaller(w.noColor)
	}
	if _, ok := w.formatter[PinError]; !ok {
		w.formatter[PinError] = formatter.NewError(w.noColor, w.formatKey)
	}

	// Ensure default extractor
	if w.fielder == nil {
		w.fielder = formatter.Map(w.formatKey)
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

func (w Writer) Println(a ...any) (int, error) {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if b, err := fmt.Fprintln(buf, a...); err != nil {
		return b, err
	}
	b, err := buf.WriteTo(w.out)
	return int(b), err
}

// WriteBytes transforms the JSON input with formatters and appends to w.Out.
func (w Writer) Write(p []byte) (int, error) {
	var msg any
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	if err := d.Decode(&msg); err != nil {
		return w.Print(string(p))
	}
	return w.WriteAny(msg)
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
		return w.Print(b)
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Array, reflect.Slice:
		return w.writeArray(a.([]any))
	}
	return w.Print(a)
}

func (w Writer) writeMap(a map[string]any) (int, error) {
	var buf = bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if len(w.includeLevels) > 0 {
		levels := formatter.GetLevels(a)
		for _, l := range levels {
			if !gu.Includes(w.includeLevels, l) {
				return 0, nil
			}
		}
	}

	for _, p := range w.pinOrder {
		w.writePinned(buf, a, p)
	}

	// Write space only if something has already been written to the buffer and we are going to write
	// a key which was not pinned
	if buf.Len() > 0 {
		for key := range a {
			if !gu.Includes(w.excludeKeys, key) {
				buf.WriteByte(defaultSep)
				break
			}
		}
	}

	w.writeFields(buf, a, "")
	b, err := buf.WriteTo(w.out)
	return int(b), err
}

func (w Writer) writeArray(a []any) (int, error) {
	b, err := w.Print("[\n")
	if err != nil {
		return b, err
	}
	for _, v := range a {
		c, err := w.WriteAny(v)
		b += c
		if err != nil {
			return b, err
		}
		if c > 0 {
			c, err = w.Print("\n")
			b += c
			if err != nil {
				return b, err
			}
		}
	}
	c, err := w.Print("]")
	b += c
	return b, err
}

// writeFields appends formatted key-value pairs to buf.
func (w Writer) writeFields(buf *bytes.Buffer, evt map[string]any, prefix string) {
	var keys = make([]string, 0, len(evt))
	for key := range evt {
		if gu.Includes(w.excludeKeys, prefix+key) {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		value := evt[key]
		if m, ok := value.(map[string]any); ok && w.flatten {
			w.writeFields(buf, m, prefix+key+".")
		} else {
			buf.WriteString(w.fielder(prefix+key, value))
		}
		// Skip space for last key
		if i < len(keys)-1 {
			buf.WriteByte(defaultSep)
		}
	}
}

// writePinned appends a formatted part to buf.
func (w Writer) writePinned(buf *bytes.Buffer, evt map[string]any, p string) {
	var s string
	if f, ok := w.formatter[p]; ok {
		s = f.Format(evt)
	} else {
		s = w.fielder(p, evt[p])
	}

	if len(s) > 0 {
		// Write space only if not the first part
		if buf.Len() > 0 {
			buf.WriteByte(defaultSep)
		}
		buf.WriteString(s)
	}
}
