# ==============================================================================
# Vars
# ==============================================================================
VERSION:=$(shell git describe --tags --abbrev=0)
DEVVERSION:=$(shell git describe --tags)
	
# ==============================================================================
# Directorys
# ==============================================================================
BUILDDIR:= .

# ==============================================================================
# Recipes
# ==============================================================================


$(BUILDDIR)/tnt: main.go 
	go build -ldflags="-s -w -X main.Version=$(DEVVERSION)" -o $@

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy





.phony: clean remake dev all tnt.mac.x86_64 tnt.mac.arm64 tnt.linux.x86_64 tnt.linux.arm64 tnt.windows.x86_64.exe tnt.windows.arm64.exe


all: $(BUILDDIR)/tnt.mac.x86_64 $(BUILDDIR)/tnt.mac.arm64 $(BUILDDIR)/tnt.linux.x86_64 $(BUILDDIR)/tnt.linux.arm64 $(BUILDDIR)/tnt.windows.x86_64.exe $(BUILDDIR)/tnt.windows.arm64.exe
	echo "done"

dev:
	find . -type f -name "*.go" | entr make remake

remake: 
	make clean
	make tnt
	
clean:
	rm -rf $(BUILDDIR)/tnt
	rm -rf $(BUILDDIR)/tnt.mac.x86_64
	rm -rf $(BUILDDIR)/tnt.mac.arm64
	rm -rf $(BUILDDIR)/tnt.linux.x86_64
	rm -rf $(BUILDDIR)/tnt.linux.arm64
	rm -rf $(BUILDDIR)/tnt.windows.x86_64.exe
	rm -rf $(BUILDDIR)/tnt.windows.arm64.exe

# ==============================================================================
# crossbuild 
# ==============================================================================
$(BUILDDIR)/tnt.mac.x86_64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILDDIR)/$@

$(BUILDDIR)/tnt.mac.arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILDDIR)/$@

$(BUILDDIR)/tnt.linux.x86_64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILDDIR)/$@

$(BUILDDIR)/tnt.linux.arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILDDIR)/$@

$(BUILDDIR)/tnt.windows.x86_64.exe:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILDDIR)/$@

$(BUILDDIR)/tnt.windows.arm64.exe:
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILDDIR)/$@


