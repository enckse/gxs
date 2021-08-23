package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"voidedtech.com/gxs/internal"
	"voidedtech.com/stock"
)

var (
	version = "development"
)

func stdin() []byte {
	b, err := stock.Stdin(false)
	if err != nil {
		stock.Die("failed to read stdin", err)
	}
	return b
}

func main() {
	file := flag.String("input", "", "file to take as an input pattern (else stdin)")
	out := flag.String("output", "", "file to save output (else stdout)")
	outMode := flag.String("format", internal.ASCIIMode, "output format")
	showVers := flag.Bool("version", false, "display version")
	option := &internal.Option{}
	flag.Func("option", "gxs options", func(s string) error {
		return option.Set(s)
	})
	flag.Parse()
	if *showVers {
		fmt.Printf("version: %s\n", version)
		return
	}
	var b []byte
	if fileName := *file; fileName == "" {
		b = stdin()
	} else {
		raw, err := os.ReadFile(fileName)
		if err != nil {
			stock.Die("unable to read file", err)
		}
		b = raw
	}
	pattern, pErr := internal.Parse(b)
	if pErr != nil && pErr.Error != nil {
		if pErr.Backtrace != nil {
			for _, line := range pErr.Backtrace {
				fmt.Fprintln(os.Stderr, line)
			}
		}
		stock.Die("unable to parse pattern", pErr.Error)
	}
	tmpl, err := internal.Build(pattern, *outMode, option)
	if err != nil {
		stock.Die("failed to template", err)
	}
	var write io.Writer
	outFile := *out
	if len(outFile) == 0 {
		write = os.Stdout
	} else {
		var b bytes.Buffer
		write = &b
	}
	if _, err := write.Write(tmpl); err != nil {
		stock.Die("failed to write output", err)
	}
}
