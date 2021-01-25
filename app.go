package main

import (
	"app/controllers"
	"app/services"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
)

func main() {
	r := gin.Default()
	r.GET("/access", controllers.Index)

	go func() {
		access := services.GetAccessInstance()

		gocron.Every(30).Minutes().Do(access.Prune)
		<-gocron.Start()
	}()

	r.Run()
}
