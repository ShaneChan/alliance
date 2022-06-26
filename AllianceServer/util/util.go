package util

import (
	"strings"
)

func DealString(content string) []string {
	return strings.Split(content, " ")
}

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
