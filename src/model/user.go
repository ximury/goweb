package model

type User struct {
	UserId         int8   `gorm:"column:user_id;AUTO_INCREMENT;comment:用户ID" json:"user_id"`
	UserName       string `gorm:"column:username;comment:用户名" json:"username"`
	Password       string `gorm:"column:password;comment:密码" json:"password"`
	RoleId         int8   `gorm:"column:role_id;comment:角色ID" json:"role_id"`
	Status         string `gorm:"column:status;comment:用户是否禁用标志位，0为禁用，1为启用" json:"status"`
	Name           string `gorm:"column:name;comment:用户真实姓名" json:"name"`
	CreateByUserId int8   `gorm:"column:create_by_user_id;comment:创建者ID" json:"create_by_user_id"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
