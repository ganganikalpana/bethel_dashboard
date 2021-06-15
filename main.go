package main

import (
	"log"
	"os"

	"github.com/niluwats/bethel_dashboard/app"
)

func main() {
	os.Setenv("AZURE_AUTH_LOCATION", "github.com/niluwats/bethel_dashboard/azure_gosdk/quickstart.auth")
	log.Println("App started")
	app.Start()
}
