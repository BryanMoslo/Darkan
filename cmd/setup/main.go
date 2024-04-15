package main

import (
	"fmt"

	"darkan/internal"
	"darkan/internal/migrations"

	"github.com/leapkit/core/db"
)

func main() {
	conn, err := internal.DB()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.RunMigrations(migrations.All, conn)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("✅ Migrations ran successfully")
}
