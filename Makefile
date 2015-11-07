DEV_PID=/tmp/quarid-dev.pid

all: bin/quarid-go bin/quarid-irc

bin/quarid-go: | deps.lock
	mkdir -p bin
	go build -o bin/quarid-go ./cmd/quaridd


bin/quarid-irc: | deps.lock
	mkdir -p bin
	go build -o bin/quarid-irc ./cmd/quaridirc

$GOPATH/bin/quarid-go: | deps.lock
	go install ./cmd/quaridd


$GOPATH/bin/quarid-irc: | deps.lock
	go install ./cmd/quaridirc

deps.lock:
	go get github.com/tools/godep
	godep restore
	touch deps.lock

.PHONY: dev dev_restart dev_kill clean

clean:
	rm bin/quarid-go

dev:
	@make dev_restart
	@fswatch -o . -e Godeps -e bin | xargs -n1 -I{}  make dev_restart || make dev_kill

dev_kill:
	@kill `cat $(DEV_PID)` || true

dev_restart: bin/quarid-go bin/quarid-irc
	@make dev_kill
	./bin/quarid-go & echo $$! > $(DEV_PID)
