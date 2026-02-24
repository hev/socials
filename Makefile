.PHONY: build install clean test run

BINARY := socials

build:
	go build -o $(BINARY) .

install:
	go install .

clean:
	rm -f $(BINARY)
	go clean

test:
	go test ./...

run: build
	./$(BINARY)
