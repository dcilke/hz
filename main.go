package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type cmd struct{}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

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
	decoder := json.NewDecoder(reader)
	writer := zerolog.NewConsoleWriter()
	for decoder.More() {
		var line interface{}
		err := decoder.Decode(&line)
		check(err, "unable to decode object")
		data, err := json.Marshal(line)
		check(err, "unable to marshal object")
		_, err = writer.Write(data)
		check(err, "unable to write object")
	}
}

func check(err error, hint string) {
	if err != nil {
		log.Fatal().Err(err).Msg(hint)
	}
}
