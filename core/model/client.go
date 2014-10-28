package model

import (
	"github.com/rafaeljusto/druns/core/protocol"
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	Id      bson.ObjectId `bson:"_id"`
	Name    string
	Classes Classes
}

func (c *Client) Apply(request *protocol.ClientRequest) protocol.Translator {
	return nil
}

func (c *Client) Protocol() *protocol.ClientResponse {
	return &protocol.ClientResponse{
		Id:      c.Id.String(),
		Name:    c.Name,
		Classes: c.Classes.Protocol(),
	}
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
