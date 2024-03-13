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

tnt.mac.arm64: main.go
	GOOS=darwin GOARCH=arm64 $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

tnt.win.x86_64.exe: main.go
	GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows" CXX="zig c++ -target x86_64-windows" $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

tnt.linux.x86_64: main.go
	GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" CXX="zig c++ -target x86_64-linux" $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@
	

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy


.phony: clean remake dev all test tst all

all: tnt.mac.arm64 tnt.win.x86_64.exe tnt.linux.x86_64
	echo "done"




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
	rm -rf $(BUILDDIR)/tnt.mac.arm64
	rm -rf $(BUILDDIR)/tnt.win.x86_64.exe
	rm -rf $(BUILDDIR)/tnt.linux.x86_64
