default: build

build: fix
	go build -v .

test: fix
	go test -v

fix: *.go
	goimports -l -w .
	gofmt -l -w .
