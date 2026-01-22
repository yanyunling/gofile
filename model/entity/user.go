package entity

import "gofile/model/common"

type User struct {
	Id        string `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	NickName  string `json:"nickName" db:"nick_name"`
	Password  string `json:"password" db:"password"`
	Enable    int    `json:"enable" db:"enable"`
	HasUpdate int    `json:"hasUpdate" db:"has_update"`
}

type UserCondition struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// 获取超级管理员信息
func GetAdminUser() User {
	return User{
		Id:        "019b6e1b-2916-74f1-be0b-8673ebdf3862",
		Username:  "admin",
		NickName:  "超级管理员",
		Password:  common.AdminPassword,
		Enable:    1,
		HasUpdate: 1,
	}
}
