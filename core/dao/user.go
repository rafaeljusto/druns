package dao

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"strings"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	SQLer       SQLer
	IP          net.IP
	Handle      string
	tableName   string
	tableFields []string
}

func NewUser(sqler SQLer, ip net.IP, handle string) User {
	return User{
		SQLer:     sqler,
		IP:        ip,
		Handle:    handle,
		tableName: "user",
		tableFields: []string{
			"id",
			"name",
			"email",
			"password",
		},
	}
}

func (dao *User) Save(u *model.User) error {
	if len(dao.Handle) == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
	}

	var operation model.LogOperation

	if u.Id > 0 {
		if err := dao.insert(u); err != nil {
			return err
		}

		operation = model.LogOperationCreation

	} else {
		if err := dao.update(u); err != nil {
			return err
		}

		operation = model.LogOperationUpdate
	}

	userLogDAO := NewUserLog(dao.SQLer, dao.IP, dao.Handle)
	return userLogDAO.save(u, operation)
}

func (dao *User) insert(u *model.User) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields[1:]),
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MaxCost)
	if err != nil {
		return core.NewError(err)
	}

	row := dao.SQLer.QueryRow(
		query,
		u.Name,
		u.Email,
		base64.StdEncoding.EncodeToString(hashedPassword),
	)

	if err := row.Scan(&u.Id); err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *User) update(u *model.User) error {
	if lastUser, err := dao.FindById(u.Id); err == nil && lastUser.Equal(*u) {
		// Nothing changed
		return nil

	} else if err != nil {
		return err
	}

	query := `
		UPDATE user
		SET name = ?,
			email = ?
		WHERE id = ?
	`

	_, err := dao.SQLer.Exec(
		query,
		u.Name,
		u.Email,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *User) FindById(id int) (model.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = ?",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	var u model.User
	var hashedPassword string

	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&hashedPassword,
	)

	if err == sql.ErrNoRows {
		return u, core.ErrNotFound

	} else if err != nil {
		return u, core.NewError(err)
	}

	return u, nil
}

func (dao *User) VerifyPassword(email mail.Address, password string) (bool, error) {
	query := fmt.Sprintf(
		"SELECT password FROM %s WHERE email = ?",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, email.Address)

	var hashedPassword string
	err := row.Scan(
		&hashedPassword,
	)

	if err == sql.ErrNoRows {
		return false, core.ErrNotFound

	} else if err != nil {
		return false, core.NewError(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}
