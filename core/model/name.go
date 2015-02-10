package model

import "strings"

type Name struct {
	value string
}

func (n *Name) Set(value string) {
	n.value = strings.TrimSpace(value)
	n.value = strings.Title(n.value)
}

func (n Name) MarshalText() ([]byte, error) {
	return []byte(n.value), nil
}

func (n *Name) UnmarshalText(data []byte) (err error) {
	n.Set(string(data))
	return
}

func (n Name) String() string {
	return n.value
}
