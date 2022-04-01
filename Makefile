fmt:
	go fmt

build:
	go mod vendor
	go build

run:
	go run main.go

unit:
	go test -v ./...

unit-coverage:
	go get golang.org/x/tools/cmd/cover
	go test -cover ./...