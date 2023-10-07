package utils

import (
	"strings"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func RemoveWhitespaces(str string) string {
	return strings.ReplaceAll(str, " ", "")
}
