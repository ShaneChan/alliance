package util

import (
	"math/rand"
	"strings"
)

func DealString(content string) []string {
	return strings.Split(content, " ")
}

// Max 求最大值
func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

// Min 求最小值
func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// RandomInt 随机一个整数
func RandomInt() int {
	return rand.Int()
}

// RandomIntBetween 随机一个M到N之间的整数
func RandomIntBetween() {
}
