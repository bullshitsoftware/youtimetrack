.PHONY: cover
cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: build
build:
	go build -o bin/yttadd  cmd/yttadd/main.go
	go build -o bin/yttconf cmd/yttconf/main.go
	go build -o bin/yttdel  cmd/yttdel/main.go
	go build -o bin/yttdet  cmd/yttdet/main.go
	go build -o bin/yttsum  cmd/yttsum/main.go
