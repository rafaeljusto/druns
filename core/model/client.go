package model

import (
	"time"

	"github.com/rafaeljusto/druns/core/protocol"
)

type Client struct {
	Id       int
	Name     string
	Birthday time.Time
}

func (c *Client) Apply(request *protocol.ClientRequest) protocol.Translator {
	return nil
}

func (c *Client) Protocol() *protocol.ClientResponse {
	return nil
}
