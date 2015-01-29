package dao

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
)

type UserLog struct {
	SQLer       SQLer
	IP          net.IP
	Handle      string
	tableName   string
	tableFields []string
}

func NewUserLog(sqler SQLer, ip net.IP, handle string) UserLog {
	return UserLog{
		SQLer:     sqler,
		IP:        ip,
		Handle:    handle,
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
		u.Id,
		u.Name,
		u.Email,
		log.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}