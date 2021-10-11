BIN      := bin/
CASES    := $(shell ls ../examples/*.gxs) $(shell ls inputs/*.gxs)
FORMATS  := html ascii
EXPECT   := $(shell find outputs -type f)

.PHONY: $(CASES) $(EXPECT)

check: clean example options expect

example: $(CASES)
  
options:
	cat inputs/readme.gxs | $(GXS) -option ascii-no-delimiter=true > $(BIN)nodelimiter.ascii

expect:  $(EXPECT)

$(EXPECT):
	diff -u $@ $(BIN)$(shell basename $@)

$(CASES):
	for f in $(FORMATS); do $(GXS) -format $$f -input $@ > $(BIN)$(shell basename $@).$$f; done

clean:
	rm -rf $(BIN)
	mkdir -p $(BIN)
