package logconsent

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	name     string
	password string
	email    string
}

type Database struct {
	Dbuser string
	Dbpw   string
	Dbname string
	Dbport string
	Dbhost string
	Db     *sql.DB
}

func (d *Database) Connect() {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		d.Dbhost, d.Dbuser, d.Dbpw, d.Dbname, d.Dbport)

	// Establish database connection.
	db, err := sql.Open("postgres", dbinfo)
	Handle(err)

	d.Db = db
}

func (d *Database) Validate(mail, pw string) bool {
	var user User
	err := d.Db.QueryRow("SELECT name, password, email FROM users WHERE email=$1",
		mail).Scan(&user.name, &user.password, &user.email)
	Handle(err)

	if user.password != pw {
		return false
	}

	return true

}
