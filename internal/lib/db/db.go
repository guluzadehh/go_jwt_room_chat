package db

import "strings"

func Placeholders(n int) string {
	if n == 0 {
		return ""
	}

	return strings.Join(strings.Split(strings.Repeat("?", n), ""), ", ")
}
