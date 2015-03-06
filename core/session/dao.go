package session

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/errors"
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

func (dao *dao) save(s *Session) error {
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

	err := row.Scan(&s.Id)
	return errors.New(err)
}

func (dao *dao) update(s *Session) error {
	if s.revision == db.Revision(s) {
		return nil
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

	return errors.New(err)
}

func (dao *dao) findById(id int) (Session, error) {
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

	if err != nil {
		return s, errors.New(err)
	}

	s.IPAddress = net.ParseIP(ipAddress)
	s.User, err = user.NewService().FindById(dao.SQLer, userId)
	s.revision = db.Revision(s)
	return s, err
}
