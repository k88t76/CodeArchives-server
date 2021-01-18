package main

import (
	"github.com/k88t76/code_archives/server/config"
	"github.com/k88t76/code_archives/server/controllers"
	"github.com/k88t76/code_archives/server/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	controllers.StartWebServer()
}
