package internal

import (
	"testing"
)

func TestNoBlocks(t *testing.T) {
	_, err := Parse([]byte(""))
	if err == nil || err.Error.Error() != "no blocks found" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`



`))
	if err == nil || err.Error.Error() != "no blocks found" {
		t.Error("wrong error")
	}
}

func TestUnclosedBlock(t *testing.T) {
	_, err := Parse([]byte(`
myblock => {`))
	if err == nil || err.Error.Error() != "unclosed block" {
		t.Error("wrong error")
	}
}

func TestExpectStart(t *testing.T) {
	_, err := Parse([]byte(`
myblock`))
	if err == nil || err.Error.Error() != "expected start of block" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`
myblock => { => {`))
	if err == nil || err.Error.Error() != "invalid start block" {
		t.Error("wrong error")
	}
}

func TestEmptyBlock(t *testing.T) {
	_, err := Parse([]byte(`
myblock => {
}
`))
	if err == nil || err.Error.Error() != "empty block found" {
		t.Error("wrong error")
	}
}

func TestUnknownAction(t *testing.T) {
	_, err := Parse([]byte(`action => {
		save
}`))
	if err == nil || err.Error.Error() != "unknown action" {
		t.Error("wrong error")
	}
}

func TestUnknownMode(t *testing.T) {
	_, err := Parse([]byte(`actions => {
		commit
}`))
	if err == nil || err.Error.Error() != "unknown mode in block" {
		t.Error("wrong error")
	}
}

func TestNoPattern(t *testing.T) {
	_, err := Parse([]byte(`action => {
		commit
}`))
	if err == nil || err.Error.Error() != "no pattern" {
		t.Error("wrong error")
	}
}

func TestOverwritePattern(t *testing.T) {
	_, err := Parse([]byte(`
pattern => {
	abc
}
pattern => {
	abc
}
action => {
		commit
}`))
	if err == nil || err.Error.Error() != "pattern not committed" {
		t.Error("wrong error")
	}
}

func TestBadStitchMode(t *testing.T) {
	_, err := Parse([]byte(`
pattern => {
	abc
}
action => {
		commit
}`))
	if err == nil || err.Error.Error() != "invalid stitch mode" {
		t.Error("wrong error")
	}
}

func TestUncommit(t *testing.T) {
	_, err := Parse([]byte(`
pattern => {
	abc
}`))
	if err == nil || err.Error.Error() != "uncommitted pattern" {
		t.Error("wrong error")
	}
}

func TestStitchModeSetting(t *testing.T) {
	_, err := Parse([]byte(`
mode => {
	xs
	yy
}`))
	if err == nil || err.Error.Error() != "incorrect stitch mode setting" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`
mode => {
	xs
}
mode => {
	te
}`))
	if err == nil || err.Error.Error() != "stitching not committed" {
		t.Error("wrong error")
	}
}

func TestNoActions(t *testing.T) {
	_, err := Parse([]byte(`
mode => {
	xs
}
mode => {
	xs
}`))
	if err == nil || err.Error.Error() != "no actions, nothing committed?" {
		t.Error("wrong error")
	}
}

func TestBadPalette(t *testing.T) {
	_, err := Parse([]byte(`palette => {
	x => y => z
}`))
	if err == nil || err.Error.Error() != "invalid palette assignment" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`palette => {
	xr => y
}`))
	if err == nil || err.Error.Error() != "only single characters allowed" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`palette => {
	x => y
	x => z
}`))
	if err == nil || err.Error.Error() != "character re-used within palette" {
		t.Error("wrong error")
	}
}

func TestUnknownSymbol(t *testing.T) {
	_, err := Parse([]byte(`palette => {
	x => y
}
mode => {
	xs
}
pattern => {
	z
}
action => {
	commit
}`))
	if err == nil || err.Error.Error() != "symbol unknown" {
		t.Error("wrong error")
	}
}

func TestParser(t *testing.T) {
	_, err := Parse([]byte(`
# allow comments
palette => {
	x => NONE
	y => red
	z => #231234
}
mode => {
	xs
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
	be
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
	_, err := Parse([]byte(`
offset => {
	1
	2
}`))
	if err == nil || err.Error.Error() != "invalid offset" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`
offset => {
	12
}`))
	if err == nil || err.Error.Error() != "offset should be Width[x]Height" {
		t.Error("wrong error")
	}
	_, err = Parse([]byte(`
offset => {
	BADx2
}`))
	if err == nil || err.Error.Error() != "strconv.Atoi: parsing \"BAD\": invalid syntax" {
		t.Error("wrong error")
		t.Error(err.Error.Error())
	}
}
