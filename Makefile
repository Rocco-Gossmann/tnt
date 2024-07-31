# ==============================================================================
# Vars
# ==============================================================================
DEVVERSION:=$(shell git describe --tags)
VERSION:=$(shell git describe --tags --abbr=0)

BUILDCMD:=go build
TARGET:=tnt

# ==============================================================================
# Directorys
# ==============================================================================
BUILDDIR:= .

# ==============================================================================
# Recipes
# ==============================================================================
GOSOURCE:=$(shell find . -name "*.go")

$(BUILDDIR)/$(TARGET): $(GOSOURCE)
	$(BUILDCMD) -ldflags="-X main.Version=$(DEVVERSION)" -o $@

$(TARGET).win.x86_64.exe: main.go
	GOOS=windows GOARCH=amd64  $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

$(TARGET).linux.x86_64: main.go
	GOOS=linux GOARCH=amd64 $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

$(TARGET).linux.arm: main.go
	GOOS=linux GOARCH=arm $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

$(TARGET).mac.arm: main.go
	GOOS=darwin GOARCH=arm64 $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

$(TARGET).mac.x86_64: main.go
	GOOS=darwin GOARCH=amd64 $(BUILDCMD) -ldflags="-w -X main.Version=$(VERSION)" -o $@

setup: go.sum
	@echo "setup done"

go.sum: go.mod
	GOPRIVATE="github.com/rocco-gossmann" go mod tidy


.phony: clean remake dev all test tst all serve server css

all: $(TARGET).win.x86_64.exe $(TARGET).linux.x86_64 $(TARGET).linux.arm $(TARGET).mac.arm $(TARGET).mac.x86_64
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
	$(shell killall -q $(TARGET))
	clear
	rm -f ./$(TARGET)
	make css
	make $(TARGET) 
	./$(TARGET) serve --db ./devdb.sqlite --debug &

remake:
	clear
	rm -f ./$(TARGET)
	make $(TARGET)

clean:
	rm -rf $(BUILDDIR)/$(TARGET)
	rm -rf $(BUILDDIR)/debug.run
	rm -rf $(BUILDDIR)/$(TARGET).win.x86_64.exe
	rm -rf $(BUILDDIR)/$(TARGET).linux.x86_64
	rm -rf $(BUILDDIR)/$(TARGET).linux.arm
	rm -rf $(BUILDDIR)/$(TARGET).mac.x86_64
	rm -rf $(BUILDDIR)/$(TARGET).mac.arm
