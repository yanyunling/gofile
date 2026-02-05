package middleware

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache *cache.Cache

// 初始化go-cache
func InitCache() {
	Cache = cache.New(time.Hour, time.Minute)
}
