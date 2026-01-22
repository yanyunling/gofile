package common

var (
	Port             string // 端口
	LogPath          string // 日志目录
	DataPath         string // 数据目录
	BasicTokenKey    string // token相关接口认证key前缀
	PrivateDirName   string // 私有的下载目录，登录可见，在数据目录下
	ProtectedDirName string // 个人的下载目录，登录可见，在数据目录下
	PublicDirName    string // 公开的下载目录，无需登录，在数据目录下
	AdminPassword    string // 超级管理员密码
)
