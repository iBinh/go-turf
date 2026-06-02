.PHONY: test bench lint cover clean

test:
	go test ./... -v -count=1

bench:
	go test ./... -bench=. -benchmem

lint:
	golangci-lint run ./...

cover:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html
