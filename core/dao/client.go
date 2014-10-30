package dao

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	clientsCollection = "clients"
)

type Client struct {
	db *mgo.Database
}

func NewClient(db *mgo.Database) Client {
	return Client{
		db: db,
	}
}

func (dao *Client) Save(client *model.Client) error {
	if len(client.Id.Hex()) == 0 {
		client.Id = bson.NewObjectId()
	}

	_, err := dao.db.C(clientsCollection).Upsert(bson.M{
		"_id": client.Id,
	}, client)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *Client) Delete(client *model.Client) error {
	if err := dao.db.C(clientsCollection).RemoveId(client.Id); err != nil {
		return core.NewError(err)
	}
	return nil
}

func (dao *Client) FindById(id string) (model.Client, error) {
	query := dao.db.C(clientsCollection).FindId(id)

	var client model.Client
	if err := query.One(client); err != nil {
		return client, core.NewError(err)
	}
	return client, nil
}

func (dao *Client) FindAll() (model.Clients, error) {
	query := dao.db.C(clientsCollection).Find(bson.M{})

	var clients model.Clients
	if err := query.All(&clients); err != nil {
		return nil, core.NewError(err)
	}

	return clients, nil
}
