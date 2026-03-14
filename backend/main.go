package main

import (
	"log"
	"os"

	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/routes"
)

func main() {
	// Init DB
	database.Init()

	// Setup Router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🌿 Aloé Raw Go server running on :%s", port)
	r.Run(":" + port)
}
