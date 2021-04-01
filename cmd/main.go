package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var (
	//go:embed "template.html"
	templateHTML string
)

func die(message string, err error) {
	fmt.Println(message)
	fmt.Println(err)
	panic("failed")
}

type (
	Cell struct {
		ID    string
		Value string
	}
	Pattern struct {
		Size        int
		Pad         string
		PadChar     string
		Interactive bool
		Cells       []Cell
		JSON        template.JS
	}
)

func (o Pattern) pad(val int) string {
	padded := fmt.Sprintf("%s%d", o.Pad, val)
	for len(padded) > len(o.Pad) {
		padded = padded[1:]
	}
	return padded
}

func (o Pattern) initCells() []Cell {
	var results []Cell
	x := 0
	for x < o.Size {
		y := 0
		for y < o.Size {
			left := o.pad(x)
			right := o.pad(y)
			val := ""
			if x == 0 {
				val = fmt.Sprintf("%d", y)
			}
			if y == 0 {
				val = fmt.Sprintf("%d", x)
			}
			if x == 0 && y == 0 {
				val = ""
			}
			cell := Cell{}
			cell.ID = fmt.Sprintf("%s%s%s", left, o.PadChar, right)
			cell.Value = val
			results = append(results, cell)
			y += 1
		}
		x += 1
	}
	return results
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
	size := flag.Int("size", 25, "pattern size")
	pad := flag.Int("pad", 4, "id padding (with 0s)")
	char := flag.String("padchar", "x", "padding character")
	bind := flag.String("bind", "", "local binding to use for server")
	file := flag.String("input", "", "file to take as an input pattern (-- stdin)")
	flag.Parse()
	fileName := *file
	b := []byte("[]")
	binding := *bind
	isCommandLine := binding == ""
	if len(fileName) == 0 {
		if isCommandLine {
			die("no input pattern and non-interactive editing", fmt.Errorf("bad configuration"))
		}
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
	tmpl, err := template.New("t").Parse(templateHTML)
	if err != nil {
		die("unable to parse template", err)
	}
	padding := *pad
	if padding < 1 {
		die("invalid padding", fmt.Errorf("< 1"))
	}
	padString := ""
	for padding > 0 {
		padString = fmt.Sprintf("0%s", padString)
		padding = padding - 1
	}
	obj := Pattern{Size: *size + 1, Pad: padString, PadChar: *char}
	obj.Cells = obj.initCells()
	obj.JSON = template.JS(string(b))
	if obj.Size <= 0 {
		die("invalid size", fmt.Errorf("< 0"))
	}
	obj.Interactive = false
	if isCommandLine {
		if err := tmpl.Execute(os.Stdout, obj); err != nil {
			die("unable to execute template", err)
		}
		return
	}
	obj.Interactive = true
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, obj); err != nil {
			fmt.Printf("failed to execute template: %v\n", err)
		}
	})
	fmt.Println(binding)
	if err := http.ListenAndServe(binding, nil); err != nil {
		die("unable to bind", err)
	}
}
