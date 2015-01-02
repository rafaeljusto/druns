package model

import (
	"github.com/rafaeljusto/druns/core/protocol"
)

type Class struct {
	Group Group
}

func (c *Class) Apply(request *protocol.ClassRequest) protocol.Translator {
	return nil
}

func (c *Class) Protocol() *protocol.ClassResponse {
	return nil
}
