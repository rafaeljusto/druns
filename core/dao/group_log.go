package dao

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type GroupLog struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewGroupLog(sqler SQLer, ip net.IP, agent int) GroupLog {
	return GroupLog{
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

func (dao *GroupLog) save(g *model.Group, operation model.LogOperation) error {
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
		g.Id,
		g.Name.String(),
		g.Weekday.String(),
		g.Time.String(),
		g.Duration.String(),
		g.Type.String(),
		g.Capacity,
		log.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
