package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/services"
	"time"
)

func Index(c *gin.Context) {
	ip := c.ClientIP()
	currentTimestamp := time.Now().Unix()
	access := services.GetAccessInstance()

	var result interface{}
	if result = access.Get(ip, currentTimestamp); result == 0 || result == 61 {
		result = "Error!"
	}

	c.JSON(http.StatusOK, result)
}
