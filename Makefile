BIN   := bin/
GXS   := $(BIN)gxs
TESTS := $(PWD)/internal/

.PHONY: $(TESTS)

all: test $(GXS)

$(GXS): $(shell find cmd/ -type f) $(shell find internal/ -type f)
	go build -o $(GXS) cmd/main.go

test: $(TESTS)

$(TESTS):
	go test -v $@

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
