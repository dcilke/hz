package emitter

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/dcilke/goj"
)

const (
	DefaultBufSize = 512

	newLine = 10
)

type Emitter struct {
	buf     *bytes.Buffer
	bufSize int
	json    func(any)
	bytes   func([]byte)
	error   func(error)
}

type Option func(*Emitter)

func WithBufSize(size int) Option {
	return func(p *Emitter) {
		p.bufSize = size
	}
}

func WithJSON(f func(any)) Option {
	return func(p *Emitter) {
		p.json = f
	}
}

func WithBytes(f func([]byte)) Option {
	return func(p *Emitter) {
		p.bytes = f
	}
}

func WithError(f func(error)) Option {
	return func(p *Emitter) {
		p.error = f
	}
}

func New(opts ...Option) *Emitter {
	p := &Emitter{
		bufSize: DefaultBufSize,
		json:    func(a any) {},
		bytes:   func(s []byte) {},
		error:   func(e error) {},
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.bufSize > 0 {
		p.buf = bytes.NewBuffer(make([]byte, 0, p.bufSize))
	}

	return p
}

func (p *Emitter) Process(stream io.Reader) {
	decoder := goj.NewDecoder(stream)
	decoder.UseNumber()

	for {
		// we will only emit objects or arrays, so if the current character is not { or [, process as byte stream
		if c, _ := decoder.Peek(); goj.IsBegin(c) {
			var msg any
			err := decoder.Decode(&msg)
			if err == nil {
				switch k := reflect.TypeOf(msg).Kind(); k {
				case reflect.Map, reflect.Array, reflect.Slice:
					p.Flush()
					p.json(msg)
				default:
					p.error(fmt.Errorf("unexpected decoded line type %q", k))
				}
				continue
			}
		}

		if err := p.processByte(decoder); err == io.EOF {
			return
		}
	}
}

func (p *Emitter) BufSize() int {
	return p.bufSize
}

func (p *Emitter) processByte(decoder *goj.Decoder) error {
	b, err := decoder.ReadByte()
	if err == io.EOF || b == newLine {
		p.Flush()
		return err
	}

	p.push(b)
	p.flushFull()
	return nil
}

func (p *Emitter) push(b byte) {
	if p.buf == nil {
		return
	}
	p.buf.WriteByte(b)
}

func (p *Emitter) flushFull() {
	if p.buf == nil {
		return
	}
	if p.buf.Len() >= p.bufSize {
		p.Flush()
	}
}

func (p *Emitter) Flush() {
	if p.buf == nil {
		return
	}
	if p.buf.Len() > 0 {
		p.bytes(p.buf.Bytes())
		p.buf.Reset()
	}
}
