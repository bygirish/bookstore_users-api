package app

import (
	"github.com/bygirish/bookstore_users-api/controllers/ping"
	"github.com/bygirish/bookstore_users-api/controllers/users"
)

func MapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("users/:userId", users.GetUser)
	router.POST("users", users.CreateUser)
	// router.GET("users/search", controllers.SearchUser)

}
