# ==============================================================================
# Vars
# ==============================================================================
DEVVERSION:=$(shell git describe --tags)

BUILDCMD:=CGO_ENABLED=1 go build

# ==============================================================================
# Directorys
# ==============================================================================
BUILDDIR:= .

# ==============================================================================
# Recipes
# ==============================================================================

$(BUILDDIR)/tnt: main.go 
	$(BUILDCMD) -ldflags="-X main.Version=$(DEVVERSION)" -o $@

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy


.phony: clean remake dev all test tst


dev:
	find . -type f -name "*.go" | entr make remake

test:
	find . -type f -name "*.go" | entr make tst 

tst:
	clear
	go test

remake: 
	make clean
	make tnt
	
clean:
	rm -rf $(BUILDDIR)/tnt

