package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var writer = zerolog.NewConsoleWriter()

type cmd struct{}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
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

	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		_, err = writer.Write(b)
		if err != nil {
			/**
			TODO: handle "pretty printed" objects
			"Better" json stream parser? The need is a parser which will handle
			the json and output the invalid characters while skipping over them.
			Read into buffer, classifying as valid or invalid json. Something
			with an api like:

			var obj interface{}
			var str string
			err := reader.Decode(&obj, &str)
			**/
			fmt.Print(string(b))
		}
	}
}

func check(err error, hint string) {
	if err != nil {
		log.Fatal().Err(err).Msg(hint)
	}
}
