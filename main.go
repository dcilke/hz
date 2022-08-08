package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dcilke/gojay"
	"github.com/jessevdk/go-flags"
)

var writer = NewConsoleWriter()

const (
	bufSize = 1024
	bufDump = bufSize / 2
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
	decoder := gojay.Stream.NewDecoder(reader)
	decoder.UseNumber()
	_buf := make([]byte, 0, bufSize)
	buf := bytes.NewBuffer(_buf)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		writer.Printf("%s", buf.String())
		buf.Reset()
		os.Exit(0)
	}()

	for {
		var msg any
		err := decoder.Decode(&msg)
		if err != nil {
			b, _ := decoder.ReadByte()
			buf.WriteByte(b)
			if buf.Len() > bufDump || b == 10 {
				writer.Printf("%s", buf.String())
				buf.Reset()
			}
			continue
		}
		if buf.Len() > 0 {
			writer.Printf("%s", buf.String())
			buf.Reset()
		}
		err = writer.WriteAny(msg)
		if err != nil {
			writer.Println(err)
		}
	}
}

func check(err error, hint string) {
	if err != nil {
		writer.Println(hint)
		writer.Println(err)
	}
}
