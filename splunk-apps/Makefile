SOURCES := $(patsubst %/default/app.conf,%,$(wildcard */default/app.conf))
TARGETS := $(SOURCES:%=%.spl)

default: $(TARGETS)

$(TARGETS): $(SOURCES)
	@echo "Building $@ from $<"
	@tar -czf $@ -C $< .
