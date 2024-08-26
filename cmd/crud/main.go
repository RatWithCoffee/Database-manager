package main

import (
	"crud/internal"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DSN = "root:love@tcp(localhost:3306)/crud?charset=utf8"
)

func main() {
	db, err := sql.Open("mysql", DSN)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	handler, err := internal.NewDbExplorer(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", handler)
}
