package payment

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/client"
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
		tableName: "payment",
		tableFields: []string{
			"id",
			"client_id",
			"status",
			"expires_at",
			"value",
		},
	}
}

func (dao *dao) save(p *Payment) error {
	if dao.agent == 0 || dao.ip == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if p.Id == 0 {
		if err := dao.insert(p); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(p); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(p, operation)
}

func (dao *dao) insert(p *Payment) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		p.Client.Id,
		p.Status,
		p.ExpiresAt,
		p.Value,
	)

	err := row.Scan(&p.Id)
	return errors.New(err)
}

func (dao *dao) update(p *Payment) error {
	if p.revision == db.Revision(p) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET client_id = $1, status = $2, expires_at = $3, value = $4 WHERE id = $5",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		p.Client.Id,
		p.Status,
		p.ExpiresAt,
		p.Value,
	)

	return errors.New(err)
}

func (dao *dao) findById(id int) (Payment, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	e, err := dao.load(row, true)
	if err != nil {
		return e, err
	}

	return e, nil
}

func (dao *dao) findByClient(clientId int) ([]Payment, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE client_id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, clientId)
	if err != nil {
		return nil, errors.New(err)
	}

	var payments []Payment

	for rows.Next() {
		p, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		payments = append(payments, p)
	}

	// We cannot load a composite object while we are iterating over the main
	// result, that's why we only load it after we finish the iteration

	clientService := client.NewService(dao.sqler)

	for i, p := range payments {
		p.Client, err = clientService.FindById(p.Client.Id)
		if err != nil {
			return nil, err
		}

		payments[i] = p
	}

	return payments, nil
}

func (dao *dao) load(row db.Row, eager bool) (Payment, error) {
	var p Payment

	err := row.Scan(
		&p.Id,
		&p.Client.Id,
		&p.Status,
		&p.ExpiresAt,
		&p.Value,
	)

	if err != nil {
		return p, errors.New(err)
	}

	if eager {
		p.Client, err = client.NewService(dao.sqler).FindById(p.Client.Id)
		if err != nil {
			return p, err
		}
	}

	p.revision = db.Revision(p)
	return p, nil
}
