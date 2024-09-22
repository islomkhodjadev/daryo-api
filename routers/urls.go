package routers

import (
	"daryo-api/controllers"
	"daryo-api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes() *gin.Engine {
	r := gin.Default()

	// Client routes
	user := r.Group("api/v1/user")
	{
		user.POST("/register", controllers.RegisterHandler)
		user.POST("/login", controllers.LoginHandler)
		user.POST("/updatetoken", middlewares.TokenAuthMiddleware(), controllers.UpdateTokenHandler)
		user.POST("/delete", middlewares.TokenAuthMiddleware(), controllers.DeleteUser)

	}

	conversationRoutes := r.Group("api/v1/conversations")
	conversationRoutes.Use(middlewares.TokenAuthMiddleware())
	{

		conversationRoutes.POST("/ai", controllers.HandleClientMessage)

	}

	return r
}
