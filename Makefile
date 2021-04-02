BIN      := bin/
GXS      := $(BIN)gxs
TESTS    := $(PWD)/internal/
EXAMPLES := $(shell ls examples/)

.PHONY: $(TESTS) $(EXAMPLES)

all: test example $(GXS)

$(GXS): $(shell find cmd/ -type f) $(shell find internal/ -type f)
	go build -o $(GXS) cmd/main.go

test: $(TESTS)

example: $(EXAMPLES)

$(EXAMPLES):
	$(GXS) -input examples/$@ > $(BIN)$@.html
	diff -u $(BIN)$@.html expected/$@.html

$(TESTS):
	go test -v $@

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
