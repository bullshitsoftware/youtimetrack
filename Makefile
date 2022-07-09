.PHONY: cover
cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test
test:
	go test -v -cover ./...
