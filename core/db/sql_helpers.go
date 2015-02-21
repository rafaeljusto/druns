package db

import (
	"fmt"
	"strings"
)

func Placeholders(tableFields []string) string {
	placeholders := make([]string, len(tableFields))

	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ", ")
}

type Row interface {
	Scan(dest ...interface{}) error
}
