# ==============================================================================
# Vars
# ==============================================================================
DEVVERSION:=$(shell git describe --tags)
VERSION:=$(shell git describe --tags --abbr=0)

BUILDCMD:=CGO_ENABLED=1 go build

# ==============================================================================
# Directorys
# ==============================================================================
BUILDDIR:= .

# ==============================================================================
# Recipes
# ==============================================================================
GOSOURCE:=$(shell find . -name "*.go")

$(BUILDDIR)/tnt: $(GOSOURCE)
	$(BUILDCMD) -ldflags="-X main.Version=$(DEVVERSION)" -o $@

tnt.win.x86_64.exe: main.go
	GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows" CXX="zig c++ -target x86_64-windows" $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

tnt.linux.x86_64: main.go
	GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" CXX="zig c++ -target x86_64-linux" $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy



.phony: clean remake dev all test tst all serve server css

all: tnt.win.x86_64.exe tnt.linux.x86_64
	echo "done"

dev:
	find . -type f -name "*.go" | entr make remake

css:
	make -C ./pkg/serve/views/

serve:
	find . -type f -iname "*.go" -o -iname "*.html" -o -iname "*.css" | entr make server

test:
	find . -type f -name "*.go" | entr make tst

tst:
	clear
	go test

server:
	$(shell killall -q tnt)
	clear
	rm -f ./tnt
	make css
	make tnt
	./tnt serve --db ./devdb.sqlite --debug &

remake:
	clear
	rm -f ./tnt
	make tnt

clean:
	rm -rf $(BUILDDIR)/tnt
	rm -rf $(BUILDDIR)/debug.run
	rm -rf $(BUILDDIR)/tnt.win.x86_64.exe
	rm -rf $(BUILDDIR)/tnt.linux.x86_64
