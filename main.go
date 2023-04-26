package main

import (
	"fmt"
	"os"

	"github.com/dcilke/hz/pkg/emitter"
	"github.com/dcilke/hz/pkg/terminator"
	"github.com/dcilke/hz/pkg/writer"
	"github.com/jessevdk/go-flags"
)

type Cmd struct {
	Strict bool     `short:"s" long:"strict" description:"strict mode"`
	Level  []string `short:"l" long:"level" description:"only output lines at this level"`
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

	bufSize := emitter.DefaultBufSize
	if cmd.Strict {
		bufSize = 0
	}

	w := writer.New(
		writer.WithLevelFilters(cmd.Level),
	)
	e := emitter.New(
		emitter.WithBufSize(bufSize),
		emitter.WithJSON(func(a any) {
			s, _ := w.WriteAny(a)
			if s > 0 {
				w.Println()
			}
		}),
		emitter.WithBytes(func(b []byte) {
			s, _ := w.Print(string(b))
			if s > 0 {
				w.Println()
			}
		}),
		emitter.WithError(func(err error) {
			fmt.Fprint(os.Stderr, fmt.Errorf("extractor error: %w", err))
		}),
	)

	terminator.OnSig(func() int {
		e.Flush()
		return 0
	})

	if len(filenames) > 0 {
		for _, filename := range filenames {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprint(os.Stderr, fmt.Errorf("unable to open %q: %w", filename, err))
			}
			e.Process(f)
		}
		return
	}

	e.Process(os.Stdin)
}
