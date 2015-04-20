package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/rafaeljusto/druns/core/class"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/password"
	"github.com/rafaeljusto/druns/web/config"
)

var (
	Logger = log.NewLogger("system")
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>\n", os.Args[0])
		return
	}

	if err := config.LoadConfig(os.Args[1]); err != nil {
		fmt.Printf("Error loading configuration file. Details: %s\n", err)
		return
	}

	if err := initializeLogger(); err != nil {
		fmt.Printf("Error initializing logger. Details: %s\n", err)
		return
	}

	if err := initializeDatabase(); err != nil {
		Logger.Critf("Error initializing database. Details: %s", err)
		return
	}
	defer db.DB.Close()

	addr, err := localAddress()
	if err != nil {
		Logger.Errorf("Error retrieving the local address. Details: %s\n", err)
		return

	} else if addr == nil {
		Logger.Errorf("Couldn't retrieve the local address")
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		Logger.Errorf("Error creating a database transaction. Details: %s", err)
		return
	}

	groups, err := group.NewService(tx).FindAll()
	if err != nil {
		Logger.Errorf("Error retrieving groups. Details: %s", err)
		return
	}

	now := time.Now()
	twoWeeksFromNow := now.Add(7 * 24 * time.Hour)
	classService := class.NewClassService(tx)
	studentService := class.NewStudentService(tx)

	for _, group := range groups {
		classes, err := classService.FindByGroupIdBetweenDates(group.Id, now, twoWeeksFromNow)
		if err != nil {
			Logger.Errorf("Error retrieving classes. Details: %s", err)
			return
		}

		var scheduleDate time.Time
		for i := 0; i < 7; i++ {
			if d := now.Add(time.Duration(24*i) * time.Hour); d.Weekday() == group.Weekday.Weekday {
				scheduleDate = d
				break
			}
		}

		if scheduleDate.IsZero() {
			// Not on schedule yet, move on
			continue
		}

		foundClass := false
		for _, c := range classes {
			if c.BeginAt.Weekday() != group.Weekday.Weekday &&
				c.BeginAt.Hour() != group.Time.Hour() &&
				c.BeginAt.Minute() != group.Time.Minute() {

				// TODO: We should remove this class
				continue
			}

			foundClass = true
		}

		if foundClass {
			// Class already created, move on
			continue
		}

		refDate := time.Date(scheduleDate.Year(), scheduleDate.Month(), scheduleDate.Day(),
			group.Time.Hour(), group.Time.Minute(), 0, 0, time.Local)

		c := class.Class{
			Group:   group,
			BeginAt: refDate,
			EndAt:   refDate.Add(group.Duration.Duration),
		}

		if err := classService.Save(addr, systemUser.Id, &c); err != nil {
			Logger.Errorf("Error saving new class. Details: %s", err)
			return
		}

		enrollments, err := enrollment.NewService(tx).FindByGroup(group.Id)
		if err != nil {
			Logger.Errorf("Error retrieving enrollments for Group %d. Details: %s", group.Id, err)
			return
		}

		for _, e := range enrollments {
			s := class.Student{
				Enrollment: e,
			}

			if err := studentService.Save(addr, systemUser.Id, &s, c); err != nil {
				Logger.Errorf("Error saving new student. Details: %s", err)
				return
			}

			c.Students = append(c.Students, s)
		}

	}

	if err := tx.Commit(); err != nil {
		Logger.Errorf("Error commiting transaction. Details: %s", err)
		return
	}
}

func initializeLogger() error {
	logAddr := net.JoinHostPort(config.DrunsConfig.Log.Host, strconv.Itoa(config.DrunsConfig.Log.Port))
	if err := log.Connect("druns", logAddr); err != nil {
		return errors.New(err)
	}
	return nil
}

func initializeDatabase() error {
	dbPassword, err := password.Decrypt(config.DrunsConfig.Database.Password)
	if err != nil {
		return err
	}

	return db.Start(
		config.DrunsConfig.Database.Host,
		config.DrunsConfig.Database.Port,
		config.DrunsConfig.Database.User,
		dbPassword,
		config.DrunsConfig.Database.Name,
	)
}

func localAddress() (net.IP, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return nil, err
	}

	if len(addrs) > 0 {
		return net.ParseIP(addrs[0]), nil
	}

	return nil, nil
}
