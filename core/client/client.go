package client

import (
	"github.com/rafaeljusto/druns/core/types"
)

type Client struct {
	Id       int
	Name     types.Name
	Birthday types.Date
	revision uint64
}
