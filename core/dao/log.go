package dao

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type Log struct {
	SQLer       SQLer
	tableName   string
	tableFields []string
}

func NewLog(sqler SQLer) Log {
	return Log{
		SQLer:     sqler,
		tableName: "log",
		tableFields: []string{
			"id",
			"agent",
			"ip_address",
			"changed_at",
			"operation",
		},
	}
}

func (dao Log) Save(log *model.Log) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields[1:]),
	)

	log.ChangedAt = time.Now()

	row := dao.SQLer.QueryRow(
		query,
		log.Agent,
		log.IPAddress.String(),
		log.ChangedAt,
		string(log.Operation),
	)

	if err := row.Scan(&log.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao Log) FindById(id int64) (model.Log, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = ?",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	var log model.Log
	var ipAddress string

	err := row.Scan(
		&log.Id,
		&log.Agent,
		&ipAddress,
		&log.ChangedAt,
		&log.Operation,
	)

	if err == sql.ErrNoRows {
		return log, core.ErrNotFound

	} else if err != nil {
		return log, core.NewError(err)
	}

	log.IPAddress = net.ParseIP(ipAddress)

	return log, nil
}
