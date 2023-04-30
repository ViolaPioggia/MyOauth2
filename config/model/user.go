package model

type User struct {
	Username int64  `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	ID       int64  `form:"id" json:"id"`
}
