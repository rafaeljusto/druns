package place

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
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
		tableName: "place",
		tableFields: []string{
			"id",
			"name",
			"address",
		},
	}
}

func (dao *dao) Save(p *Place) error {
	if dao.Agent == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
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

	logDAO := newDAOLog(dao.SQLer, dao.IP, dao.Agent)
	return logDAO.save(p, operation)
}

func (dao *dao) insert(p *Place) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.SQLer.QueryRow(
		query,
		p.Name,
		p.Address,
	)

	if err := row.Scan(&p.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *dao) update(p *Place) error {
	if lastPlace, err := dao.FindById(p.Id); err == nil && lastPlace.Equal(*p) {
		// Nothing changed
		return nil

	} else if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, address = $2 WHERE id = $3",
		dao.tableName,
	)

	_, err := dao.SQLer.Exec(
		query,
		p.Name,
		p.Address,
		p.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *dao) FindById(id int) (Place, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	p, err := dao.load(row)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (dao *dao) FindAll() (Places, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, core.NewError(err)
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

	if err == sql.ErrNoRows {
		return p, core.ErrNotFound

	} else if err != nil {
		return p, core.NewError(err)
	}

	return p, nil
}
