package internal

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"sort"
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
	HTMLMode             = "html"
	ASCIIMode            = "ascii"
	asciiSymbols         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
	asciiSep             = "."
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
	colorMap struct {
		input  string
		output string
		count  int
	}
	Pattern struct {
		size    int
		pad     int
		entries []entry
		colors  []colorMap
	}
	asciiCell struct {
		top     bool
		bot     bool
		left    bool
		right   bool
		found   bool
		value   string
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
	var legend []string
	for _, mapped := range p.colors {
		legend = append(legend, fmt.Sprintf("color: %s (count %d)", mapped.input, mapped.count))
	}
	sort.Strings(legend)
	obj.Legend = legend
	return obj, nil
}

func (p Pattern) findASCIIEdges(y, x int) asciiCell {
	obj := asciiCell{}
	for _, entry := range p.entries {
		for _, cell := range entry.cells {
			if cell.x == x && cell.y == y {
				obj.found = true
				obj.entries = append(obj.entries, entry)
				switch entry.mode {
				case isBottomEdge:
					obj.bot = true
				case isTopEdge:
					obj.top = true
				case isLeftEdge:
					obj.left = true
				case isRightEdge:
					obj.right = true
				}
			}
		}
	}
	return obj
}

func ascii(p Pattern) ([]byte, error) {
	size := p.size + 2
	row := 0
	var array [][]asciiCell
	colorMap := make(map[string]string)
	colorPos := 0
	var warnings []string
	for row <= size {
		col := 0
		array = append(array, []asciiCell{})
		for col <= size {
			self := p.findASCIIEdges(row, col)
			above := p.findASCIIEdges(row-1, col)
			below := p.findASCIIEdges(row+1, col)
			left := p.findASCIIEdges(row, col+1)
			right := p.findASCIIEdges(row, col-1)
			self.top = self.top || above.bot
			self.bot = self.bot || below.top
			self.right = self.right || left.left
			self.left = self.left || right.right
			self.value = " "
			if self.found {
				hasHLine := false
				hasVLine := false
				hasTLBR := false
				hasTRBL := false
				color := ""
				isStitch := false
				for _, entry := range self.entries {
					switch entry.mode {
					case isTopLeftBottomRight:
						self.value = "\\"
						if hasHLine || hasVLine || hasTRBL {
							return nil, fmt.Errorf("unable to perform tlbr with other line")
						}
						hasTLBR = true
					case isTopRightBottomLeft:
						self.value = "/"
						if hasHLine || hasVLine || hasTLBR {
							return nil, fmt.Errorf("unable to perform trbl with other line")
						}
						hasTRBL = true
					case isHorizontalLine:
						self.value = "-"
						if hasTLBR || hasTRBL {
							return nil, fmt.Errorf("unable to make horizontal line with tlbr/trbl")
						}
						hasHLine = true
					case isVerticalLine:
						self.value = "|"
						if hasTLBR || hasTRBL {
							return nil, fmt.Errorf("unable to make vertical line with tlbr/trbl")
						}
						hasVLine = true
					case isXStitch:
						color = entry.color
						isStitch = true
					}
				}
				if hasVLine && hasHLine {
					self.value = "+"
				}
				if isStitch {
					if color == "" {
						return nil, fmt.Errorf("no color found")
					}
					if self.value != " " {
						warnings = append(warnings, "cannot have stitch+line in ASCII pattern")
					}
					symbol := ""
					if val, ok := colorMap[color]; ok {
						symbol = val
					} else {
						if colorPos < len(asciiSymbols) {
							symbol = fmt.Sprintf("%c", asciiSymbols[colorPos])
							colorPos += 1
						}
						colorMap[color] = symbol
					}
					self.value = symbol
				}
			}
			array[row] = append(array[row], self)
			col += 1
		}
		row += 1
	}

	var raw bytes.Buffer
	for _, row := range array {
		raw.WriteString("\n")
		for _, idx := range []int{0, 1} {
			for _, cell := range row {
				switch idx {
				case 0:
					raw.WriteString(asciiSep)
					if cell.top {
						raw.WriteString("-")
					} else {
						raw.WriteString(" ")
					}
				case 1:
					if cell.left {
						raw.WriteString("|")
					} else {
						raw.WriteString(" ")
					}
					raw.WriteString(cell.value)
				}
			}
			if idx == 0 {
				raw.WriteString("\n")
			}
		}
	}

	lines := reverse(strings.Split(raw.String(), "\n"))
	var final []string
	var prev string
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			prev = line
			continue
		}
		if strings.TrimSpace(strings.Replace(line, asciiSep, "", -1)) == "" {
			prev = line
			continue
		}
		final = append(final, prev)
		final = append(final, lines[idx:]...)
		break
	}

	var b bytes.Buffer
	for _, line := range reverse(final) {
		b.WriteString(fmt.Sprintf("%s\n", line))
	}
	b.WriteString("\n")
	b.WriteString("---\n")
	var legend []string
	for k, v := range colorMap {
		count := 0
		input := k
		for _, color := range p.colors {
			if color.output == k {
				count = color.count
				input = color.input
				break
			}
		}
		legend = append(legend, (fmt.Sprintf("color: %s => %s (count: %d)\n", v, input, count)))
	}
	sort.Strings(legend)
	for _, line := range legend {
		b.WriteString(line)
	}
	sort.Strings(warnings)
	tracked := make(map[string]int)
	for _, warning := range warnings {
		tracked[warning] += 1
	}
	for _, warning := range warnings {
		if count, ok := tracked[warning]; ok {
			b.WriteString(fmt.Sprintf("WARN: %s [%d]\n", warning, count))
			delete(tracked, warning)
		}
	}
	return b.Bytes(), nil
}

func reverse(array []string) []string {
	s := array
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Build(p Pattern, mode string) ([]byte, error) {
	switch mode {
	case HTMLMode:
		return html(p)
	case ASCIIMode:
		return ascii(p)
	}
	return nil, fmt.Errorf("unknown mode: %s", mode)
}

func html(p Pattern) ([]byte, error) {
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
