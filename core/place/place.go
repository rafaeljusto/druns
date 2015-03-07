package place

import (
	"github.com/rafaeljusto/druns/core/types"
)

type Place struct {
	Id       int
	Name     types.Name
	Address  types.Address
	revision uint64
}
