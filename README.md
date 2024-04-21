# SQL Manager
Simple SQL database resource manager.

# Why?
Some SQL database drivers, like Postgres, do not manage their own resources for queries. This results in errors during high traffic when many connections are active at once. This resource manager implements a [semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore) to manage resource access and ensure that no more than the max number of connections is active at a time.

***This only designed for simple use cases. More advanced applications should look into connection pooling, sharding, etc.***

***Only tested with the "postgres" Go driver***

# Usage
Common methods are implemented on the `SQLManager` struct. Can be used as a drop-in replacement for `sql.DB` in most simple cases.

```go
package main

import (
	"fmt"
	"sync"

	"github.com/PretendoNetwork/sql-manager"
)

func main() {
	// Match these to your database setup
	driver := "postgres"
	uri := "postgres://username:password@localhost:5432/db"
	maxConnections := 4

	manager, err := sqlmanager.NewSQLManager(driver, uri, maxConnections)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	_, err = manager.Exec(`CREATE TABLE IF NOT EXISTS test (
		id bigserial PRIMARY KEY
	)`)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// Despite only 4 connections being allowed at once, all 100 calls succeed without issue
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := manager.Exec("INSERT INTO test DEFAULT VALUES")
			if err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()

	fmt.Println("Done")
}
```

# Tests
To test 4 environment variables must be set

| Name                               | Description                                                          | Example                                          |
| ---------------------------------- | -------------------------------------------------------------------- | ------------------------------------------------ |
| `SQL_MANAGER_TEST_DRIVER`          | SQL driver name.                                                     | `postgres`                                       |
| `SQL_MANAGER_TEST_URI`             | SQL database connection string.                                      | `postgres://username:password@localhost:5432/db` |
| `SQL_MANAGER_TEST_MAX_CONNECTIONS` | Number of max connections. Should match your database configuration. | 4                                                |
| `SQL_MANAGER_TEST_OPERATIONS`      | Number of operations to do at once.                                  | 100                                              |
