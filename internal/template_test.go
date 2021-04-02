package internal

import (
	"fmt"
	"testing"
)

func TestNewJSONPattern(t *testing.T) {
	_, err := NewJSONPattern(0)
	if err == nil {
		t.Error("invalid request, size is invalid")
	}
	_, err = NewJSONPattern(1)
	if err != nil {
		t.Error("valid JSON result")
	}
}

func testCell(cell Cell, expectedID, expectedValue string, t *testing.T) {
	if cell.ID != expectedID {
		t.Error(fmt.Sprintf("%s != %s", cell.ID, expectedID))
	}
	if cell.Value != expectedValue {
		t.Error(fmt.Sprintf("%s != %s", cell.Value, expectedValue))
	}
}

func TestToPattern(t *testing.T) {
	j, err := NewJSONPattern(2)
	if err != nil {
		t.Error("valid JSON result")
	}
	pattern := j.ToPattern()
	if pattern.Size != 3 || pattern.Pad != "000" || pattern.PadChar != "x" {
		t.Error("invalid conversion")
	}
	if len(pattern.Cells) != 9 {
		t.Error("invalid cells")
	}
	testCell(pattern.Cells[0], "000x000", "", t)
	testCell(pattern.Cells[1], "000x001", "1", t)
	testCell(pattern.Cells[2], "000x002", "2", t)
	testCell(pattern.Cells[3], "001x000", "1", t)
	testCell(pattern.Cells[4], "001x001", "", t)
	testCell(pattern.Cells[5], "001x002", "", t)
	testCell(pattern.Cells[6], "002x000", "2", t)
	testCell(pattern.Cells[7], "002x001", "", t)
	testCell(pattern.Cells[8], "002x002", "", t)
	if len(j.Entries) != len(pattern.JSON) {
		t.Error("pattern JSON not assigned")
	}
}

func TestBuild(t *testing.T) {
	j, err := NewJSONPattern(1)
	if err != nil {
		t.Error("pattern is valid")
	}
	b, err := Build(j)
	if err != nil || len(b) == 0 {
		t.Error("invalid building result")
	}
}
