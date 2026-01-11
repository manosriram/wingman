build:
	go build -o wingman ./cmd/main.go

run:
	./wingman

test:
	go test ./... -v

.PHONY: all test
