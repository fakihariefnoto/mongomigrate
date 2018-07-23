package main

import (
	orc "github.com/fakihariefnoto/mongomigrate/controller/orchestrator"
	config "github.com/fakihariefnoto/mongomigrate/utils/config"
)

func main() {
	config.Init()
	orc.Init()
	orc.ShowApp()
}
