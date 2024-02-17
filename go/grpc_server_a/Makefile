.PHONY: \
	build \
	run \
	lint

build:
	go build -o dist/api cmd/api/*

run:
	go run cmd/api/*

lint:
	golangci-lint run ./...

test:
	go test -race ./...

testv:
	go test -race -v ./...

gen:
	go generate ./...
	$(call fmt)

define fmt
	go fmt ./...
	goimports -w cmd internal
endef
