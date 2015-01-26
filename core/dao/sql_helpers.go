package dao

import (
	"fmt"
	"strings"
)

func placeholders(tableFields []string) string {
	placeholders := make([]string, len(tableFields))

	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ", ")
}
