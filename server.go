package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/6210530017/assessment/config"

	_ "github.com/lib/pq"
)

func Setup(url string) (*sql.DB, func()) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	teardown := func() {
		db.Close()
	}

	return db, teardown
}

func main() {
	config := config.NewConfig()

	db, teardown := Setup(config.DB_url)
	defer teardown()

	fmt.Printf("%#v", db)
}
