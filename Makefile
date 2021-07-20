BIN      := bin/
GXS      := $(BIN)gxs
TESTS    := $(PWD)/internal/
FLAGS    := -trimpath -buildmode=pie -mod=readonly -modcacherw
CASES    := $(shell ls examples/*.gxs) $(shell ls tests/inputs/*.gxs)
FORMATS  := html ascii
EXPECT   := $(shell find tests/outputs -type f)

.PHONY: $(TESTS) $(CASES) $(EXPECT)

all: $(GXS) check

check: test example options expect

$(GXS): $(shell find cmd/ -type f) $(shell find internal/ -type f)
	go build -o $(GXS) $(FLAGS) cmd/main.go

test: $(TESTS)

example: $(CASES)
  
options:
	cat tests/inputs/readme.gxs | $(GXS) -option ascii-no-delimiter=true > $(BIN)nodelimiter.ascii

expect:  $(EXPECT)

$(EXPECT):
	diff -u $@ $(BIN)$(shell basename $@)

$(CASES):
	for f in $(FORMATS); do $(GXS) -format $$f -input $@ > $(BIN)$(shell basename $@).$$f; done

$(TESTS):
	go test -v $@

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
