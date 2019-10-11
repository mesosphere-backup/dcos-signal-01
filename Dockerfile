FROM golang:1.13.1

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.20.0
