package payment

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

type daoLog struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newDAOLog(sqler db.SQLer, ip net.IP, agent int) daoLog {
	return daoLog{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "payment_log",
		tableFields: []string{
			"id",
			"client_id",
			"status",
			"expires_at",
			"value",
			"log_id",
		},
	}
}

func (dao *daoLog) save(p *Payment, operation dblog.Operation) error {
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
		p.Id,
		p.Client.Id,
		p.Status,
		p.ExpiresAt,
		p.Value,
		dbLog.Id,
	)

	if err != nil {
		return errors.New(err)
	}

	return nil
}
