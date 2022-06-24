package service

import (
	"github.com/gin-gonic/gin"
	"main/middleware"
	"main/module/user"
)

func userRouter(e *gin.Engine) {
	authMiddleware := middleware.AuthMiddleWare()
	e.Use(authMiddleware.MiddlewareFunc())
	{
		e.GET("/getUserByPage", user.GetUserByPageHandler)
	}
}
