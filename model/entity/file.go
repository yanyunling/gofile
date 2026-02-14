package entity

type FileInfo struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UpdateTime int64  `json:"updateTime"`
	IsDir      bool   `json:"isDir"`
}
