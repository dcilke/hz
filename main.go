package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dcilke/gu"
	"github.com/dcilke/heron"
	"github.com/dcilke/hz/pkg/writer"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

const (
	newline = "\n"
)

var cfgPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "." // fallback to current directory
	}
	cfgPath = filepath.Join(home, ".config", "hz", "config.yml")
}

type Cmd struct {
	Level    []string `short:"l" long:"level" description:"only output lines at this level" yaml:"level"`
	Strict   bool     `short:"s" long:"strict" description:"exclude non JSON output" yaml:"strict"`
	Flat     bool     `short:"f" long:"flat" description:"flatten objects" yaml:"flat"`
	Vertical bool     `short:"v" long:"vertical" description:"vertical output" yaml:"vertical"`
	Raw      bool     `short:"r" long:"raw" description:"raw output" yaml:"plain"`
	NoPin    bool     `short:"n" long:"no-pin" description:"exclude pinning of fields" yaml:"noPin"`
}

func main() {
	var cmd Cmd
	// read in config file, if it exists
	err := loadDefaults(&cmd)
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("WARN: unable to load config: %w", err), "\n")
	}

	// parse command line flags
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

	opts := []writer.Option{
		writer.WithLevelFilters(cmd.Level),
		writer.WithFlatten(cmd.Flat),
		writer.WithVertical(cmd.Vertical),
		writer.WithColor(!cmd.Raw),
	}

	if cmd.NoPin {
		opts = append(opts, writer.WithPinOrder([]string{}))
	}

	w := writer.New(opts...)
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

func loadDefaults(cfg *Cmd) error {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil
	}

	bytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("unable to read config file: %w", err)
	}
	return yaml.Unmarshal(bytes, &cfg)
}
