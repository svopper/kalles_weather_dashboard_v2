package main

import "github.com/svopper/kalles_weather_dashboard_v2/app/server"

func main() {
	router := server.InstantiateControllers()
	router.Run(":8080")
}
