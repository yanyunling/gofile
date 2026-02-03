package entity

import "gofile/model/common"

type User struct {
	Id            string `json:"id" db:"id"`
	Username      string `json:"username" db:"username"`
	Password      string `json:"password" db:"password"`
	Enable        int    `json:"enable" db:"enable"`
	PublicAuth    int    `json:"publicAuth" db:"public_auth"`
	ProtectedAuth int    `json:"protectedAuth" db:"protected_auth"`
	PrivateAuth   int    `json:"privateAuth" db:"private_auth"`
}

type UserCondition struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

// 获取超级管理员信息
func GetAdminUser() User {
	return User{
		Id:            "019b6e1b-2916-74f1-be0b-8673ebdf3862",
		Username:      "admin",
		Password:      common.AdminPassword,
		Enable:        1,
		PublicAuth:    1,
		ProtectedAuth: 1,
		PrivateAuth:   1,
	}
}
