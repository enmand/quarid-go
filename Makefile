all: bin/quarid-go

bin/quarid-go: godep
	mkdir -p bin
	go build -o bin/quarid-go .

$GOPATH/bin/quarid-go: godep
	go install .

.PHONY: godep
godep:
	go get github.com/tools/godep
	godep restore
