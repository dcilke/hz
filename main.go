package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var writer = NewConsoleWriter()

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

	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		_, err = writer.Write(b)
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
