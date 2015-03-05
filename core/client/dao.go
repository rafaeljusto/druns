package client

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

type dao struct {
	SQLer       db.SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func newDAO(sqler db.SQLer, ip net.IP, agent int) dao {
	return dao{
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

func (dao *dao) Save(c *Client) error {
	if dao.Agent == 0 || dao.IP == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if c.Id == 0 {
		if err := dao.insert(c); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(c); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newDAOLog(dao.SQLer, dao.IP, dao.Agent)
	return logDAO.save(c, operation)
}

func (dao *dao) insert(c *Client) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.SQLer.QueryRow(
		query,
		c.Name,
		c.Birthday,
	)

	if err := row.Scan(&c.Id); err != nil {
		return errors.New(err)
	}

	return nil
}

func (dao *dao) update(c *Client) error {
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
		return errors.New(err)
	}

	return nil
}

func (dao *dao) FindById(id int) (Client, error) {
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

func (dao *dao) FindAll() (Clients, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY name",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, errors.New(err)
	}

	var clients Clients

	for rows.Next() {
		c, err := dao.load(rows)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		clients = append(clients, c)
	}

	return clients, nil
}

func (dao *dao) load(row db.Row) (Client, error) {
	var c Client

	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Birthday,
	)

	return c, errors.New(err)
}
