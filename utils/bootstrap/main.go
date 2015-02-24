package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/password"
	"github.com/rafaeljusto/druns/core/place"
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

	e, err := core.NewEmail(email)
	if err != nil {
		fmt.Printf("Invalid e-mail. Details: %s\n", err)
		os.Exit(4)
	}

	u := user.User{
		Name:     core.NewName(name),
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
	row := tx.QueryRow("INSERT INTO adm_user(id, name) " +
		"VALUES (DEFAULT, 'BOOTSTRAP') RETURNING id")

	var id int
	if err := row.Scan(&id); err != nil {
		fmt.Printf("Error inserting bootstrap user. Details: %s\n", err)
		os.Exit(8)
	}

	if users, err := user.NewService().FindAll(tx); err != nil {
		fmt.Printf("Error retrieving users. Details: %s\n", err)
		os.Exit(9)

	} else if len(users) > 0 {
		fmt.Println("Database already initialized")
		return
	}

	if err := user.NewService().Save(tx, addr, id, &u); err != nil {
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
		{Name: core.NewName("Erika Santiago"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Lula Knight"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Geraldine Howard"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Maureen Lawrence"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ted Mason"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Katrina Richardson"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Angelo Douglas"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Lorena Snyder"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Rafael Byrd"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Preston Soto"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Max Francis"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ernestine Jennings"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Roosevelt Romero"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Wallace Chapman"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Heidi Ramirez"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Julio Owens"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Allison Freeman"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Cheryl Gutierrez"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Marcia Gomez"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Donnie Norton"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Carol Lynch"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Thelma Bowers"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Claire Ramsey"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Nina Stone"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jose Bryant"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Leigh Hogan"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Daryl Fox"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Violet Williams"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Johanna Lucas"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Edna Clayton"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Tim Blake"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jay Morgan"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Georgia Torres"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Leland Martin"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ismael Tate"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Janis Simpson"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Aaron Rodgers"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Warren Wheeler"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Sally Cohen"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Elvira James"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Fredrick Dean"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Horace Tyler"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jeanne Nunez"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Earnest Chandler"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Milton Sutton"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Myra Horton"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Camille Wong"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Derrick Powers"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Gladys Higgins"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Hector Gross"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Todd Joseph"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Roberta Mills"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Faith Saunders"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Rick Carter"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Gina Perez"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Clifford Lopez"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Herman Barnes"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jesse Perkins"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Marie Becker"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Leslie Farmer"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Phillip Matthews"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Henry Reynolds"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Pablo Santos"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jeffrey Rose"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Rene Wood"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Stewart Curtis"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Olga Munoz"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Kari Tucker"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Bradley Henry"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Allen Garcia"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ethel Salazar"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jesus Conner"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Joan Berry"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Rosemarie Allen"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ella Ortega"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Mathew Stokes"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Madeline Park"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Joanna Arnold"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Robin Lowe"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Phil Goodwin"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Francis Swanson"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Maggie Mckinney"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Eric Lane"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Travis Edwards"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Nicholas Butler"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Bonnie Leonard"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Melody Luna"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Laurie Bush"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Brian Palmer"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ann Bennett"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Dixie Meyer"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Ken Moody"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Kurt Ray"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Wilbert Cooper"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Lindsey Marsh"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Josefina Hale"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Jennifer Ross"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Mark Nguyen"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Brittany Stephens"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
		{Name: core.NewName("Gerardo Stanley"), Birthday: core.NewDate(time.Unix(rand.Int63n(478000000), 0))},
	}

	for _, c := range clients {
		if err := client.NewService().Save(tx, ip, agent, &c); err != nil {
			fmt.Printf("Error creating a client. Details: %s\n", err)
			return false
		}
	}

	return true
}

func generatePlaces(tx db.Transaction, ip net.IP, agent int) bool {
	places := []place.Place{
		{Name: core.NewName("Campolim Park"), Address: core.NewAddress("200 Domingos Julio Av.")},
		{Name: core.NewName("Water Park"), Address: core.NewAddress("3560 Dom Aguirre Av.")},
	}

	for _, p := range places {
		if err := place.NewService().Save(tx, ip, agent, &p); err != nil {
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
			Name:     core.NewName(fmt.Sprintf("Group %d", i)),
			Place:    place.Place{Id: rand.Intn(2) + 1},
			Weekday:  core.NewWeekday(time.Weekday(rand.Intn(7))),
			Time:     core.NewTime(time.Date(0, 0, 0, rand.Intn(24), rand.Intn(60), 0, 0, time.UTC)),
			Duration: core.NewDuration(time.Duration(rand.Intn(120)) * time.Minute),
			Type:     groupType,
			Capacity: rand.Intn(20) + 10,
		}

		if err := group.NewService().Save(tx, ip, agent, &g); err != nil {
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

		if err := enrollment.NewService().Save(tx, ip, agent, &e); err != nil {
			fmt.Printf("Error creating an enrollment. Details: %s\n", err)
			return false
		}
	}

	return true
}
