package main

import (
	"log"
	"os"

	"github.com/svopper/kalles_weather_dashboard_v2/app/server"
)

func main() {
	router := server.InstantiateControllers()
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("PORT environment variable not set. Defaulting to 8080")
		port = "8080"
	}
	router.Run(":" + port)
}
