package helper

import "strings"

func BuildSearchValue(value string) string {
	return "%" + strings.Replace(value, "%", "\\%", -1) + "%"
}
