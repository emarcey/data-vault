# only set if you need to use a local network
NETWORK=local-dev_default

VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+([-a-zA-Z0-9]*)?)' version.go)

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

docker-build:
	docker build -t emarcey/data-vault:${VERSION} .

docker-run:
	docker run \
		-p 6666:6666 \
		--net=${NETWORK} \
		-v ${CURDIR}:/go/src/github.com/emarcey/data-vault \
		emarcey/data-vault:${VERSION}