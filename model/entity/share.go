package entity

type Share struct {
	Id         string `json:"id" db:"id"`
	ParentDir  string `json:"parentDir" db:"parent_dir"`
	Path       string `json:"path" db:"path"`
	Name       string `json:"name" db:"name"`
	Username   string `json:"username" db:"username"`
	ShareHours int    `json:"shareHours" db:"share_hours"`
	StartTime  int64  `json:"startTime" db:"start_time"`
	EndTime    int64  `json:"endTime" db:"end_time"`
}
