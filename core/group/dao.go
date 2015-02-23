package group

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/place"
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
		tableName: "client_group",
		tableFields: []string{
			"id",
			"name",
			"place_id",
			"weekday",
			"time",
			"duration",
			"type",
			"capacity",
		},
	}
}

func (dao *dao) Save(g *Group) error {
	if dao.Agent == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if g.Id == 0 {
		if err := dao.insert(g); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(g); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newDAOLog(dao.SQLer, dao.IP, dao.Agent)
	return logDAO.save(g, operation)
}

func (dao *dao) insert(g *Group) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.SQLer.QueryRow(
		query,
		g.Name,
		g.Place.Id,
		g.Weekday,
		g.Time,
		g.Duration,
		g.Type,
		g.Capacity,
	)

	if err := row.Scan(&g.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *dao) update(g *Group) error {
	if lastGroup, err := dao.FindById(g.Id); err == nil && lastGroup.Equal(*g) {
		// Nothing changed
		return nil

	} else if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, weekday = $2, time = $3, duration = $4, type = $5, capacity = $6 WHERE id = $7",
		dao.tableName,
	)

	_, err := dao.SQLer.Exec(
		query,
		g.Name,
		g.Place.Id,
		g.Weekday,
		g.Time,
		g.Duration,
		g.Type,
		g.Capacity,
		g.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *dao) FindById(id int) (Group, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	g, err := dao.load(row, true)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (dao *dao) FindAll() (Groups, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, core.NewError(err)
	}

	var groups Groups

	for rows.Next() {
		g, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		groups = append(groups, g)
	}

	// We cannot load a composite object while we are iterating over the main
	// result, that's why we only load it after we finish the iteration
	placeService := place.NewService()
	for i, g := range groups {
		if g.Place, err = placeService.FindById(dao.SQLer, g.Place.Id); err != nil {
			return nil, err
		}
		groups[i] = g
	}

	return groups, nil
}

func (dao *dao) load(row db.Row, eager bool) (Group, error) {
	var g Group

	err := row.Scan(
		&g.Id,
		&g.Name,
		&g.Place.Id,
		&g.Weekday,
		&g.Time,
		&g.Duration,
		&g.Type,
		&g.Capacity,
	)

	if err == sql.ErrNoRows {
		return g, core.ErrNotFound

	} else if err != nil {
		return g, core.NewError(err)
	}

	if eager {
		g.Place, err = place.NewService().FindById(dao.SQLer, g.Place.Id)
	}

	return g, err
}