BINARY := simpleclient

.PHONY: build run clean

build:
	CGO_ENABLED=0 go build -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

clean:
	rm -f $(BINARY)

