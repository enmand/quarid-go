DEV_PID=/tmp/quarid-dev.pid

all: bin/quarid-irc

Gopkg.lock:
	dep ensure

bin/quarid-irc:  Gopkg.lock
	mkdir -p bin
	go build -o bin/quarid-irc ./cmd/quaridirc

$GOPATH/bin/quarid-go: | Gopkg.lock
	go install ./cmd/quaridd


$GOPATH/bin/quarid-irc: | Gopkg.lock
	go install ./cmd/quaridirc

.PHONY: dev dev_restart dev_kill clean

metalint:
	gometalinter \
    --concurrency=2 --deadline=1m --sort=path \
    --disable=dupl --disable=vetshadow --enable=misspell \
		--enable nakedret --vendor \
    ./...

clean:
	rm bin/quarid-go

dev:
	@make dev_restart
	@fswatch -o . -e vendor -e bin | xargs -n1 -I{}  make dev_restart || make dev_kill

dev_kill:
	@kill `cat $(DEV_PID)` || true

dev_restart: bin/quarid-go bin/quarid-irc
	@make dev_kill
	./bin/quarid-go & echo $$! > $(DEV_PID)
