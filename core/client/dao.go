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
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newDAO(sqler db.SQLer, ip net.IP, agent int) dao {
	return dao{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "client",
		tableFields: []string{
			"id",
			"name",
			"birthday",
			"weight",
		},
	}
}

func (dao *dao) save(c *Client) error {
	if dao.agent == 0 || dao.ip == nil {
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

	logDAO := newDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(c, operation)
}

func (dao *dao) insert(c *Client) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		c.Name,
		c.Birthday,
		c.Weight,
	)

	err := row.Scan(&c.Id)
	return errors.New(err)
}

func (dao *dao) update(c *Client) error {
	if c.revision == db.Revision(c) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, birthday = $2, weight = $3 WHERE id = $4",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		c.Name,
		c.Birthday,
		c.Weight,
		c.Id,
	)

	return errors.New(err)
}

func (dao *dao) findById(id int) (Client, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	c, err := dao.load(row)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (dao *dao) findAll() ([]Client, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY name",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query)
	if err != nil {
		return nil, errors.New(err)
	}

	var clients []Client

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
		&c.Weight,
	)

	c.revision = db.Revision(c)
	return c, errors.New(err)
}
