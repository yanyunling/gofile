package common

type TokenResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	IsAdmin      bool   `json:"isAdmin"`
	Username     string `json:"username"`
	NickName     string `json:"nickName"`
	Icon         string `json:"icon"`
	Functions    string `json:"functions"`
}

type TokenCache struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserId       string `json:"userId"`
	IsAdmin      bool   `json:"isAdmin"`
}

type SignIn struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CaptchaId string `json:"captchaId"`
	CaptchaX  int    `json:"captchaX"`
	CaptchaY  int    `json:"captchaY"`
}
