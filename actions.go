package main

import (
	"log"
	"os"
)

type CLIActions struct {
	actions map[string]func()
}

func New() *CLIActions {
	var actions CLIActions
	actions.init()
	return &actions
}

func (a *CLIActions) init() {
	a.actions = make(map[string]func())
	a.actions["start"] = func() { return }
}

func (a *CLIActions) Add(s string, h func()) {
	a.actions[s] = h
	return
}
func (a *CLIActions) Run(s string) {
	a.actions[s]()
	return
}

func printHello() {
	log.Print("Yoba work")
}

func gtfo() {
	log.Print("MAMKU YEBAL")
	os.Exit(0)
}
func cliInit() *CLIActions {
	cActions := New()
	cActions.Add("hell", gtfo)
	return cActions
}
