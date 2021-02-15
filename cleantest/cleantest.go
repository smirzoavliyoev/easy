package cleantest

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/ory/dockertest/v3"
)

// bad code
// im thinking on how to refactor this and make it more flexible and
// simple for using

//type CleanTest struct {
//	pool     *dockertest.Pool
//	resource *dockertest.Resource
//	dbConn   *sql.DB
//	driver   database.Driver
//	mg       *migrate.Migrate
//
//	//
//
//	deferStack []func()
//	err        error
//}

// so...lets start
// for simplicity i implemented it so
// next step it will be implemented  as module and use options to let user
// inject his own configuration

// not detailed algo:
// 1. we have to run docker container (if there is no docker container it  will pull it
// from docker hub)

// 2. we have to init migrations

// also all ran containers and migrations must close after operation
// thats wy i used deferStack (i know its bad code but for clarity)

// ok
// algo in details:
// in test main
// first of all user will run New() function that will init container and run migrations
// user will use defer function CLose  to close migrations and stop docker container
// between New() and Close test will pass
// example in user side :
//	cleantest.New()
//	defer cleantest.Close()
// for more see example in transactionRepo test

// during running New function,  deferStack slice sums all defer functions
// why?
// cause of Close()
// if New function has an error RunDef will run to close migrations and stop docker container
// if there is not errors then we have not to stop container and down migrations
// in user side user will use Close() function to close all containers and migrations

// vars
var (
	//container vars
	pool     *dockertest.Pool
	resource *dockertest.Resource

	//database conn
	dbConn *sql.DB

	// driver connections and migration
	driver database.Driver
	mg     *migrate.Migrate

	//

	// sum of defer functions
	deferStack []func()
	err        error
)

// new clean test
func New() {
	//run docker container
	// details in function
	initDockerContainerForTesting()
	// run migrations
	// details in function
	initMigrations()

}

// close all used services
func Close() {
	RunDef()
}

// run defer functions
func RunDef() {
	for i := len(deferStack) - 1; i >= 0; i-- {
		deferStack[i]()
	}
}
