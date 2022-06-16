package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/svopper/kalles_weather_dashboard_v2/app/server/ocean"
	"github.com/svopper/kalles_weather_dashboard_v2/app/server/temperature"
)

func InstantiateControllers() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.HTML(http.StatusInternalServerError, "error.go.tmpl", gin.H{
				"error": err,
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.LoadHTMLGlob("app/server/templates/*")
	router.Static("/assets", "app/server/assets")
	router.Static("/.well-known/acme-challenge", "app/server/.well-known/acme-challenge")

	router.GET("/", temperature.GetIndex)
	router.GET("/ocean", ocean.GetOcean)

	return router
}
