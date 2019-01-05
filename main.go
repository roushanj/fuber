package main

import (
	C "./controller"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/cabrequest", C.RequestRide)
	r.GET("/confirmride", C.ConfirmRequest)
	r.POST("/endride", C.EndTrip)

	r.Run(":8080")

}
