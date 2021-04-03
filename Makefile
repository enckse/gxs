BIN      := bin/
GXS      := $(BIN)gxs
TESTS    := $(PWD)/internal/
EXAMPLES := $(shell ls examples/ | grep "\.gxs$$")

.PHONY: $(TESTS) $(EXAMPLES)

all: $(GXS) check

check: test example 

$(GXS): $(shell find cmd/ -type f) $(shell find internal/ -type f)
	go build -o $(GXS) cmd/main.go

test: $(TESTS)

example: $(EXAMPLES)

$(EXAMPLES):
	$(GXS) -format html -input examples/$@ > $(BIN)$@.html
	diff -u $(BIN)$@.html expected/$@.html
	$(GXS) -format ascii -input examples/$@ > $(BIN)$@.ascii
	diff -u $(BIN)$@.ascii expected/$@.ascii

$(TESTS):
	go test -v $@

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
