package dao

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type Group struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewGroup(sqler SQLer, ip net.IP, agent int) Group {
	return Group{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "client_group",
		tableFields: []string{
			"id",
			"name",
			"weekday",
			"time",
			"duration",
			"type",
			"capacity",
		},
	}
}

func (dao *Group) Save(g *model.Group) error {
	if dao.Agent == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
	}

	var operation model.LogOperation

	if g.Id == 0 {
		if err := dao.insert(g); err != nil {
			return err
		}

		operation = model.LogOperationCreation

	} else {
		if err := dao.update(g); err != nil {
			return err
		}

		operation = model.LogOperationUpdate
	}

	GroupLogDAO := NewGroupLog(dao.SQLer, dao.IP, dao.Agent)
	return GroupLogDAO.save(g, operation)
}

func (dao *Group) insert(g *model.Group) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields[1:]),
	)

	row := dao.SQLer.QueryRow(
		query,
		g.Name.String(),
		g.Weekday.String(),
		g.Time.String(),
		g.Duration.String(),
		g.Type.String(),
		g.Capacity,
	)

	if err := row.Scan(&g.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *Group) update(g *model.Group) error {
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
		g.Name.String(),
		g.Weekday.String(),
		g.Time.String(),
		g.Duration.String(),
		g.Type.String(),
		g.Capacity,
		g.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *Group) FindById(id int) (model.Group, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	g, err := dao.load(row)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (dao *Group) FindAll() (model.Groups, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, core.NewError(err)
	}

	var groups model.Groups

	for rows.Next() {
		g, err := dao.load(rows)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		groups = append(groups, g)
	}

	return groups, nil
}

func (dao *Group) load(row row) (model.Group, error) {
	var g model.Group
	var name string
	var weekday, duration, groupType []byte

	err := row.Scan(
		&g.Id,
		&name,
		&weekday,
		&g.Time.Time,
		&duration,
		&groupType,
		&g.Capacity,
	)

	if err == sql.ErrNoRows {
		return g, core.ErrNotFound

	} else if err != nil {
		return g, core.NewError(err)
	}

	g.Name.Set(name)

	if err := g.Weekday.UnmarshalText(weekday); err != nil {
		return g, err
	}
	if err := g.Duration.UnmarshalText(duration); err != nil {
		return g, err
	}
	if err := g.Type.UnmarshalText(groupType); err != nil {
		return g, err
	}

	return g, nil
}
