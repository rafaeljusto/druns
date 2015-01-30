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
	Agent       int
	tableName   string
	tableFields []string
}

func NewUser(sqler SQLer, ip net.IP, agent int) User {
	return User{
		SQLer:     sqler,
		IP:        ip,
		Agent:     agent,
		tableName: "adm_user",
		tableFields: []string{
			"id",
			"name",
			"email",
			"password",
		},
	}
}

func (dao *User) Save(u *model.User) error {
	if dao.Agent == 0 || dao.IP == nil {
		return core.NewError(fmt.Errorf("No log information defined to persist information"))
	}

	var operation model.LogOperation

	if u.Id == 0 {
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

	userLogDAO := NewUserLog(dao.SQLer, dao.IP, dao.Agent)
	return userLogDAO.save(u, operation)
}

func (dao *User) insert(u *model.User) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		placeholders(dao.tableFields[1:]),
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
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

	query := fmt.Sprintf(
		"UPDATE %s SET name = $1, email = $2 WHERE id = $3",
		dao.tableName,
	)

	_, err := dao.SQLer.Exec(
		query,
		u.Name,
		u.Email,
		u.Id,
	)

	if err != nil {
		return core.NewError(err)
	}

	return nil
}

func (dao *User) FindById(id int) (model.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, id)

	u, err := dao.load(row)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (dao *User) FindByEmail(email string) (model.User, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, email)

	u, err := dao.load(row)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (dao *User) FindAll() ([]model.User, error) {
	// Avoid selecting the BOOTSTRAP user
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email != ''",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.SQLer.Query(query)
	if err != nil {
		return nil, core.NewError(err)
	}

	var users []model.User

	for rows.Next() {
		u, err := dao.load(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (dao *User) load(row row) (model.User, error) {
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
		"SELECT password FROM %s WHERE email = $1",
		dao.tableName,
	)

	row := dao.SQLer.QueryRow(query, email.Address)

	var base64Password string
	err := row.Scan(
		&base64Password,
	)

	if err == sql.ErrNoRows {
		return false, core.ErrNotFound

	} else if err != nil {
		return false, core.NewError(err)
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(base64Password)
	if err != nil {
		return false, core.NewError(err)
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return false, core.NewError(err)
	}

	return true, nil
}
