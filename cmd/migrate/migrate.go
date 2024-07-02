package main

import (
	"fmt"
	"lamoi/internal"
	"lamoi/internal/migrations"

	"github.com/leapkit/leapkit/core/db"
)

func main() {
	conn, err := internal.DB()
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
	}

	err = db.RunMigrations(migrations.All, conn)
	if err != nil {
		fmt.Println("Error running migrations: ", err)
	}
}
