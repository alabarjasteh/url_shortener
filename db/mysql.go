package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alabarjasteh/url-shortener/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	db *sql.DB
}

func NewMySql(c *config.Config) Interface {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Mysql.User, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Dbname)
	db, err := sql.Open(c.Mysql.Driver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic("Cannot connect to db")
	}
	return &MySql{
		db: db,
	}
}

func (m *MySql) Save(shortlink, originallink string) error {
	var exists int
	err := m.db.QueryRow("SELECT EXISTS( SELECT * FROM pastes WHERE BINARY shortlink = ?)", shortlink).Scan(&exists)
	if err != nil {
		return err
	}

	formatedTime := time.Now().Format("2006-01-02 15:04:05")
	if exists == 1 {
		updateTimeStamp, err := m.db.Prepare("UPDATE pastes SET created_at=? WHERE shortlink=?")
		if err != nil {
			return err
		}
		defer updateTimeStamp.Close()
		updateTimeStamp.Exec(formatedTime, shortlink)
		return nil
	}
	stmt, err := m.db.Prepare("INSERT INTO pastes(shortlink, originallink, created_at) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(shortlink, originallink, formatedTime)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySql) Load(shortlink string) (string, error) {
	var originalLink string
	err := m.db.QueryRow("SELECT originallink FROM pastes WHERE shortlink = ?", shortlink).Scan(&originalLink)
	if err != nil {
		return "", ErrNotFound
	}
	return originalLink, nil
}
