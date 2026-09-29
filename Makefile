BINARY := simpleclient

.PHONY: build run clean

build:
	@test -f go.mod || go mod init simpleclient
	CGO_ENABLED=0 go build -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

clean:
	rm -f $(BINARY)

