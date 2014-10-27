package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	Id      bson.ObjectId `bson:"_id"`
	Name    string
	Classes []Class
}
