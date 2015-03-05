package db

import (
	"fmt"
	"hash/crc64"
	"io"
	"reflect"
)

func Revision(object interface{}) uint64 {
	value := reflect.ValueOf(object)
	valueType := value.Type()

	h := crc64.New(crc64.MakeTable(crc64.ECMA))

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		if field.Tag.Get("revision") == "-" {
			continue
		}

		if value.Field(i).CanInterface() {
			io.WriteString(h, fmt.Sprintf("{%#v}", value.Field(i).Interface()))
		}
	}

	return h.Sum64()
}
