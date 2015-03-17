package class

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

type studentDAOLog struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newStudentDAOLog(sqler db.SQLer, ip net.IP, agent int) studentDAOLog {
	return studentDAOLog{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "student_log",
		tableFields: []string{
			"id",
			"class_id",
			"enrollment_id",
			"attended",
			"log_id",
		},
	}
}

func (dao *studentDAOLog) save(s *Student, c Class, operation dblog.Operation) error {
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
		c.Id,
		s.Enrollment.Id,
		s.Attended,
		dbLog.Id,
	)

	if err != nil {
		return errors.New(err)
	}

	return nil
}
