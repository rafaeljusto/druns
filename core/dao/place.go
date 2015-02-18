package dao

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type Place struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewPlace(sqler SQLer, ip net.IP, agent int) Place {
	return Place{
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

func (dao *Place) Save(p *model.Place) error {
	if dao.Agent == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
	}

	var operation model.LogOperation

	if p.Id == 0 {
		if err := dao.insert(p); err != nil {
			return err
		}

		operation = model.LogOperationCreation

	} else {
		if err := dao.update(p); err != nil {
			return err
		}

		operation = model.LogOperationUpdate
	}

	placeLogDAO := NewPlaceLog(dao.SQLer, dao.IP, dao.Agent)
	return placeLogDAO.save(p, operation)
}

func (dao *Place) insert(p *model.Place) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields[1:]),
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

func (dao *Place) update(p *model.Place) error {
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

func (dao *Place) FindById(id int) (model.Place, error) {
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

func (dao *Place) FindAll() (model.Places, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, core.NewError(err)
	}

	var places model.Places

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

func (dao *Place) load(row row) (model.Place, error) {
	var p model.Place

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
