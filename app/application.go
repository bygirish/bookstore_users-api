package app

import (
	"github.com/bygirish/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	MapUrls()

	logger.Info("about to start the application")
	router.Run(":8080")

}
