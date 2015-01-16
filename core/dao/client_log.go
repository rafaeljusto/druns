package dao

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type ClientLog struct {
	SQLer       SQLer
	IP          net.IP
	Handle      string
	tableName   string
	tableFields []string
}

func NewClientLog(sqler SQLer, ip net.IP, handle string) ClientLog {
	return ClientLog{
		SQLer:     sqler,
		IP:        ip,
		Handle:    handle,
		tableName: "client_log",
		tableFields: []string{
			"id",
			"name",
			"birthday",
			"id_log",
		},
	}
}

func (dao *ClientLog) save(c *model.Client, operation model.LogOperation) error {
	log := model.Log{
		Handle:    dao.Handle,
		IPAddress: dao.IP,
		ChangedAt: time.Now(),
		Operation: operation,
	}

	logDAO := NewLog(dao.SQLer)
	if err := logDAO.Save(&log); err != nil {
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields),
	)

	_, err := dao.SQLer.Exec(
		query,
		c.Id,
		c.Name,
		c.Birthday,
		log.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
