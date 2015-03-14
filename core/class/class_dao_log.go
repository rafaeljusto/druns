package class

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

type classDAOLog struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newClassDAOLog(sqler db.SQLer, ip net.IP, agent int) classDAOLog {
	return classDAOLog{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "class_log",
		tableFields: []string{
			"id",
			"client_group_id",
			"class_date",
			"log_id",
		},
	}
}

func (dao *classDAOLog) save(c *Class, operation dblog.Operation) error {
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
		c.Id,
		c.Group.Id,
		c.Date,
		dbLog.Id,
	)

	if err != nil {
		return errors.New(err)
	}

	return nil
}
