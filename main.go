package main

import (
	"github.com/k88t76/CodeArchives-server/config"
	"github.com/k88t76/CodeArchives-server/server"
	"github.com/k88t76/CodeArchives-server/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	server.StartWebServer()
}
