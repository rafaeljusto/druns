package place

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
		tableName: "place",
		tableFields: []string{
			"id",
			"name",
			"address",
		},
	}
}

func (dao *dao) save(p *Place) error {
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

func (dao *dao) insert(p *Place) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		p.Name,
		p.Address,
	)

	err := row.Scan(&p.Id)
	return errors.New(err)
}

func (dao *dao) update(p *Place) error {
	if p.revision == db.Revision(p) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, address = $2 WHERE id = $3",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		p.Name,
		p.Address,
		p.Id,
	)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (dao *dao) findById(id int) (Place, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	p, err := dao.load(row)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (dao *dao) findAll() (Places, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query)
	if err != nil {
		return nil, errors.New(err)
	}

	var places Places

	for rows.Next() {
		p, err := dao.load(rows)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		places = append(places, p)
	}

	return places, nil
}

func (dao *dao) load(row db.Row) (Place, error) {
	var p Place

	err := row.Scan(
		&p.Id,
		&p.Name,
		&p.Address,
	)

	p.revision = db.Revision(p)
	return p, errors.New(err)
}
