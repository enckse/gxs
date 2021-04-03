package internal

import (
	"testing"
)

func TestSetInvalid(t *testing.T) {
	o := &Option{}
	err := o.Set("a=x=y")
	if err == nil || err.Error() != "invalid key=value pair" {
		t.Error("is invalid")
	}
	err = o.Set("a=x")
	if err == nil || err.Error() != "unknown option" {
		t.Error("bad option")
	}
	err = o.Set("ascii-no-delimiter=abc")
	if err == nil || err.Error() != "invalid boolean value" {
		t.Error("bad boolean")
	}
}

func TestSetASCIIDelimiter(t *testing.T) {
	o := &Option{}
	err := o.Set("ascii-no-delimiter=true")
	if err != nil || !o.asciiNoDelimiter {
		t.Error("valid")
	}
	err = o.Set("ascii-no-delimiter=false")
	if err != nil || o.asciiNoDelimiter {
		t.Error("valid")
	}
}
