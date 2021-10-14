package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"voidedtech.com/stock"
)

type (
	flossColor struct {
		input    string
		resolved string
	}
	patternOffset struct {
		x int
		y int
	}
	patternBlock struct {
		lines []string
		mode  string
		err   error
	}
	// ParserError is an internal error associated with parsing patterns.
	ParserError struct {
		Error     error
		Backtrace []string
	}
	patternAction struct {
		palette    map[string]flossColor
		stitchMode string
		pattern    []string
		offset     patternOffset
	}
)

const (
	parserBlockStart = " => {"
	defaultBlock     = ""
	paletteAssign    = " => "
	noColor          = "NONE"
)

// NewParsingError returns a new gxs error for parsing.
func NewParsingError(message string) error {
	return stock.NewBasicCategoryError("parsing", message)
}

func next(stream []string) (patternBlock, int) {
	idx := 0
	blockCount := 0
	block := patternBlock{mode: defaultBlock}
	for idx < len(stream) {
		line := strings.TrimSpace(stream[idx])
		if strings.HasPrefix(line, "#") {
			line = ""
		}
		if len(line) > 0 {
			if blockCount > 0 {
				if line == "}" {
					if len(block.lines) == 0 {
						return patternBlock{err: NewParsingError("empty block found")}, 0
					}
					return block, idx + 1
				}
				block.lines = append(block.lines, line)
			} else {
				if strings.HasSuffix(line, parserBlockStart) {
					blockCount++
					mode, err := getBlockMode(line)
					if err != nil {
						return patternBlock{err: err}, 0
					}
					block.mode = mode
				} else {
					if strings.HasSuffix(line, "}") && strings.Contains(line, parserBlockStart) {
						sub := line[0 : len(line)-1]
						parts := strings.Split(sub, parserBlockStart)
						if len(parts) == 2 {
							mode, err := getBlockMode(sub)
							if err != nil {
								return patternBlock{err: NewParsingError("unable to read single line block")}, 0
							}
							block.mode = mode
							block.lines = []string{parts[1]}
							return block, 1
						}
						return patternBlock{err: NewParsingError("single-line start of block invalid")}, 0
					}
					return patternBlock{err: NewParsingError("expected start of block")}, 0
				}
			}
		}
		idx++
	}
	if blockCount > 0 {
		return patternBlock{err: NewParsingError(fmt.Sprintf("unclosed block at block: %d", blockCount))}, 0
	}
	return block, idx
}

func getBlockMode(line string) (string, error) {
	modeSection := strings.Split(line, parserBlockStart)
	if len(modeSection) != 2 {
		return "", NewParsingError("invalid start block")
	}
	return modeSection[0], nil
}

func (b patternBlock) isMatch(is string) bool {
	return len(b.lines) == 1 && b.lines[0] == is
}

func (b patternBlock) toError(message string) *ParserError {
	return &ParserError{Error: NewParsingError(message), Backtrace: b.lines}
}

func parseBlocks(blocks []patternBlock) ([]patternAction, *ParserError) {
	var actions []patternAction
	var action patternAction
	colorLookup := colors()
	for _, block := range blocks {
		switch block.mode {
		case "palette":
			action.palette = make(map[string]flossColor)
			for _, line := range block.lines {
				parts := strings.Split(line, paletteAssign)
				if len(parts) != 2 {
					return nil, block.toError("invalid palette assignment")
				}
				char := parts[0]
				color := parts[1]
				if len(char) != 1 {
					return nil, block.toError("only single characters allowed")
				}
				rawColor := color
				if val, ok := colorLookup[color]; ok {
					color = val
				}
				if _, ok := action.palette[char]; ok {
					return nil, block.toError("character re-used within palette")
				}
				action.palette[char] = flossColor{input: rawColor, resolved: color}
			}
		case "pattern":
			if len(action.pattern) > 0 {
				return nil, block.toError("pattern not committed")
			}
			action.pattern = block.lines
		case "action":
			if !block.isMatch("commit") {
				return nil, block.toError("unknown action")
			}
			if len(action.pattern) == 0 {
				return nil, block.toError("no pattern")
			}
			switch action.stitchMode {
			case isLeftEdge, isRightEdge, isTopEdge, isBottomEdge, isXStitch, isHorizontalLine, isVerticalLine, isTopLeftBottomRight, isTopRightBottomLeft:
				break
			default:
				return nil, block.toError("invalid stitch mode")
			}
			actions = append(actions, action)
			action.pattern = []string{}
			action.stitchMode = ""
		case "offset":
			if len(block.lines) != 1 {
				return nil, block.toError("invalid offset")
			}
			parts := strings.Split(block.lines[0], "x")
			if len(parts) != 2 {
				return nil, block.toError("offset should be Width[x]Height")
			}
			x, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, &ParserError{Error: err, Backtrace: block.lines}
			}
			y, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, &ParserError{Error: err, Backtrace: block.lines}
			}
			action.offset = patternOffset{x: x, y: y}
		case "mode":
			if len(block.lines) != 1 {
				return nil, block.toError("incorrect stitch mode setting")
			}
			line := block.lines[0]
			if action.stitchMode != "" {
				if action.stitchMode != line {
					return nil, block.toError("stitching not committed")
				}
			}
			action.stitchMode = line
		default:
			return nil, block.toError("unknown mode in block")
		}
	}
	if len(action.pattern) != 0 {
		return nil, &ParserError{Error: NewParsingError("uncommitted pattern")}
	}
	return actions, nil
}

func (a patternAction) toPatternError(message string) *ParserError {
	return &ParserError{Error: NewParsingError(message), Backtrace: a.pattern}
}

func buildPattern(actions []patternAction) (Pattern, *ParserError) {
	var entries []entry
	var maxSize = -1
	colorLegend := make(map[string]int)
	reverseColors := make(map[string]string)
	for _, action := range actions {
		tracking := make(map[string]map[string][]cell)
		for rawHeight, line := range action.pattern {
			height := rawHeight + action.offset.y
			if height > maxSize {
				maxSize = height
			}
			for rawWidth, chr := range line {
				width := rawWidth + action.offset.x
				if width > maxSize {
					maxSize = width
				}
				symbol := fmt.Sprintf("%c", chr)
				if color, ok := action.palette[symbol]; ok {
					if _, hasColor := tracking[color.resolved]; !hasColor {
						tracking[color.resolved] = make(map[string][]cell)
					}
					curColor := tracking[color.resolved]
					if _, hasMode := curColor[action.stitchMode]; !hasMode {
						curColor[action.stitchMode] = []cell{}
					}
					modeSet := curColor[action.stitchMode]
					modeSet = append(modeSet, cell{x: width + 1, y: height + 1})
					curColor[action.stitchMode] = modeSet
					tracking[color.resolved] = curColor
					reverseColors[color.resolved] = color.input
				} else {
					return Pattern{}, action.toPatternError("symbol unknown")
				}
			}
		}
		for color, modes := range tracking {
			if color == noColor {
				continue
			}
			count := 0
			for mode, cells := range modes {
				entry := entry{cells: cells, mode: mode, color: color}
				count += len(cells)
				entries = append(entries, entry)
			}
			if _, ok := colorLegend[color]; !ok {
				colorLegend[color] = 0
			}
			colorLegend[color] += count
		}
	}
	pattern, err := NewPattern(maxSize + 1)
	if err != nil {
		return pattern, &ParserError{Error: err}
	}
	var colorMapping []colorMap
	for k, v := range colorLegend {
		if lookup, ok := reverseColors[k]; ok {
			mapped := colorMap{input: lookup, output: k, count: v}
			colorMapping = append(colorMapping, mapped)
			continue
		}
		return pattern, &ParserError{Error: NewParsingError("unable to reverse map color")}
	}
	pattern.colors = colorMapping
	pattern.entries = entries
	return pattern, nil
}

func parseActions(b []byte) ([]patternAction, *ParserError) {
	lines := strings.Split(string(b), "\n")
	var blocks []patternBlock
	for {
		block, read := next(lines)
		if block.err != nil {
			return nil, &ParserError{Error: block.err, Backtrace: lines}
		}
		if read == 0 {
			break
		}
		var inserts []string
		if block.mode != defaultBlock {
			if block.mode == "include" {
				for _, line := range block.lines {
					data, err := os.ReadFile(line)
					if err != nil {
						return nil, &ParserError{Error: err, Backtrace: block.lines}
					}
					inserts = append(inserts, strings.Split(string(data), "\n")...)
				}
			} else {
				blocks = append(blocks, block)
			}
		}
		newLines := inserts
		newLines = append(newLines, lines[read:]...)
		lines = newLines
	}
	if len(blocks) == 0 {
		return nil, &ParserError{Error: NewParsingError("no blocks found")}
	}
	actions, pErr := parseBlocks(blocks)
	if pErr != nil {
		return nil, pErr
	}
	if len(actions) == 0 {
		return nil, &ParserError{Error: NewParsingError("no actions, nothing committed?")}
	}
	return actions, nil
}

// Parse handles parsing a pattern.
func Parse(b []byte) (Pattern, *ParserError) {
	actions, err := parseActions(b)
	if err != nil {
		return Pattern{}, err
	}
	return buildPattern(actions)
}
