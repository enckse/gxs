package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

const (
	gridLocation     = "x"
	isBottomEdge     = "bottomedge"
	isTopEdge        = "topedge"
	isRightEdge      = "rightedge"
	isLeftEdge       = "leftedge"
	isXStitch        = "xstitch"
	isHorizontalLine = "hline"
	isVerticalLine   = "vline"
	vLineSymbol      = "|"
	hLineSymbol      = "-----"
	hLinePartial     = "--"
)

var (
	//go:embed "template.html"
	templateHTML string
)

type (
	Cell struct {
		ID    string
		Value template.HTML
		Style template.CSS
	}
	HTMLPattern struct {
		Size    int
		padding string
		Cells   []Cell
	}
	cell struct {
		x int
		y int
	}
	entry struct {
		cells []cell
		mode  string
		color string
	}
	Pattern struct {
		size    int
		pad     int
		entries []entry
	}
)

func NewPattern(size int) (Pattern, error) {
	if size < 1 {
		return Pattern{}, fmt.Errorf("invalid size <= 0")
	}
	padding := len(fmt.Sprintf("%d", size)) + 2
	return Pattern{pad: padding, size: size}, nil
}

func (o HTMLPattern) pad(val int) string {
	padded := fmt.Sprintf("%s%d", o.padding, val)
	for len(padded) > len(o.padding) {
		padded = padded[1:]
	}
	return padded
}

func (o HTMLPattern) initCells(j Pattern) []Cell {
	var results []Cell
	x := 0
	for x < o.Size {
		y := 0
		for y < o.Size {
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
			cell.ID = o.newID(y, x)
			cell.Value = template.HTML(val)
			style, hvLine := j.layout(y, x)
			cell.Style = template.CSS(style)
			if hvLine != "" {
				cell.Value = template.HTML(hvLine)
			}
			results = append(results, cell)
			y += 1
		}
		x += 1
	}
	return results
}

func (o HTMLPattern) newID(x, y int) string {
	left := o.pad(x)
	right := o.pad(y)
	return fmt.Sprintf("%s%s%s", left, gridLocation, right)
}

func (p Pattern) layout(x, y int) (string, string) {
	var style []string
	vLineColor := ""
	hLineColor := ""
	for _, e := range p.entries {
		for _, c := range e.cells {
			if c.x == x && c.y == y {
				s := ""
				switch e.mode {
				case isVerticalLine:
					vLineColor = e.color
				case isHorizontalLine:
					hLineColor = e.color
				case isBottomEdge:
					s = "border-bottom-style: solid; border-bottom-color: "
				case isTopEdge:
					s = "border-top-style: solid; border-top-color: "
				case isRightEdge:
					s = "border-right-style: solid; border-right-color: "
				case isLeftEdge:
					s = "border-left-style: solid; border-left-color: "
				case isXStitch:
					s = "background-color: "
				}
				if s != "" {
					style = append(style, fmt.Sprintf("%s %s", s, e.color))
				}
			}
		}
	}
	sub := ""
	if vLineColor != "" {
		sub = colorLine(vLineSymbol, vLineColor)
	}
	if hLineColor != "" {
		if sub == "" {
			sub = colorLine(hLineSymbol, hLineColor)
		} else {
			hSub := colorLine(hLinePartial, hLineColor)
			sub = fmt.Sprintf("%s%s%s", hSub, colorLine(vLineSymbol, vLineColor), hSub)
		}
	}
	return strings.Join(style, ";"), sub
}

func colorLine(symbol, color string) string {
	return fmt.Sprintf("<div style=\"color: %s\">%s</div>", color, symbol)
}

func (p Pattern) ToHTMLPattern() HTMLPattern {
	padString := ""
	padding := p.pad
	for padding > 0 {
		padString = fmt.Sprintf("0%s", padString)
		padding = padding - 1
	}
	obj := HTMLPattern{Size: p.size + 1, padding: padString}
	obj.Cells = obj.initCells(p)
	return obj
}

func Build(p Pattern) ([]byte, error) {
	obj := p.ToHTMLPattern()
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
