package internal

import (
	"fmt"
	"html/template"
	"testing"
)

func TestNewPattern(t *testing.T) {
	_, err := NewPattern(0)
	if err == nil {
		t.Error("invalid request, size is invalid")
	}
	_, err = NewPattern(1)
	if err != nil {
		t.Error("valid JSON result")
	}
}

func testCell(cell Cell, expectedID, expectedValue string, t *testing.T) {
	if cell.ID != expectedID {
		t.Error(fmt.Sprintf("%s != %s", cell.ID, expectedID))
	}
	if cell.Value != template.HTML(expectedValue) {
		t.Error(fmt.Sprintf("%s != %s", cell.Value, expectedValue))
	}
}

func TestToHTMLPattern(t *testing.T) {
	j, err := NewPattern(2)
	if err != nil {
		t.Error("valid JSON result")
	}
	pattern, err := j.ToHTMLPattern()
	if err != nil || pattern.Size != 3 {
		t.Error("invalid conversion")
	}
	if len(pattern.Cells) != 9 {
		t.Error("invalid cells")
	}
	testCell(pattern.Cells[0], "000x000", "", t)
	testCell(pattern.Cells[1], "001x000", "1", t)
	testCell(pattern.Cells[2], "002x000", "2", t)
	testCell(pattern.Cells[3], "000x001", "1", t)
	testCell(pattern.Cells[4], "001x001", "", t)
	testCell(pattern.Cells[5], "002x001", "", t)
	testCell(pattern.Cells[6], "000x002", "2", t)
	testCell(pattern.Cells[7], "001x002", "", t)
	testCell(pattern.Cells[8], "002x002", "", t)
}

func TestHTMLBuild(t *testing.T) {
	j, err := NewPattern(1)
	if err != nil {
		t.Error("pattern is valid")
	}
	b, err := Build(j, "html")
	if err != nil || len(b) == 0 {
		t.Error("invalid building result")
	}
}

func TestASCIIBuild(t *testing.T) {
	j, err := NewPattern(1)
	if err != nil {
		t.Error("pattern is valid")
	}
	b, err := Build(j, "ascii")
	if err != nil || len(b) == 0 {
		t.Error("invalid building result")
	}
}
