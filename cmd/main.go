package main

import (
	"fmt"
	"flag"
	"net/http"
	"html/template"
	_ "embed"
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
		ID string
		Value string
	}
	Pattern struct {
		Size int
		Pad string
		PadChar string
		Cells []Cell
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

func main() {
	size := flag.Int("size", 25, "pattern size")
	pad := flag.Int("pad", 4, "id padding (with 0s)")
	char := flag.String("padchar", "x", "padding character")
	bind := flag.String("bind", ":10987", "local binding to use for server")
	flag.Parse()
	template, err := template.New("t").Parse(templateHTML)
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
	if obj.Size <= 0 {
		die("invalid size", fmt.Errorf("< 0"))
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := template.Execute(w, obj)
		if err != nil {
			fmt.Printf("failed to execute template: %v\n", err)
		}
    })
	binding := *bind
	fmt.Println(binding)
	if err := http.ListenAndServe(binding, nil); err != nil {
		die("unable to bind", err)
	}
}
