package class

import (
	"net"
	"time"

	"github.com/rafaeljusto/druns/core/db"
)

type ClassService struct {
	sqler db.SQLer
}

func NewClassService(sqler db.SQLer) ClassService {
	return ClassService{sqler}
}

func (s ClassService) Save(ip net.IP, agent int, c *Class) error {
	dao := newClassDAO(s.sqler, ip, agent)
	return dao.save(c)
}

func (s ClassService) FindById(id int) (Class, error) {
	dao := newClassDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s ClassService) FindAll() ([]Class, error) {
	dao := newClassDAO(s.sqler, nil, 0)
	return dao.findAll()
}

func (s ClassService) FindByGroupIdBetweenDates(groupId int, begin, end time.Time) ([]Class, error) {
	dao := newClassDAO(s.sqler, nil, 0)
	return dao.findByGroupIdBetweenDates(groupId, begin, end)
}

type StudentService struct {
	sqler db.SQLer
}

func NewStudentService(sqler db.SQLer) StudentService {
	return StudentService{sqler}
}

func (s StudentService) Save(ip net.IP, agent int, student *Student, c Class) error {
	dao := newStudentDAO(s.sqler, ip, agent)
	return dao.save(student, c)
}

func (s StudentService) FindById(id int) (Student, error) {
	dao := newStudentDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s StudentService) FindAll() ([]Student, error) {
	dao := newStudentDAO(s.sqler, nil, 0)
	return dao.findAll()
}
