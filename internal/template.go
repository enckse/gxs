package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

const (
	gridLocation         = "x"
	isBottomEdge         = "bottomedge"
	isTopEdge            = "topedge"
	isRightEdge          = "rightedge"
	isLeftEdge           = "leftedge"
	isXStitch            = "xstitch"
	isTopLeftBottomRight = "tlbrline"
	isTopRightBottomLeft = "trblline"
	isHorizontalLine     = "hline"
	isVerticalLine       = "vline"
	vLineSymbol          = "|"
	hLineSymbol          = "---"
	hLinePartial         = "-"
	fontSize             = "font-size: 6pt"
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
		Legend  []string
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
		legend  []string
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

func (o HTMLPattern) initCells(j Pattern) ([]Cell, error) {
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
			style, hvLine, err := j.layout(y, x)
			if err != nil {
				return results, err
			}
			cell.Style = template.CSS(style)
			if hvLine != "" {
				cell.Value = template.HTML(hvLine)
			}
			results = append(results, cell)
			y += 1
		}
		x += 1
	}
	return results, nil
}

func (o HTMLPattern) newID(x, y int) string {
	left := o.pad(x)
	right := o.pad(y)
	return fmt.Sprintf("%s%s%s", left, gridLocation, right)
}

func (p Pattern) layout(x, y int) (string, string, error) {
	var style []string
	vLineColor := ""
	hLineColor := ""
	tlbrColor := ""
	trblColor := ""
	for _, e := range p.entries {
		for _, c := range e.cells {
			if c.x == x && c.y == y {
				s := ""
				switch e.mode {
				case isVerticalLine:
					vLineColor = e.color
					style = append(style, fontSize)
				case isHorizontalLine:
					hLineColor = e.color
					style = append(style, fontSize)
				case isTopLeftBottomRight:
					tlbrColor = e.color
					style = append(style, fontSize)
				case isTopRightBottomLeft:
					trblColor = e.color
					style = append(style, fontSize)
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
	hadHVLine := vLineColor != "" || hLineColor != ""
	if hadHVLine {
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
	}
	hadXLine := tlbrColor != "" || trblColor != ""
	if hadXLine {
		if hadHVLine {
			return "", "", fmt.Errorf("unable to do horizontal/vertical line with cross line")
		}
		if tlbrColor != "" && trblColor != "" {
			return "", "", fmt.Errorf("perform a cross stitch, not 2 different lines")
		}
		if tlbrColor != "" {
			sub = colorLine("\\", tlbrColor)
		}
		if trblColor != "" {
			sub = colorLine("/", trblColor)
		}
	}
	return strings.Join(style, ";"), sub, nil
}

func colorLine(symbol, color string) string {
	return fmt.Sprintf("<div style=\"color: %s\">%s</div>", color, symbol)
}

func (p Pattern) ToHTMLPattern() (HTMLPattern, error) {
	padString := ""
	padding := p.pad
	for padding > 0 {
		padString = fmt.Sprintf("0%s", padString)
		padding = padding - 1
	}
	obj := HTMLPattern{Size: p.size + 1, padding: padString}
	cells, err := obj.initCells(p)
	if err != nil {
		return obj, err
	}
	obj.Cells = cells
	obj.Legend = p.legend
	return obj, nil
}

func Build(p Pattern) ([]byte, error) {
	obj, err := p.ToHTMLPattern()
	if err != nil {
		return nil, err
	}
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
