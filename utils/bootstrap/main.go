package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"time"

	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/password"
	"github.com/rafaeljusto/druns/core/place"
	"github.com/rafaeljusto/druns/core/types"
	"github.com/rafaeljusto/druns/core/user"
	"github.com/rafaeljusto/druns/web/config"
)

var (
	configPath string
	name       string
	email      string
	pwd        string
	test       bool
)

func init() {
	flag.StringVar(&configPath, "config", "", "configuration file path of the web server")
	flag.StringVar(&name, "name", "", "administrator's name")
	flag.StringVar(&email, "email", "", "administrator's e-mail")
	flag.StringVar(&pwd, "password", "", "administrator's password")
	flag.BoolVar(&test, "test", false, "test environment")
}

func main() {
	flag.Parse()

	if !checkInputs() {
		usage()
		os.Exit(1)
	}

	if err := config.LoadConfig(configPath); err != nil {
		fmt.Printf("Error loading configuration file. Details: %s\n", err)
		os.Exit(2)
	}

	if err := initializeDatabase(); err != nil {
		fmt.Printf("Error initializing database. Details: %s\n", err)
		os.Exit(3)
	}
	defer db.DB.Close()

	e, err := types.NewEmail(email)
	if err != nil {
		fmt.Printf("Invalid e-mail. Details: %s\n", err)
		os.Exit(4)
	}

	u := user.User{
		Name:     types.NewName(name),
		Email:    e,
		Password: pwd,
	}

	addr, err := localAddress()
	if err != nil {
		fmt.Printf("Error retrieving the local address. Details: %s\n", err)
		os.Exit(5)

	} else if addr == nil {
		fmt.Printf("Couldn't retrieve the local address")
		os.Exit(6)
	}

	tx, err := db.DB.Begin()
	if err != nil {
		fmt.Printf("Error starting database transaction. Details: %s\n", err)
		os.Exit(7)
	}

	// Bootstrap user doesn't have password to avoid using it in running enviroment
	row := tx.QueryRow("INSERT INTO adm_user(id, name, email) " +
		"VALUES (DEFAULT, 'System', 'system@druns.com.br') RETURNING id")

	var id int
	if err := row.Scan(&id); err != nil {
		fmt.Printf("Error inserting bootstrap user. Details: %s\n", err)
		os.Exit(8)
	}

	if users, err := user.NewService(tx).FindAll(); err != nil {
		fmt.Printf("Error retrieving users. Details: %s\n", err)
		os.Exit(9)

	} else if len(users) > 0 {
		fmt.Println("Database already initialized")
		return
	}

	if err := user.NewService(tx).Save(addr, id, &u); err != nil {
		fmt.Printf("Error saving the user. Details: %s\n", err)
		os.Exit(10)
	}

	if test {
		if !generateClients(tx, addr, id) {
			os.Exit(11)
		}

		if !generatePlaces(tx, addr, id) {
			os.Exit(12)
		}

		if !generateGroups(tx, addr, id) {
			os.Exit(13)
		}

		if !generateEnrollments(tx, addr, id) {
			os.Exit(14)
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Error commiting database transaction. Details: %s\n", err)
		os.Exit(15)
	}

	fmt.Println("Bootstrap runned successfully")
}

func usage() {
	fmt.Printf("Usage: %s <-config 'path'> <-email 'email'> "+
		"<-name 'name'> <-password 'password'>\n", os.Args[0])
	flag.PrintDefaults()
}

func checkInputs() bool {
	ok := true

	if len(configPath) == 0 {
		fmt.Println("Configuration path not informed!")
		ok = false
	}

	if len(name) == 0 {
		fmt.Println("Name not informed!")
		ok = false
	}

	if len(email) == 0 {
		fmt.Println("E-mail not informed!")
		ok = false
	}

	if _, err := mail.ParseAddress(email); err != nil {
		fmt.Printf("Invalid e-mail. Details: %s\n", err)
		os.Exit(4)
	}

	if len(pwd) == 0 {
		fmt.Println("Password not informed!")
		ok = false
	}

	return ok
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

func generateClients(tx db.Transaction, ip net.IP, agent int) bool {
	clients := []client.Client{
		{Name: types.NewName("Erika Santiago"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Lula Knight"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Geraldine Howard"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Maureen Lawrence"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ted Mason"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Katrina Richardson"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Angelo Douglas"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Lorena Snyder"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Rafael Byrd"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Preston Soto"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Max Francis"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ernestine Jennings"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Roosevelt Romero"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Wallace Chapman"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Heidi Ramirez"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Julio Owens"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Allison Freeman"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Cheryl Gutierrez"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Marcia Gomez"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Donnie Norton"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Carol Lynch"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Thelma Bowers"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Claire Ramsey"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Nina Stone"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jose Bryant"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Leigh Hogan"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Daryl Fox"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Violet Williams"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Johanna Lucas"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Edna Clayton"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Tim Blake"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jay Morgan"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Georgia Torres"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Leland Martin"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ismael Tate"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Janis Simpson"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Aaron Rodgers"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Warren Wheeler"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Sally Cohen"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Elvira James"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Fredrick Dean"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Horace Tyler"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jeanne Nunez"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Earnest Chandler"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Milton Sutton"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Myra Horton"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Camille Wong"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Derrick Powers"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Gladys Higgins"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Hector Gross"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Todd Joseph"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Roberta Mills"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Faith Saunders"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Rick Carter"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Gina Perez"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Clifford Lopez"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Herman Barnes"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jesse Perkins"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Marie Becker"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Leslie Farmer"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Phillip Matthews"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Henry Reynolds"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Pablo Santos"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jeffrey Rose"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Rene Wood"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Stewart Curtis"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Olga Munoz"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Kari Tucker"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Bradley Henry"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Allen Garcia"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ethel Salazar"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jesus Conner"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Joan Berry"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Rosemarie Allen"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ella Ortega"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Mathew Stokes"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Madeline Park"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Joanna Arnold"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Robin Lowe"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Phil Goodwin"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Francis Swanson"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Maggie Mckinney"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Eric Lane"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Travis Edwards"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Nicholas Butler"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Bonnie Leonard"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Melody Luna"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Laurie Bush"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Brian Palmer"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ann Bennett"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Dixie Meyer"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Ken Moody"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Kurt Ray"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Wilbert Cooper"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Lindsey Marsh"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Josefina Hale"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Jennifer Ross"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Mark Nguyen"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Brittany Stephens"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: types.NewName("Gerardo Stanley"), Birthday: types.NewDate(time.Unix(rand.Int63n(478000000), 0))},
	}

	for _, c := range clients {
		if err := client.NewService(tx).Save(ip, agent, &c); err != nil {
			fmt.Printf("Error creating a client. Details: %s\n", err)
			return false
		}
	}

	return true
}

func generatePlaces(tx db.Transaction, ip net.IP, agent int) bool {
	places := []place.Place{
		{Name: types.NewName("Campolim Park"), Address: types.NewAddress("200 Domingos Julio Av.")},
		{Name: types.NewName("Water Park"), Address: types.NewAddress("3560 Dom Aguirre Av.")},
	}

	for _, p := range places {
		if err := place.NewService(tx).Save(ip, agent, &p); err != nil {
			fmt.Printf("Error creating a place. Details: %s\n", err)
			return false
		}
	}

	return true
}

func generateGroups(tx db.Transaction, ip net.IP, agent int) bool {
	groupTypes := []string{group.TypeOnce, group.TypeWeekley}

	for i := 1; i <= 10; i++ {
		var groupType group.Type
		groupType.Set(groupTypes[rand.Intn(len(groupTypes))])

		g := group.Group{
			Name:  types.NewName(fmt.Sprintf("Group %d", i)),
			Place: place.Place{Id: rand.Intn(2) + 1},
			Schedules: []group.Schedule{
				{
					Weekday:  types.NewWeekday(time.Weekday(rand.Intn(5) + 1)),
					Time:     types.NewTime(time.Date(0, 0, 0, rand.Intn(17)+6, rand.Intn(60), 0, 0, time.Local)),
					Duration: types.NewDuration(time.Duration(rand.Intn(120)) * time.Minute),
				},
			},
			Type:     groupType,
			Capacity: rand.Intn(20) + 10,
		}

		if err := group.NewService(tx).Save(ip, agent, &g); err != nil {
			fmt.Printf("Error creating a group. Details: %s\n", err)
			return false
		}
	}

	return true
}

func generateEnrollments(tx db.Transaction, ip net.IP, agent int) bool {
	enrollmentTypes := []string{enrollment.TypeRegular, enrollment.TypeReplacement, enrollment.TypeReservation}
	cache := make(map[int][]int)

	for i := 1; i <= 200; i++ {
		var enrollmentType enrollment.Type
		enrollmentType.Set(enrollmentTypes[rand.Intn(len(enrollmentTypes))])

		var clientId, groupId int

		for {
			clientId = rand.Intn(100) + 1
			groupId = rand.Intn(10) + 1

			created := false
			for _, id := range cache[groupId] {
				if id == clientId {
					created = true
					break
				}
			}

			if !created {
				cache[groupId] = append(cache[groupId], clientId)
				break
			}
		}

		e := enrollment.Enrollment{
			Type:   enrollmentType,
			Client: client.Client{Id: clientId},
			Group:  group.Group{Id: groupId},
		}

		if err := enrollment.NewService(tx).Save(ip, agent, &e); err != nil {
			fmt.Printf("Error creating an enrollment. Details: %s\n", err)
			return false
		}
	}

	return true
}
