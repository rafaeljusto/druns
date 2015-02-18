package dao

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type PlaceLog struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewPlaceLog(sqler SQLer, ip net.IP, agent int) PlaceLog {
	return PlaceLog{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "place_log",
		tableFields: []string{
			"id",
			"name",
			"address",
			"log_id",
		},
	}
}

func (dao *PlaceLog) save(p *model.Place, operation model.LogOperation) error {
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
		p.Id,
		p.Name,
		p.Address,
		log.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
