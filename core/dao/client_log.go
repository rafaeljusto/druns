package dao

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type ClientLog struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewClientLog(sqler SQLer, ip net.IP, agent int) ClientLog {
	return ClientLog{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "client_log",
		tableFields: []string{
			"id",
			"name",
			"birthday",
			"log_id",
		},
	}
}

func (dao *ClientLog) save(c *model.Client, operation model.LogOperation) error {
	log := model.NewLog(dao.Agent, dao.IP, operation)

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
		c.Name.String(),
		c.Birthday.String(),
		log.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
