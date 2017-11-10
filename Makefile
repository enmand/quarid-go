DEV_PID=/tmp/quarid-dev.pid

all: bin/quarid-irc

deps.lock: glide.lock
	glide up
	touch $@

bin/quarid-irc:  deps.lock
	mkdir -p bin
	go build -o bin/quarid-irc ./cmd/quaridirc

$GOPATH/bin/quarid-go: | deps.lock
	go install ./cmd/quaridd


$GOPATH/bin/quarid-irc: | deps.lock
	go install ./cmd/quaridirc

.PHONY: dev dev_restart dev_kill clean

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
