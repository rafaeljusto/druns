package group

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

type scheduleDAO struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newScheduleDAO(sqler db.SQLer, ip net.IP, agent int) scheduleDAO {
	return scheduleDAO{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "client_group_schedule",
		tableFields: []string{
			"id",
			"client_group_id",
			"weekday",
			"time",
			"duration",
		},
	}
}

func (dao *scheduleDAO) save(s *Schedule, groupID int) error {
	if dao.agent == 0 || dao.ip == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if s.Id == 0 {
		if err := dao.insert(s, groupID); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(s); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newScheduleDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(s, groupID, operation)
}

func (dao *scheduleDAO) insert(s *Schedule, groupID int) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		groupID,
		s.Weekday,
		s.Time,
		s.Duration,
	)

	err := row.Scan(&s.Id)
	return errors.New(err)
}

func (dao *scheduleDAO) update(s *Schedule) error {
	if s.revision == db.Revision(s) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET weekday = $1, time = $2, duration = $3 WHERE id = $4",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		s.Weekday,
		s.Time,
		s.Duration,
		s.Id,
	)

	return errors.New(err)
}

func (dao *scheduleDAO) findById(id int) (Schedule, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	s, err := dao.load(row)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (dao *scheduleDAO) findByGroup(groupID int) ([]Schedule, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE client_group_id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, groupID)
	if err != nil {
		return nil, errors.New(err)
	}

	var schedules []Schedule

	for rows.Next() {
		s, err := dao.load(rows)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, s)
	}

	return schedules, nil
}

func (dao *scheduleDAO) load(row db.Row) (Schedule, error) {
	var s Schedule
	var groupID int

	err := row.Scan(
		&s.Id,
		&groupID,
		&s.Weekday,
		&s.Time,
		&s.Duration,
	)

	if err != nil {
		return s, errors.New(err)
	}

	s.revision = db.Revision(s)
	return s, err
}
