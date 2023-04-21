package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/dcilke/goj"
	"github.com/jessevdk/go-flags"
)

var writer = NewConsoleWriter()

const (
	bufDump = 512
)

type cmd struct{}

func main() {
	var c cmd
	parser := flags.NewParser(&c, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "[FILE]"
	filenames, err := parser.Parse()
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		parser.WriteHelp(os.Stdout)
		return
	}
	check(err, "unable to parse arguments")

	if len(filenames) > 0 {
		for _, filename := range filenames {
			f, err := os.Open(filename)
			check(err, fmt.Sprintf("unable to open %q", filename))
			process(f)
		}
		return
	}

	process(os.Stdin)
}

func process(file *os.File) {
	reader := bufio.NewReader(file)
	decoder := goj.NewDecoder(reader)
	decoder.UseNumber()
	_buf := make([]byte, 0, bufDump)
	buf := bytes.NewBuffer(_buf)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		flush(buf)
		os.Exit(0)
	}()

	for {
		// we will only pretty print objects, so if the current character is not { or [, process as byte stream
		if c, _ := decoder.Peek(); goj.IsBegin(c) {
			var msg any
			_ = decoder.Decode(&msg)
			t := reflect.TypeOf(msg)
			if t != nil {
				switch t.Kind() {
				case reflect.Map, reflect.Array, reflect.Slice:
					flush(buf)
					err := writer.WriteAny(msg)
					if err != nil {
						writer.Println(err)
					}
					continue
				case reflect.String: // json.Number
					for _, b := range []byte(msg.(json.Number)) {
						buf.WriteByte(b)
					}
					continue
				default:
					writer.Println(fmt.Sprintf("[hz] internal error, unexpected decoded line type %q", t.Kind()))
				}
			}
		}

		b, err := decoder.ReadByte()
		if err == io.EOF {
			flush(buf)
			return
		}
		buf.WriteByte(b)
		if buf.Len() >= bufDump || b == 10 {
			flush(buf)
		}
	}
}

func flush(buf *bytes.Buffer) {
	if buf.Len() > 0 {
		writer.Printf("%s", buf.String())
		buf.Reset()
	}
}

func check(err error, hint string) {
	if err != nil {
		writer.Println(hint)
		writer.Println(err)
	}
}
