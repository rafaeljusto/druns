package session

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/user"
)

type dao struct {
	SQLer       db.SQLer
	tableName   string
	tableFields []string
}

func newDAO(sqler db.SQLer) dao {
	return dao{
		SQLer:     sqler,
		tableName: "session",
		tableFields: []string{
			"id",
			"adm_user_id",
			"ip_address",
			"created_at",
			"last_access_at",
		},
	}
}

func (dao *dao) Save(s *Session) error {
	if s.Id == 0 {
		return dao.insert(s)
	} else {
		return dao.update(s)
	}
}

func (dao *dao) insert(s *Session) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.SQLer.QueryRow(
		query,
		s.User.Id,
		s.IPAddress.String(),
		s.CreatedAt,
		s.LastAccessAt,
	)

	if err := row.Scan(&s.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *dao) update(s *Session) error {
	if lastSession, err := dao.FindById(s.Id); err == nil && lastSession.Equal(*s) {
		// Nothing changed
		return nil

	} else if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"UPDATE %s SET last_access_at = $1 WHERE id = $2",
		dao.tableName,
	)

	_, err := dao.SQLer.Exec(
		query,
		s.LastAccessAt,
		s.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *dao) FindById(id int) (Session, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	var s Session
	var userId int
	var ipAddress string

	err := row.Scan(
		&s.Id,
		&userId,
		&ipAddress,
		&s.CreatedAt,
		&s.LastAccessAt,
	)

	if err == sql.ErrNoRows {
		return s, core.ErrNotFound

	} else if err != nil {
		return s, core.NewError(err)
	}

	s.IPAddress = net.ParseIP(ipAddress)
	s.User, err = user.NewService().FindById(dao.SQLer, userId)
	return s, err
}
