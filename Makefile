################################################################################
# Vars
################################################################################
VERSION:=$(shell date '+%Y%m%d-%H%M%S')
	
################################################################################
# Directorys
################################################################################
BUILDDIR:= .

################################################################################
# Recipes
################################################################################
$(BUILDDIR)/tnt: main.go 
	go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $@

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy

.phony: clean remake dev

dev:
	find . -type f -name "*.go" | entr make remake

remake: 
	make clean
	make tnt
	
clean:
	rm -rf $(BUILDDIR)/tnt
