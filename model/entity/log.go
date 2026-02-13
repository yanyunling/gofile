package entity

type LogLevel string

const (
	Debug LogLevel = "Debug"
	Info  LogLevel = "Info"
	Warn  LogLevel = "Warn"
	Error LogLevel = "Error"
)

type Log struct {
	Id         string   `json:"id" db:"id"`
	Title      string   `json:"title" db:"title"`
	Content    string   `json:"content" db:"content"`
	Level      LogLevel `json:"level" db:"level"`
	Username   string   `json:"username" db:"username"`
	CreateTime int64    `json:"createTime" db:"create_time"`
}

type LogCondition struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Level     LogLevel `json:"level"`
	Username  string   `json:"username"`
	StartTime int64    `json:"startTime"`
	EndTime   int64    `json:"endTime"`
}
