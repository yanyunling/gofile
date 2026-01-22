// 随机字符串工具类
package util

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

var (
	letters    = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	numbers    = []rune("0123456789")
	hexLetters = []rune("0123456789abcdef")
)

// 生成uuid-v7
func UUID() string {
	result, err := uuid.NewV7()
	if err != nil {
		return random(8, hexLetters) + "-" + random(4, hexLetters) + "-" + random(4, hexLetters) + "-" + random(4, hexLetters) + "-" + random(12, hexLetters)
	}
	return result.String()
}

// 生成没有连字符的uuid-v7
func UUIDNoHyphen() string {
	return strings.ReplaceAll(UUID(), "-", "")
}

// 随机生成字符串
func RandomString(length int) string {
	return random(length, letters)
}

// 随机生成数字字符串
func RandomNumber(length int) string {
	return random(length, numbers)
}

// 从对应字符集中随机生成指定长度字符串
func random(length int, arr []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = arr[rand.Intn(len(arr))]
	}
	return string(b)
}
