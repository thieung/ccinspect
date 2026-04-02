.PHONY: build install uninstall test clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/ccinspect ./cmd/ccinspect

install: build
	@echo "Installing ccinspect to /usr/local/bin..."
	@sudo cp bin/ccinspect /usr/local/bin/ccinspect
	@echo "✓ ccinspect installed. Run 'ccinspect --help' to get started."

uninstall:
	@sudo rm -f /usr/local/bin/ccinspect
	@echo "✓ ccinspect uninstalled."

test:
	go test ./...

clean:
	go clean
	$(RM) -r bin/
