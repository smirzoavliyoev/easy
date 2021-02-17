package easy

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"time"
	// _ "github.com/golang-migrate/migrate/v4/source/file"
	// _ "github.com/lib/pq"
)

// init migrations
// cause migration packe works with sql driver thats why i do not use pgxpool
// but dont worry it just for running migrations
// i will close it after running

// algo:
// connect to postgres in container
// add close connection to deferStack
// set configurations
// after set connection driver to migration
// run migration
// add Down migration to defer stack

// more here

func initMigrations(migrattionFileName string) {

	// get connection to db in docker container
	dbConn, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://mobi:test@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "mobidb"))

	// if there is an error run all defer functions
	if err != nil {
		RunDef()
		panic(err)
	}

	// add close db to deferSrack
	deferStack = append(deferStack, func() {
		dbConn.Close()
	})

	//set configuration
	// its not important
	dbConn.SetConnMaxIdleTime(time.Minute * 5)
	dbConn.SetConnMaxLifetime(time.Minute * 5)
	dbConn.SetMaxIdleConns(0)
	dbConn.SetMaxOpenConns(0)

	err = dbConn.Ping()

	if err != nil {
		RunDef()
		panic(err)
	}

	// dirver for migrations
	driver, err = postgres.WithInstance(dbConn, &postgres.Config{
		StatementTimeout: time.Minute * 5,
	})

	if err != nil {
		RunDef()
		panic(err)
	}

	// close driver in deferStack
	deferStack = append(deferStack, func() {
		driver.Close()
	})

	// new migration

	if migrattionFileName == "" {
		migrattionFileName = "file://./../../../migrations/"
	}

	mg, err = migrate.NewWithDatabaseInstance(
		migrattionFileName,
		"mobidb",
		driver,
	)

	if err != nil {
		RunDef()
		panic(err)
	}

	// up migrations
	err = mg.Up()

	if err != nil {
		RunDef()
		panic(err)
	}

	// add Down migrations in deferStack
	deferStack = append(deferStack, func() {
		mg.Down()
	})

	//migrations was init-ed and upped
}
