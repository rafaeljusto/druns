package class

import (
	"fmt"
	"net"
	"strings"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/dblog"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/errors"
)

type studentDAO struct {
	sqler       db.SQLer
	ip          net.IP
	agent       int
	tableName   string
	tableFields []string
}

func newStudentDAO(sqler db.SQLer, ip net.IP, agent int) studentDAO {
	return studentDAO{
		sqler:     sqler,
		ip:        ip,
		agent:     agent,
		tableName: "student",
		tableFields: []string{
			"id",
			"class_id",
			"enrollment_id",
			"attended",
		},
	}
}

func (dao *studentDAO) save(s *Student, c Class) error {
	if dao.agent == 0 || dao.ip == nil {
		return errors.New(fmt.Errorf("No log information defined to persist information"))
	}

	var operation dblog.Operation

	if s.Id == 0 {
		if err := dao.insert(s, c); err != nil {
			return err
		}

		operation = dblog.OperationCreation

	} else {
		if err := dao.update(s); err != nil {
			return err
		}

		operation = dblog.OperationUpdate
	}

	logDAO := newStudentDAOLog(dao.sqler, dao.ip, dao.agent)
	return logDAO.save(s, c, operation)
}

func (dao *studentDAO) insert(s *Student, c Class) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (DEFAULT, %s) RETURNING id",
		dao.tableName,
		strings.Join(dao.tableFields, ", "),
		db.Placeholders(dao.tableFields[1:]),
	)

	row := dao.sqler.QueryRow(
		query,
		c.Id,
		s.Enrollment.Id,
		s.Attended,
	)

	err := row.Scan(&s.Id)
	return errors.New(err)
}

func (dao *studentDAO) update(s *Student) error {
	if s.revision == db.Revision(s) {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE %s SET attended = $1 WHERE id = $2",
		dao.tableName,
	)

	_, err := dao.sqler.Exec(
		query,
		s.Attended,
		s.Id,
	)

	return errors.New(err)
}

func (dao *studentDAO) findById(id int) (Student, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	row := dao.sqler.QueryRow(query, id)

	s, err := dao.load(row, true)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (dao *studentDAO) findByClass(classId int) ([]Student, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE class_id = $1",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query, classId)
	if err != nil {
		return nil, errors.New(err)
	}

	var students []Student

	for rows.Next() {
		s, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		students = append(students, s)
	}

	enrollmentService := enrollment.NewService(dao.sqler)
	for i, s := range students {
		students[i].Enrollment, err = enrollmentService.FindById(s.Enrollment.Id)
		if err != nil {
			return nil, err
		}
	}

	return students, nil
}

func (dao *studentDAO) findAll() ([]Student, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(dao.tableFields, ", "),
		dao.tableName,
	)

	rows, err := dao.sqler.Query(query)
	if err != nil {
		return nil, errors.New(err)
	}

	var students []Student

	for rows.Next() {
		s, err := dao.load(rows, false)
		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, err
		}

		students = append(students, s)
	}

	return students, nil
}

func (dao *studentDAO) load(row db.Row, eager bool) (Student, error) {
	var s Student
	var classId int

	err := row.Scan(
		&s.Id,
		&classId,
		&s.Enrollment.Id,
		&s.Attended,
	)

	if err != nil {
		return s, errors.New(err)
	}

	if eager {
		s.Enrollment, err = enrollment.NewService(dao.sqler).FindById(s.Enrollment.Id)
	}

	s.revision = db.Revision(s)
	return s, err
}
