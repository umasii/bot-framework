package Utilities

import "strings"

func IsRandomSize(size string) bool {
	lowercaseSize := strings.ToLower(size)
	return lowercaseSize == "random" || lowercaseSize == "rand" || lowercaseSize == "any"
}