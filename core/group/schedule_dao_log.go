package group

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

type scheduleDAOLog struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newScheduleDAOLog(sqler db.SQLer, ip net.IP, agent int) scheduleDAOLog {
	return scheduleDAOLog{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "client_group_schedule_log",
		tableFields: []string{
			"id",
			"client_group_id",
			"weekday",
			"time",
			"duration",
			"log_id",
		},
	}
}

func (dao *scheduleDAOLog) save(s *Schedule, groupID int, operation dblog.Operation) error {
	dbLog, err := dblog.NewService(dao.sqler).Create(dao.agent, dao.ip, operation)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields),
	)

	_, err = dao.sqler.Exec(
		query,
		s.Id,
		groupID,
		s.Weekday,
		s.Time,
		s.Duration,
		dbLog.Id,
	)

	if err != nil {
		return errors.New(err)
	}

	return nil
}
