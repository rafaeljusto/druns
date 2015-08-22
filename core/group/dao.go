package group

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/place"
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
		tableName: "client_group",
		tableFields: []string{
			"id",
			"name",
			"place_id",
			"type",
			"capacity",
		},
	}
}

func (dao *dao) save(g *Group) error {
	if dao.agent == 0 || dao.ip == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
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

	logDAO := newDAOLog(dao.sqler, dao.ip, dao.agent)
	if err := logDAO.save(g, operation); err != nil {
		return err
	}

	scheduleDAO := newScheduleDAO(dao.sqler, dao.ip, dao.agent)

	for i, s := range g.Schedules {
		if err := scheduleDAO.save(&s, g.Id); err != nil {
			return err
		}

		g.Schedules[i] = s
	}

	return nil
}

func (dao *dao) insert(g *Group) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		g.Name,
		g.Place.Id,
		g.Type,
		g.Capacity,
	)

	err := row.Scan(&g.Id)
	return errors.New(err)
}

func (dao *dao) update(g *Group) error {
	if g.revision == db.Revision(g) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, weekday = $2, time = $3, duration = $4, type = $5, capacity = $6 WHERE id = $7",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		g.Name,
		g.Place.Id,
		g.Type,
		g.Capacity,
		g.Id,
	)

	return errors.New(err)
}

func (dao *dao) findById(id int) (Group, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	g, err := dao.load(row, true)
	if err != nil {
		return g, err
	}

	scheduleDAO := newScheduleDAO(dao.sqler, dao.ip, dao.agent)
	if g.Schedules, err = scheduleDAO.findByGroup(g.Id); err != nil {
		return g, err
	}

	return g, nil
}

func (dao *dao) findAll() ([]Group, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY name",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query)
	if err != nil {
		return nil, errors.New(err)
	}

	var groups []Group

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
	scheduleDAO := newScheduleDAO(dao.sqler, dao.ip, dao.agent)
	placeService := place.NewService(dao.sqler)

	for i, g := range groups {
		if g.Schedules, err = scheduleDAO.findByGroup(g.Id); err != nil {
			return nil, err
		}

		if g.Place, err = placeService.FindById(g.Place.Id); err != nil {
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
		&g.Type,
		&g.Capacity,
	)

	if err != nil {
		return g, errors.New(err)
	}

	if eager {
		g.Place, err = place.NewService(dao.sqler).FindById(g.Place.Id)
	}

	g.revision = db.Revision(g)
	return g, err
}
