DEV_PID=/tmp/quarid-dev.pid

all: bin/quarid-irc

CMDS=\
	github.com/UnnoTed/fileb0x

# Vendor command or code generation

define CMD_VENDOR_BIN
vendor/bin/$(notdir $(1)): vendor/$(1) | vendor
	go build -a -o $$@ ./vendor/$(1)
BINS+=vendor/bin/$(notdir $(1))
vendor/$(1): Gopkg.lock
	dep ensure -v --vendor-only
endef

$(foreach cmd,$(CMDS),$(eval $(call CMD_VENDOR_BIN,$(cmd))))

generated: vendor/bin/fileb0x
	./vendor/bin/fileb0x langsupport.yaml

Gopkg.lock:
	dep ensure

bin/%: | generated Gopkg.lock
	mkdir -p bin
	go build -o bin/$* ./cmd/$*

$GOPATH/bin/quarid-go: | Gopkg.lock
	go install ./cmd/quaridd


$GOPATH/bin/quarid-irc: engines-js | Gopkg.lock
	go install ./cmd/quaridirc

.PHONY: dev dev_restart dev_kill clean

metalint:
	gometalinter \
    --concurrency=2 --deadline=1m --sort=path \
    --disable=dupl --disable=vetshadow --enable=misspell \
		--enable nakedret --vendor \
    ./...

clean: 
	rm -fR bin/
	rm -fR generated
	
dev:
	@make dev_restart
	@fswatch -o . -e vendor -e bin | xargs -n1 -I{}  make dev_restart || make dev_kill

dev_kill:
	@kill `cat $(DEV_PID)` || true

dev_restart: bin/quarid-go bin/quarid-irc
	@make dev_kill
	./bin/quarid-go & echo $$! > $(DEV_PID)
