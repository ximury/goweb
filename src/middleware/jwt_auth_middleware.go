package middleware

import (
	"crypto/md5"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"logger"
	"main/module/user"
	"time"
)

type JwtUser struct {
	UserName string
}

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func AuthMiddleWare() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		// 中间件名称
		Realm: "gin-jwt",
		Key:   []byte("secret key"),
		// token 过期时间
		Timeout: 24 * time.Hour,
		// token 刷新最大时间
		MaxRefresh: 24 * time.Hour,
		// 身份验证的 key 值
		IdentityKey: identityKey,
		// 登录期间的回调的函数
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(JwtUser); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		// 解析并设置用户身份信息
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return JwtUser{
				UserName: claims[identityKey].(string),
			}
		},
		// 根据登录信息对用户进行身份验证的回调函数
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVars login
			if err := c.ShouldBind(&loginVars); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userName := loginVars.Username
			password := loginVars.Password
			res := user.SelectByUsername(userName)
			if res != nil && MD5(password) == res.Password {
				return JwtUser{
					UserName: userName,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		// 接收用户信息并编写授权规则
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(JwtUser); ok {
				return true
			}
			return false
		},
		// 自定义处理未进行授权的逻辑
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// token 检索模式，用于提取 token，默认值为 header:Authorization
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		logger.Debugf("JWT err: %v" + err.Error())
	}
	// https://jwt.io/ 解析
	return authMiddleware
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
