package cleantest

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

// init docker container

//Algo

// create new pool
// set port bindings for connection [set host and ip for port in container]
// set repo name and tag and other options for pulling and running docker container
// retry to connect cause we have to wait while postgres will ready for queries
// details below

func initDockerContainerForTesting() {
	//new pool
	pool, err = dockertest.NewPool("")

	if err != nil {
		panic(err)
	}
	//port binding
	var port docker.Port = "5432/tcp"
	var portBinding = docker.PortBinding{
		HostIP:   "127.0.0.1",
		HostPort: "1234",
	}
	var portBindings = map[docker.Port][]docker.PortBinding{
		port: []docker.PortBinding{portBinding},
	}

	// set configuration to run/pull&run docker container
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=mobi",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=mobidb",
			"listen_addresses = '*'",
		},

		ExposedPorts: []string{"5432/tcp"},
		PortBindings: portBindings,
	},
		// this for auto clean up
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})

	// check for error
	if err != nil {
		panic(err)
	}

	// if there is not errors add close container to deferStack
	deferStack = append(deferStack, func() {
		pool.Purge(resource)
	})

	// set expire time for closing container
	// cause some times after error not dependent on code
	// container still running
	// and we have to set expire date to close container after some time
	resource.Expire(3)

	// maxwait for waiting for upping postgres in container
	pool.MaxWait = time.Second * 5

	// retry for waiting for postgres upping
	// MaxWait time
	if err = pool.Retry(func() error {
		var err error
		db1, err := sql.Open("postgres", fmt.Sprintf(
			"postgres://mobi:test@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "mobidb"))
		if err != nil {
			return err
		}
		defer db1.Close()
		return db1.Ping()

	}); err != nil {
		// if there is an error run all defer functions and out
		RunDef()
		fmt.Printf("Could not connect to docker: %s", err)
		return
	}

	// container was ran successfully
	// nice)

}
