BIN := bin/
GXS := $(BIN)gxs

all: $(GXS)

$(GXS): $(shell find cmd/ -type f)
	go build -o $(GXS) cmd/main.go

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
