BIN      := bin/
GXS      := $(BIN)gxs
FLAGS    := -ldflags "-X main.version=$(shell git log --format=%h -n 1)" -trimpath -buildmode=pie -mod=readonly -modcacherw
CASES    := $(shell ls examples/*.gxs) $(shell ls tests/inputs/*.gxs)
FORMATS  := html ascii
EXPECT   := $(shell find tests/outputs -type f)

.PHONY: $(CASES) $(EXPECT)

all: $(GXS) check

check: test example options expect

$(GXS): $(shell find cmd/ -type f) $(shell find internal/ -type f)
	go build -o $(GXS) $(FLAGS) cmd/main.go

test:
	go test -v ./...

example: $(CASES)
  
options:
	cat tests/inputs/readme.gxs | $(GXS) -option ascii-no-delimiter=true > $(BIN)nodelimiter.ascii

expect:  $(EXPECT)

$(EXPECT):
	diff -u $@ $(BIN)$(shell basename $@)

$(CASES):
	for f in $(FORMATS); do $(GXS) -format $$f -input $@ > $(BIN)$(shell basename $@).$$f; done

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
