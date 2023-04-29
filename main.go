package main

import (
	"fmt"
	"os"

	"github.com/dcilke/gu"
	"github.com/dcilke/heron"
	"github.com/dcilke/hz/pkg/writer"
	"github.com/jessevdk/go-flags"
)

const (
	newline = "\n"
)

type Cmd struct {
	Strict bool     `short:"s" long:"strict" description:"strict mode"`
	Level  []string `short:"l" long:"level" description:"only output lines at this level"`
	Flat   bool     `short:"f" long:"flat" description:"flatten output"`
	Vert   bool     `short:"v" long:"vert" description:"vertical output"`
}

func main() {
	var cmd Cmd
	parser := flags.NewParser(&cmd, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "[FILE]"
	filenames, err := parser.Parse()
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		parser.WriteHelp(os.Stdout)
		return
	}
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("unable to parse arguments: %w", err))
	}

	bufSize := heron.DefaultBufSize
	if cmd.Strict {
		bufSize = 0
	}

	w := writer.New(
		writer.WithLevelFilters(cmd.Level),
		writer.WithFlatten(cmd.Flat),
		writer.WithVertical(cmd.Vert),
	)
	// didnl is used to prevent double newlines since we want to ensure each JSON
	// objects is on its own line but we want to preserve as much of the output
	// as possible
	didnl := true
	h := heron.New(
		heron.WithBufSize(bufSize),
		heron.WithJSON(func(a any) {
			s, _ := w.WriteAny(a)
			if s > 0 {
				didnl = true
				w.Println()
			}
		}),
		heron.WithBytes(func(b []byte) {
			sb := string(b)
			if sb == newline && didnl {
				return
			}
			didnl = false
			_, _ = w.Print(sb)
		}),
		heron.WithError(func(err error) {
			fmt.Fprint(os.Stderr, fmt.Errorf("extractor error: %w", err))
		}),
	)

	gu.Terminator(func() int {
		h.Flush()
		return 0
	})

	if len(filenames) > 0 {
		for _, filename := range filenames {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprint(os.Stderr, fmt.Errorf("unable to open %q: %w", filename, err))
			}
			h.Process(f)
		}
		return
	}

	h.Process(os.Stdin)
}
