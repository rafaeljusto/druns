package dao

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type UserLog struct {
	SQLer       SQLer
	IP          net.IP
	Agent       int
	tableName   string
	tableFields []string
}

func NewUserLog(sqler SQLer, ip net.IP, agent int) UserLog {
	return UserLog{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "adm_user_log",
		tableFields: []string{
			"id",
			"name",
			"email",
			"log_id",
		},
	}
}

func (dao *UserLog) save(u *model.User, operation model.LogOperation) error {
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
		u.Id,
		u.Name.String(),
		u.Email.String(),
		log.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}
