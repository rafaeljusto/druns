package dao

import (
	"strings"
)

func placeholders(tableFields []string) string {
	placeholders := make([]string, len(tableFields))

	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = "?"
	}

	return strings.Join(placeholders, ", ")
}
