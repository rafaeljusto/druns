package group

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
		tableName: "client_group_log",
		tableFields: []string{
			"id",
			"name",
			"weekday",
			"time",
			"duration",
			"type",
			"capacity",
			"log_id",
		},
	}
}

func (dao *daoLog) save(g *Group, operation dblog.Operation) error {
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
		g.Id,
		g.Name,
		g.Weekday,
		g.Time,
		g.Duration,
		g.Type,
		g.Capacity,
		dbLog.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
