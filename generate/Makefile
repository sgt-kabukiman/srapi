default: build

build: fix
	go build -v .

fix: *.go
	goimports -l -w .
	gofmt -l -w .
