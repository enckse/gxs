package internal_test

import (
	"html/template"
	"testing"

	"voidedtech.com/gxs/internal"
)

func TestNewPattern(t *testing.T) {
	_, err := internal.NewPattern(0)
	if err == nil {
		t.Error("invalid request, size is invalid")
	}
	_, err = internal.NewPattern(1)
	if err != nil {
		t.Error("valid JSON result")
	}
}

func checkCell(t *testing.T, cell internal.Cell, expectedID, expectedValue string) {
	if cell.ID != expectedID {
		t.Errorf("%s != %s", cell.ID, expectedID)
	}
	if cell.Value != template.HTML(expectedValue) {
		t.Errorf("%s != %s", cell.Value, expectedValue)
	}
}

func TestToHTMLPattern(t *testing.T) {
	j, err := internal.NewPattern(2)
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
	checkCell(t, pattern.Cells[0], "000x000", "")
	checkCell(t, pattern.Cells[1], "001x000", "1")
	checkCell(t, pattern.Cells[2], "002x000", "2")
	checkCell(t, pattern.Cells[3], "000x001", "1")
	checkCell(t, pattern.Cells[4], "001x001", "")
	checkCell(t, pattern.Cells[5], "002x001", "")
	checkCell(t, pattern.Cells[6], "000x002", "2")
	checkCell(t, pattern.Cells[7], "001x002", "")
	checkCell(t, pattern.Cells[8], "002x002", "")
}

func TestHTMLBuild(t *testing.T) {
	j, err := internal.NewPattern(1)
	if err != nil {
		t.Error("pattern is valid")
	}
	b, err := internal.Build(j, "html", &internal.Option{})
	if err != nil || len(b) == 0 {
		t.Error("invalid building result")
	}
}

func TestASCIIBuild(t *testing.T) {
	j, err := internal.NewPattern(1)
	if err != nil {
		t.Error("pattern is valid")
	}
	b, err := internal.Build(j, "ascii", &internal.Option{})
	if err != nil || len(b) == 0 {
		t.Error("invalid building result")
	}
}
