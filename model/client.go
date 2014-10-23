package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Client struct {
	Id      bson.ObjectId `bson:"_id"`
	Name    string
	Classes []Class
}

type Class struct {
	Weekday time.Weekday
	Time    time.Time
}
