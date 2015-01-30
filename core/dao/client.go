package dao

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type Client struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewClient(sqler SQLer, ip net.IP, agent int) Client {
	return Client{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "client",
		tableFields: []string{
			"id",
			"name",
			"birthday",
		},
	}
}

func (dao *Client) Save(c *model.Client) error {
	if dao.Agent == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
	}

	var operation model.LogOperation

	if c.Id > 0 {
		if err := dao.insert(c); err != nil {
			return err
		}

		operation = model.LogOperationCreation

	} else {
		if err := dao.update(c); err != nil {
			return err
		}

		operation = model.LogOperationUpdate
	}

	clientLogDAO := NewClientLog(dao.SQLer, dao.IP, dao.Agent)
	return clientLogDAO.save(c, operation)
}

func (dao *Client) insert(c *model.Client) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields[1:]),
	)

	row := dao.SQLer.QueryRow(
		query,
		c.Name,
		c.Birthday,
	)

	if err := row.Scan(&c.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *Client) update(c *model.Client) error {
	if lastClient, err := dao.FindById(c.Id); err == nil && lastClient.Equal(*c) {
		// Nothing changed
		return nil

	} else if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, birthday = $2 WHERE id = $3",
		dao.tableName,
	)

	_, err := dao.SQLer.Exec(
		query,
		c.Name,
		c.Birthday,
		c.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *Client) FindById(id int) (model.Client, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	c, err := dao.load(row)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (dao *Client) FindAll() (model.Clients, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, core.NewError(err)
	}

	var clients model.Clients

	for rows.Next() {
		c, err := dao.load(rows)
		if err != nil {
			return nil, err
		}

		clients = append(clients, c)
	}

	return clients, nil
}

func (dao *Client) load(row row) (model.Client, error) {
	var c model.Client

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Birthday,
	)

	if err == sql.ErrNoRows {
		return c, core.ErrNotFound

	} else if err != nil {
		return c, core.NewError(err)
	}

	return c, nil
}
