package main

import (
	"os"

	"github.com/niluwats/bethel_dashboard/app"
	"github.com/niluwats/bethel_dashboard/logger"
)

func main() {
	os.Setenv("AZURE_AUTH_LOCATION", "github.com/niluwats/bethel_dashboard/azure_gosdk/quickstart.auth")
	logger.Info("Starting application")
	app.Start()
}
