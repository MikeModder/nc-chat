package server

import (
	"fmt"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
	createTable = `CREATE TABLE 'users' (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	'Name'	TEXT,
	'Password'	TEXT,
	'Mode'	INTEGER
	);`
)

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "nc-chat.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}

	fmt.Println("[info] connected to database")
	return nil
}

// CreateUser creates a new user and stores their information in the database
func CreateUser(user, password string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into users(name, password, mode) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	_, err = stmt.Exec(user, string(hash), ModeUser)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CheckPassword(in, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(in))
}

func GetUserByID(id int) (User, error) {
	tx, err := db.Begin()
	if err != nil {
		return User{}, err
	}

	stmt, err := tx.Prepare("select id, name, password, mode from users where id = ?")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return User{}, err
	}

	var u User
	// Do we really need a for loop here? There should only ever be one result...
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.Mode)
		if err != nil {
			return User{}, err
		}
	}

	return u, nil
}

func GetUserByName(name string) (User, error) {
	tx, err := db.Begin()
	if err != nil {
		return User{}, err
	}

	stmt, err := tx.Prepare("select id, name, password, mode from users where name = ?")
	if err != nil {
		return User{}, err
	}

	rows, err := stmt.Query(name)
	if err != nil {
		return User{}, err
	}

	var u User
	for rows.Next() {
		rows.Scan(&u.ID, &u.Name, &u.Password, &u.Mode)
	}

	return u, nil
}

//func UserExists(name string) // Probably not needed
