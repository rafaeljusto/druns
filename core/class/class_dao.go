package class

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/group"
)

type classDAO struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newClassDAO(sqler db.SQLer, ip net.IP, agent int) classDAO {
	return classDAO{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "class",
		tableFields: []string{
			"id",
			"client_group_id",
			"class_date",
		},
	}
}

func (dao *classDAO) save(c *Class) error {
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

	logDAO := newClassDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(c, operation)
}

func (dao *classDAO) insert(c *Class) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		c.Group.Id,
		c.Date,
	)

	err := row.Scan(&c.Id)
	return errors.New(err)
}

func (dao *classDAO) update(c *Class) error {
	if c.revision == db.Revision(c) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET class_date = $1 WHERE id = $2",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		c.Date,
		c.Id,
	)

	return errors.New(err)
}

func (dao *classDAO) findById(id int) (Class, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	c, err := dao.load(row, true)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (dao *classDAO) findAll() ([]Class, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query)
	if err != nil {
		return nil, errors.New(err)
	}

	var classes []Class

	for rows.Next() {
		c, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		classes = append(classes, c)
	}

	return classes, nil
}

func (dao *classDAO) findByGroupIdBetweenDates(groupId int, begin, end time.Time) ([]Class, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE client_group_id = $1 AND class_date >= $2 AND class_date <= $3",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, groupId, begin, end)
	if err != nil {
		return nil, errors.New(err)
	}

	var classes []Class

	for rows.Next() {
		c, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		classes = append(classes, c)
	}

	return classes, nil
}

func (dao *classDAO) load(row db.Row, eager bool) (Class, error) {
	var c Class

	err := row.Scan(
		&c.Id,
		&c.Group.Id,
		&c.Date,
	)

	if err != nil {
		return c, errors.New(err)
	}

	if eager {
		c.Group, err = group.NewService(dao.sqler).FindById(c.Group.Id)
	}

	c.revision = db.Revision(c)
	return c, err
}
