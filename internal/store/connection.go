package store

import (
	"database/sql"
	"fmt"

	// a blank import is necessary here
	_ "github.com/go-sql-driver/mysql"
)

// Config has database settings which are set by environmental variables.
type Config struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	HOST     string `env:"DB_HOST"`
	PORT     int    `env:"DB_PORT" envDefault:"3306"`
}

// Open checks the connection to the database and returns DB struct, cleanup function, and error.
func Open(c *Config) (*sql.DB, func(), error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.User, c.Password, c.HOST, c.PORT, c.Name))
	if err != nil {
		return nil, func() {}, err
	}
	err = db.Ping()
	if err != nil {
		return nil, func() {}, err
	}
	return db, func() { _ = db.Close() }, nil
}
