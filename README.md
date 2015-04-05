# guarid-go

[![Build Status](https://drone.io/github.com/enmand/quarid-go/status.png)](https://drone.io/github.com/enmand/quarid-go/latest)
[![docs examples](https://sourcegraph.com/api/repos/github.com/enmand/quarid-go/.badges/docs-examples.svg)](https://sourcegraph.com/github.com/enmand/quarid-go)

## Quick Start  ##
### To install: ###

    go get bitbucket.org/enmand/quarid-go

### To run ###
If your `$PATH` includes `$GOPATH/bin`: run `quarid-go --config config.json`
Othwise, run: `$GOPATH/bin/quarid-go --config config.json`

## Getting Started
### Requirements

This bot has been tested with *Go 1.4.1*, although *Go 1.3* should work as well.

### Your $GOPATH

Make sure your [`$GOPATH`](https://golang.org/doc/code.html#GOPATH) is
configured for your shell. You should also update your `$PATH` to include
`$GOPATH/bin`, so you can run `quarid-go`.

### Installing Quarid locally

If you don't want to install `Quarid` in the above way for a binary `quarid-go`
binary in `$GOPATH/bin`, you can see the *Hacking* section to find out more.

## Hacking

Quarid is a flexible IRC bot, built using a Go IRC framework called 
[go-ircevent](https://github.com/thoj/go-ircevent).

### Check out repository

	mkdir -p $GOPATH/src/bitbucket.org/enmand/
	cd $GOPATH/src/bitbucket.org/enmand
	git clone git@bitbucket.org:enmand/quarid.go

### Plugins

Plugins should have a file named *`<PluginName>_plugin.go`. The plugin should
call the [`RegisterPlugin`](https://gowalker.org/bitbucket.org/enmand/quarid-go#RegisterPlugin)
method.

See `random_plugin.go` for an example.

### Local build

To build a local build of `quarid-go` run `make bin/quarid-go`
