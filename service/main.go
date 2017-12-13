package main

import (
	"args"
	"database"
	"server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.InitializeDB(*args.DB)
	server.Start()
}
