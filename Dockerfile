# syntax=docker/dockerfile:1
FROM golang:1.18-alpine as builder
ADD . /go/src/github.com/emarcey/data-vault
WORKDIR /go/src/github.com/emarcey/data-vault

# RUN go get ./data-vault
RUN go install -buildvcs=false

EXPOSE 6666

ENTRYPOINT ["/go/bin/data-vault"]
