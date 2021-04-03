BIN      := bin/
GXS      := $(BIN)gxs
TESTS    := $(PWD)/internal/
EXAMPLES := $(shell ls examples/ | grep "\.gxs$$")
FORMATS  := html ascii
EXPECT   := $(shell find expected -type f)

.PHONY: $(TESTS) $(EXAMPLES) $(EXPECT)

all: $(GXS) check

check: test example expect

$(GXS): $(shell find cmd/ -type f) $(shell find internal/ -type f)
	go build -o $(GXS) cmd/main.go

test: $(TESTS)

example: $(EXAMPLES)
  
expect:  $(EXPECT)

$(EXPECT):
	diff -u $@ $(BIN)$(shell basename $@)

$(EXAMPLES):
	for f in $(FORMATS); do $(GXS) -format $$f -input examples/$@ > $(BIN)$@.$$f; done

$(TESTS):
	go test -v $@

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
