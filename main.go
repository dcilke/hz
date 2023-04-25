package main

import (
	"fmt"
	"os"

	"github.com/dcilke/hz/pkg/processor"
	"github.com/dcilke/hz/pkg/terminator"
	"github.com/dcilke/hz/pkg/writer"
	"github.com/jessevdk/go-flags"
)

type Cmd struct {
	Strict bool `short:"s" long:"strict" description:"strict mode"`
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

	p := processor.New(writer.New(), processor.WithStrict(cmd.Strict))

	terminator.OnSig(func() int {
		p.Flush()
		return 0
	})

	if len(filenames) > 0 {
		for _, filename := range filenames {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprint(os.Stderr, fmt.Errorf("unable to open %q: %w", filename, err))
			}
			p.Process(f)
		}
		return
	}

	p.Process(os.Stdin)
}
