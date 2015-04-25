package user

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"strings"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/golang.org/x/crypto/bcrypt"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/errors"
)

var systemEmail = "system@druns.com.br"

type dao struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newDAO(sqler db.SQLer, ip net.IP, agent int) dao {
	return dao{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "adm_user",
		tableFields: []string{
			"id",
			"name",
			"email",
			"password",
		},
	}
}

func (dao *dao) save(u *User) error {
	if dao.agent == 0 || dao.ip == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if u.Id == 0 {
		if err := dao.insert(u); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(u); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(u, operation)
}

func (dao *dao) insert(u *User) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return errors.New(err)
	}

	row := dao.sqler.QueryRow(
		query,
		u.Name,
		u.Email,
		base64.StdEncoding.EncodeToString(hashedPassword),
	)

	err = row.Scan(&u.Id)
	return errors.New(err)
}

func (dao *dao) update(u *User) error {
	if u.revision == db.Revision(u) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, email = $2 WHERE id = $3",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		u.Name,
		u.Email,
		u.Id,
	)

	return errors.New(err)
}

func (dao *dao) findById(id int) (User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	u, err := dao.load(row)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (dao *dao) findByEmail(email string) (User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, email)

	u, err := dao.load(row)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (dao *dao) systemUser() (User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE name = $1 AND email = $2",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, "System", systemEmail)

	u, err := dao.load(row)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (dao *dao) findAll() ([]User, error) {
	// Avoid selecting the BOOTSTRAP user
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email != $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, systemEmail)
	if err != nil {
		return nil, errors.New(err)
	}

	var users []User

	for rows.Next() {
		u, err := dao.load(rows)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (dao *dao) load(row db.Row) (User, error) {
	var u User
	var hashedPassword string

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&hashedPassword,
	)

	u.revision = db.Revision(u)
	return u, errors.New(err)
}

func (dao *dao) verifyPassword(email mail.Address, password string) (bool, error) {
	query := fmt.Sprintf(
		"SELECT password FROM %s WHERE email = $1",
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, email.Address)

	var base64Password string
	err := row.Scan(
		&base64Password,
	)

	if err != nil {
		return false, errors.New(err)
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(base64Password)
	if err != nil {
		return false, errors.New(err)
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return false, errors.New(err)
	}

	return true, nil
}
