package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"voidedtech.com/gxs/internal"
)

func die(message string, err error) {
	fmt.Fprintln(os.Stderr, message)
	fmt.Fprintln(os.Stderr, err)
	panic("failed")
}

func stdin() []byte {
	scanner := bufio.NewScanner(os.Stdin)
	var b bytes.Buffer
	for scanner.Scan() {
		b.WriteString(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		die("failed to read stdin", err)
	}
	return b.Bytes()
}

func main() {
	file := flag.String("input", "", "file to take as an input pattern (else stdin)")
	out := flag.String("output", "", "file to save output (else stdout)")
	flag.Parse()
	fileName := *file
	b := []byte("")
	if len(fileName) == 0 {
		die("no input file found", fmt.Errorf("input file required"))
	} else {
		if fileName == "--" {
			b = stdin()
		} else {
			raw, err := os.ReadFile(fileName)
			if err != nil {
				die("unable to read file", err)
			}
			b = raw
		}
	}
	pattern, pErr := internal.Parse(b)
	if pErr != nil && pErr.Error != nil {
		if pErr.Backtrace != nil {
			for _, line := range pErr.Backtrace {
				fmt.Fprintln(os.Stderr, line)
			}
		}
		die("unable to parse pattern", pErr.Error)
	}
	tmpl, err := internal.Build(pattern)
	if err != nil {
		die("failed to template", err)
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
		die("failed to write output", err)
	}
}
