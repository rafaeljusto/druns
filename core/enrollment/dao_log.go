package enrollment

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
)

type daoLog struct {
	SQLer       db.SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func newDAOLog(sqler db.SQLer, ip net.IP, agent int) daoLog {
	return daoLog{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "enrollment_log",
		tableFields: []string{
			"id",
			"client_id",
			"group_id",
			"type",
			"log_id",
		},
	}
}

func (dao *daoLog) save(e *Enrollment, operation dblog.Operation) error {
	dbLog, err := dblog.NewService().Create(dao.SQLer, dao.Agent, dao.IP, operation)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields),
	)

	_, err = dao.SQLer.Exec(
		query,
		e.Id,
		e.Client.Id,
		e.Group.Id,
		e.Type,
		dbLog.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
