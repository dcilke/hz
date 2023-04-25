package processor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/dcilke/goj"
)

const (
	defaultBufSize = 512

	newline = 10
)

type writer interface {
	Error(v error) (int, error)
	Printf(format string, v ...interface{}) (int, error)
	WriteAny(any any) (int, error)
}

type Processor struct {
	writer  writer
	buf     *bytes.Buffer
	bufSize int
}

type Option func(*Processor)

func WithBufSize(size int) Option {
	return func(p *Processor) {
		p.buf = bytes.NewBuffer(make([]byte, 0, size))
	}
}

func New(writer writer, opts ...Option) *Processor {
	p := &Processor{
		writer:  writer,
		bufSize: defaultBufSize,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.buf = bytes.NewBuffer(make([]byte, 0, p.bufSize))

	return p
}

func (p *Processor) Process(file *os.File) {
	reader := bufio.NewReader(file)
	decoder := goj.NewDecoder(reader)
	decoder.UseNumber()

	for {
		// we will only pretty print objects, so if the current character is not { or [, process as byte stream
		if c, _ := decoder.Peek(); goj.IsBegin(c) {
			var msg any
			_ = decoder.Decode(&msg)
			t := reflect.TypeOf(msg)
			if t != nil {
				switch t.Kind() {
				case reflect.Map, reflect.Array, reflect.Slice:
					p.Flush()
					_, err := p.writer.WriteAny(msg)
					if err != nil {
						_, _ = p.writer.Error(err)
					}
					continue
				case reflect.String: // json.Number
					for _, b := range []byte(msg.(json.Number)) {
						p.push(b)
					}
					continue
				default:
					_, _ = p.writer.Error(fmt.Errorf("unexpected decoded line type %q", t.Kind()))
				}
			}
		}

		b, err := decoder.ReadByte()
		if err == io.EOF {
			p.Flush()
			return
		}

		p.push(b)
		if p.buf.Len() >= p.bufSize || b == newline {
			p.Flush()
		}
	}
}

func (p *Processor) push(b byte) {
	p.buf.WriteByte(b)
}

func (p *Processor) Flush() {
	if p.buf.Len() > 0 {
		p.writer.Printf("%s", p.buf.String())
		p.buf.Reset()
	}
}
