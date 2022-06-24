package module

import (
	"github.com/gin-gonic/gin"
	"main/middleware"
)

func LoginHandler(c *gin.Context) {
	middleware.AuthMiddleWare().LoginHandler(c)
}
