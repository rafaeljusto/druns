package enrollment

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/group"
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
		tableName: "enrollment",
		tableFields: []string{
			"id",
			"client_id",
			"client_group_id",
			"type",
		},
	}
}

func (dao *dao) save(e *Enrollment) error {
	if dao.agent == 0 || dao.ip == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if e.Id == 0 {
		if err := dao.insert(e); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(e); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(e, operation)
}

func (dao *dao) insert(e *Enrollment) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		e.Client.Id,
		e.Group.Id,
		e.Type,
	)

	err := row.Scan(&e.Id)
	return errors.New(err)
}

func (dao *dao) update(e *Enrollment) error {
	if e.revision == db.Revision(e) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET client_id = $1, client_group_id = $2, type = $3 WHERE id = $4",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		e.Client.Id,
		e.Group.Id,
		e.Type,
		e.Id,
	)

	return errors.New(err)
}

func (dao *dao) findById(id int) (Enrollment, error) {
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

func (dao *dao) findByClient(clientId int) (Enrollments, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE client_id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, clientId)
	if err != nil {
		return nil, errors.New(err)
	}

	var enrollments Enrollments

	for rows.Next() {
		e, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		enrollments = append(enrollments, e)
	}

	// We cannot load a composite object while we are iterating over the main
	// result, that's why we only load it after we finish the iteration

	clientService := client.NewService(dao.sqler)
	groupService := group.NewService(dao.sqler)

	for i, e := range enrollments {
		e.Client, err = clientService.FindById(e.Client.Id)
		if err != nil {
			return nil, err
		}

		e.Group, err = groupService.FindById(e.Group.Id)
		if err != nil {
			return nil, err
		}

		enrollments[i] = e
	}

	return enrollments, nil
}

func (dao *dao) findByGroup(groupId int) (Enrollments, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE client_group_id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, groupId)
	if err != nil {
		return nil, errors.New(err)
	}

	var enrollments Enrollments

	for rows.Next() {
		e, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		enrollments = append(enrollments, e)
	}

	// We cannot load a composite object while we are iterating over the main
	// result, that's why we only load it after we finish the iteration

	clientService := client.NewService(dao.sqler)
	groupService := group.NewService(dao.sqler)

	for i, e := range enrollments {
		e.Client, err = clientService.FindById(e.Client.Id)
		if err != nil {
			return nil, err
		}

		e.Group, err = groupService.FindById(e.Group.Id)
		if err != nil {
			return nil, err
		}

		enrollments[i] = e
	}

	return enrollments, nil
}

func (dao *dao) load(row db.Row, eager bool) (Enrollment, error) {
	var e Enrollment

	err := row.Scan(
		&e.Id,
		&e.Client.Id,
		&e.Group.Id,
		&e.Type,
	)

	if err != nil {
		return e, errors.New(err)
	}

	if eager {
		e.Client, err = client.NewService(dao.sqler).FindById(e.Client.Id)
		if err != nil {
			return e, err
		}

		e.Group, err = group.NewService(dao.sqler).FindById(e.Group.Id)
		if err != nil {
			return e, err
		}
	}

	e.revision = db.Revision(e)
	return e, nil
}
