package internal_test

import (
	"testing"

	"voidedtech.com/gxs/internal"
)

func TestNoBlocks(t *testing.T) {
	_, err := internal.Parse([]byte(""))
	if err == nil || err.Error.Error() != "parsing: no blocks found" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`



`))
	if err == nil || err.Error.Error() != "parsing: no blocks found" {
		t.Error("wrong error")
	}
}

func TestUnclosedBlock(t *testing.T) {
	_, err := internal.Parse([]byte(`
myblock => {`))
	if err == nil || err.Error.Error() != "parsing: unclosed block at block: 1" {
		t.Error("wrong error")
	}
}

func TestExpectStart(t *testing.T) {
	_, err := internal.Parse([]byte(`
myblock`))
	if err == nil || err.Error.Error() != "parsing: expected start of block" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`
myblock => { => {`))
	if err == nil || err.Error.Error() != "parsing: invalid start block" {
		t.Error("wrong error")
	}
}

func TestEmptyBlock(t *testing.T) {
	_, err := internal.Parse([]byte(`
myblock => {
}
`))
	if err == nil || err.Error.Error() != "parsing: empty block found" {
		t.Error("wrong error")
	}
}

func TestUnknownAction(t *testing.T) {
	_, err := internal.Parse([]byte(`action => {
		save
}`))
	if err == nil || err.Error.Error() != "parsing: unknown action" {
		t.Error("wrong error")
	}
}

func TestUnknownMode(t *testing.T) {
	_, err := internal.Parse([]byte(`actions => {
		commit
}`))
	if err == nil || err.Error.Error() != "parsing: unknown mode in block" {
		t.Error("wrong error")
	}
}

func TestNoPattern(t *testing.T) {
	_, err := internal.Parse([]byte(`action => {
		commit
}`))
	if err == nil || err.Error.Error() != "parsing: no pattern" {
		t.Error("wrong error")
	}
}

func TestOverwritePattern(t *testing.T) {
	_, err := internal.Parse([]byte(`
pattern => {
	abc
}
pattern => {
	abc
}
action => {
		commit
}`))
	if err == nil || err.Error.Error() != "parsing: pattern not committed" {
		t.Error("wrong error")
	}
}

func TestBadStitchMode(t *testing.T) {
	_, err := internal.Parse([]byte(`
pattern => {
	abc
}
action => {
		commit
}`))
	if err == nil || err.Error.Error() != "parsing: invalid stitch mode" {
		t.Error("wrong error")
	}
}

func TestUncommit(t *testing.T) {
	_, err := internal.Parse([]byte(`
pattern => {
	abc
}`))
	if err == nil || err.Error.Error() != "parsing: uncommitted pattern" {
		t.Error("wrong error")
	}
}

func TestStitchModeSetting(t *testing.T) {
	_, err := internal.Parse([]byte(`
mode => {
	xstitch
	yy
}`))
	if err == nil || err.Error.Error() != "parsing: incorrect stitch mode setting" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`
mode => {
	xstitch
}
mode => {
	topedge
}`))
	if err == nil || err.Error.Error() != "parsing: stitching not committed" {
		t.Error("wrong error")
	}
}

func TestNoActions(t *testing.T) {
	_, err := internal.Parse([]byte(`
mode => {
	xstitch
}
mode => {
	xstitch
}`))
	if err == nil || err.Error.Error() != "parsing: no actions, nothing committed?" {
		t.Error("wrong error")
	}
}

func TestBadPalette(t *testing.T) {
	_, err := internal.Parse([]byte(`palette => {
	x => y => z
}`))
	if err == nil || err.Error.Error() != "parsing: invalid palette assignment" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`palette => {
	xr => y
}`))
	if err == nil || err.Error.Error() != "parsing: only single characters allowed" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`palette => {
	x => y
	x => z
}`))
	if err == nil || err.Error.Error() != "parsing: character re-used within palette" {
		t.Error("wrong error")
	}
}

func TestUnknownSymbol(t *testing.T) {
	_, err := internal.Parse([]byte(`palette => {
	x => y
}
mode => {
	xstitch
}
pattern => {
	z
}
action => {
	commit
}`))
	if err == nil || err.Error.Error() != "parsing: symbol unknown" {
		t.Error("wrong error")
	}
}

func TestParser(t *testing.T) {
	_, err := internal.Parse([]byte(`
# allow comments
palette => {
	x => NONE
	y => red
	z => #231234
}
mode => {
	xstitch
}
pattern => {
	zxxxy
	xxzyy
	xxyyy

	xxxxx
}
action => {
	commit
}
palette => {
	x => NONE
	y => red
	r => #231234
}
mode => {
	bottomedge
}
pattern => {
	xxxxx
	xrrrr
	xxxxx
	r
}
offset => {
	1x2
}
action => {
	commit
}
`))
	if err != nil {
		t.Error("is valid")
	}
}

func TestBadOffset(t *testing.T) {
	_, err := internal.Parse([]byte(`
offset => {
	1
	2
}`))
	if err == nil || err.Error.Error() != "parsing: invalid offset" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`
offset => {
	12
}`))
	if err == nil || err.Error.Error() != "parsing: offset should be Width[x]Height" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`
offset => {
	BADx2
}`))
	if err == nil || err.Error.Error() != "strconv.Atoi: parsing \"BAD\": invalid syntax" {
		t.Error("wrong error")
		t.Error(err.Error.Error())
	}
}

func TestBadInclude(t *testing.T) {
	_, err := internal.Parse([]byte(`
include => {
	1
	2
}`))
	if err == nil || err.Error.Error() != "open 1: no such file or directory" {
		t.Error("wrong error")
	}
}

func TestSingleLineParserError(t *testing.T) {
	_, err := internal.Parse([]byte(`
action => {
	commit}
`))
	if err == nil || err.Error.Error() != "parsing: unclosed block at block: 1" {
		t.Error("wrong error")
	}
	_, err = internal.Parse([]byte(`
action => {action => {commit}
`))
	if err == nil || err.Error.Error() != "parsing: single-line start of block invalid" {
		t.Error("wrong error")
		t.Error(err.Error.Error())
	}
}

func TestSingleLineParser(t *testing.T) {
	_, err := internal.Parse([]byte(`
# allow comments
palette => {
	x => NONE
	y => red
	z => #231234
}
mode => {
	xstitch
}
pattern => {
	zxxxy
	xxzyy
	xxyyy

	xxxxx
}
action => {commit}
palette => {
	x => NONE
	y => red
	r => #231234
}
mode => {bottomedge}
pattern => {
	xxxxx
	xrrrr
	xxxxx
	r
}
offset => {1x2}
action => {
	commit
}
`))
	if err != nil {
		t.Error("is valid")
	}
}
