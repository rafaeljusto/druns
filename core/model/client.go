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

func (c *Client) Equal(other Client) bool {
	return *c == other
}

/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////

type Clients []Client

func (c Clients) Protocol() []protocol.ClientResponse {
	var response []protocol.ClientResponse
	for _, client := range c {
		response = append(response, *client.Protocol())
	}
	return response
}
