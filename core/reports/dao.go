package reports

import (
	"time"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/group"
)

type dao struct {
	sqler db.SQLer
}

func newDAO(sqler db.SQLer) dao {
	return dao{
		sqler: sqler,
	}
}

func (dao *dao) incomingPerGroup(month time.Time, classValue float64) ([]Incoming, error) {
	query := `SELECT class.client_group_id, COUNT(student)
		FROM class, student, enrollment
		WHERE student.class_id = class.id 
		AND student.enrollment_id = enrollment.id
		AND EXTRACT(MONTH FROM class.begin_at) = $1
		AND EXTRACT(YEAR FROM class.begin_at) = $2
		AND enrollment.type = $3
		GROUP BY class.id`

	rows, err := dao.sqler.Query(
		query,
		month.Month(),
		month.Year(),
		enrollment.TypeRegular,
	)

	if err != nil {
		return nil, errors.New(err)
	}

	incomingsTmp := make(map[int]struct {
		realQuantity     int
		foreseenQuantity int
	})

	for rows.Next() {
		var groupId, students int

		err := rows.Scan(
			&groupId,
			&students,
		)

		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, errors.New(err)
		}

		incoming := incomingsTmp[groupId]
		incoming.realQuantity += students
		incomingsTmp[groupId] = incoming
	}

	query = ` SELECT class.client_group_id, COUNT(student)
		FROM class, student, enrollment
		WHERE student.class_id = class.id 
		AND student.enrollment_id = enrollment.id
		AND EXTRACT(MONTH FROM class.begin_at) = $1
		AND EXTRACT(YEAR FROM class.begin_at) = $2
		AND enrollment.type = $3
		GROUP BY class.id`

	rows, err = dao.sqler.Query(
		query,
		month.Month(),
		month.Year(),
		enrollment.TypeReservation,
	)

	if err != nil {
		return nil, errors.New(err)
	}

	for rows.Next() {
		var groupId, students int

		err := rows.Scan(
			&groupId,
			&students,
		)

		if err != nil {
			// TODO: Check ErrNotFound and ignore it
			return nil, errors.New(err)
		}

		incoming := incomingsTmp[groupId]
		incoming.foreseenQuantity += students
		incomingsTmp[groupId] = incoming
	}

	var incomings []Incoming
	groupService := group.NewService(dao.sqler)

	for groupId, students := range incomingsTmp {
		g, err := groupService.FindById(groupId)
		if err != nil {
			return nil, err
		}

		incomings = append(incomings, Incoming{
			Group:    g,
			Month:    month,
			Foreseen: float64(students.foreseenQuantity) * classValue,
			Value:    float64(students.realQuantity) * classValue,
		})

	}
	return incomings, nil
}
