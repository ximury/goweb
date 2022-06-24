package service

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"logger"
	"main/config"
	"main/middleware"
	"main/module"
)

type Option func(*gin.Engine)

var options []Option

// Include 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init() *gin.Engine {
	r := gin.New()
	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logger.Fatalf("Gin set trusted proxies failed! err: #%v", err)
	}
	r.Use(middleware.GinWebLog())
	r.Use(gin.Recovery())
	swagHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger/*any", swagHandler)

	authMiddleware := middleware.AuthMiddleWare()

	r.POST("/login", module.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	Include(userRouter)

	for _, opt := range options {
		opt(r)
	}
	return r
}

func StartApi() {
	// 初始化路由
	r := Init()
	configBase, err := config.GetChannelConfig()
	if err != nil {
		logger.Fatalf("Get config failed! err: #%v", err)
	}
	if err := r.Run(configBase.Webapi.Uri); err != nil {
		logger.Fatalf("Run web server failed! err: #%v", err)
	}
}
