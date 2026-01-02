build:
	go build -o wingman ./cmd/main.go

run:
	./wingman

test:
	go test ./test/... -v

.PHONY: all test
