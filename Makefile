.PHONY: build fmt
.DEFAULT_GOAL := build

build:
	CGO_ENABLED=0 GO111MODULES=on go build -ldflags="-s -w" -o aws-whois .

fmt:
	gofmt -w .
