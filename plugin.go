package main

import (
	"fmt"
	"regexp"

	"github.com/thoj/go-ircevent"
)

var availablePlugins map[string]Plugin
var activePlugins map[string]Plugin

func init() {
	availablePlugins = map[string]Plugin{}
	activePlugins = map[string]Plugin{}
}

type TriggerCommands map[*regexp.Regexp]Command

type Command struct {
	Handler func(con *irc.Connection, config *Config, m chan *Message)
}

type Plugin struct {
	Commands map[*regexp.Regexp]Command
}

func RegisterPlugin(name string, plugin Plugin) {
	fmt.Printf("Loading plugin %s\n", name)

	availablePlugins[name] = plugin
}

func InitializePlugin(config *Config) {
	for _, name := range config.Plugins {
		if plugin, ok := availablePlugins[name]; ok {
			fmt.Printf("Activating plugin: %v\n", name)
			activePlugins[name] = plugin
		}
	}
}

func findCommand(msgCh chan *Message, cmdCh chan Command) {
	msg := <-msgCh

	for _, plugin := range activePlugins {
		go func(plugin Plugin) {
			for trigger, cmd := range plugin.Commands {
				if trigger.Match([]byte(msg.Message())) {
					cmdCh <- cmd
				}
			}
		}(plugin)
	}
}
