package main

import (
	"github.com/niluwats/bethel_dashboard/app"
	"github.com/niluwats/bethel_dashboard/logger"
)

func main() {
	logger.Info("Starting application")
	app.Start()
}
