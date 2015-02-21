package dblog

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
)

type dao struct {
	SQLer       db.SQLer
	tableName   string
	tableFields []string
}

func newDAO(sqler db.SQLer) dao {
	return dao{
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

func (dao dao) Save(dbLog *DBLog) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	dbLog.ChangedAt = time.Now()

	row := dao.SQLer.QueryRow(
		query,
		dbLog.Agent,
		dbLog.IPAddress.String(),
		dbLog.ChangedAt,
		string(dbLog.Operation),
	)

	if err := row.Scan(&dbLog.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao dao) FindById(id int64) (DBLog, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = ?",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	var dbLog DBLog
	var ipAddress string

	err := row.Scan(
		&dbLog.Id,
		&dbLog.Agent,
		&ipAddress,
		&dbLog.ChangedAt,
		&dbLog.Operation,
	)

	if err == sql.ErrNoRows {
		return dbLog, core.ErrNotFound

	} else if err != nil {
		return dbLog, core.NewError(err)
	}

	dbLog.IPAddress = net.ParseIP(ipAddress)

	return dbLog, nil
}
