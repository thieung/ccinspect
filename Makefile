.PHONY: build install test clean

build:
	go build -o bin/ccinspect ./cmd/cli

install:
	go install ./cmd/cli

test:
	go test ./...

clean:
	go clean
	$(RM) -r bin/
