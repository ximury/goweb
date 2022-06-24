package user

import (
	"github.com/gin-gonic/gin"
)

func GetUserByPageHandler(c *gin.Context) {
	type Param struct {
		Page     int `form:"page" json:"page" binding:"required"`
		PageSize int `form:"pageSize" json:"pageSize" binding:"required"`
	}
	var param Param

	if err := c.ShouldBind(&param); err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "参数错误",
		})
		return
	}
	total, data := GetUserByPage(param.Page, param.PageSize)
	if data == nil {
		c.JSON(200, gin.H{
			"status_code": 200,
			"message":     "获取用户失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status_code": 200,
		"message":     "获取用户成功",
		"total":       total,
		"data":        data,
	})
}
