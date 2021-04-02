package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

const (
	padCharacter = "x"
)

var (
	//go:embed "template.html"
	templateHTML string
)

type (
	Cell struct {
		ID    string
		Value string
	}
	Pattern struct {
		Size    int
		Pad     string
		PadLength int
		PadChar string
		Cells   []Cell
		JSON    []Entry
	}
	Entry struct {
		Cells []string `json:"cells"`
		Mode  string   `json:"mode"`
		Color string   `json:"color"`
	}
	JSONPattern struct {
		size    int
		pad     int
		Entries []Entry
	}
)

func NewJSONPattern(size int) (JSONPattern, error) {
	if size < 1 {
		return JSONPattern{}, fmt.Errorf("invalid size <= 0")
	}
	padding := len(fmt.Sprintf("%d", size)) + 2
	return JSONPattern{pad: padding, size: size}, nil
}

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

func (p JSONPattern) ToPattern() Pattern {
	padString := ""
	padding := p.pad
	for padding > 0 {
		padString = fmt.Sprintf("0%s", padString)
		padding = padding - 1
	}
	obj := Pattern{Size: p.size + 1, PadLength: len(padString), Pad: padString, PadChar: padCharacter}
	obj.Cells = obj.initCells()
	obj.JSON = p.Entries
	return obj
}

func Build(p JSONPattern) ([]byte, error) {
	obj := p.ToPattern()
	tmpl, err := template.New("t").Parse(templateHTML)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, obj); err != nil {
		return nil, err
	}
	return b.Bytes(), err
}
