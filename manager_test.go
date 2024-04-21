package sqlmanager

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

func TestManager(t *testing.T) {
	maxConnections, err := strconv.Atoi(os.Getenv("SQL_MANAGER_TEST_MAX_CONNECTIONS"))
	if err != nil {
		t.Fatal(err)
	}

	operations, err := strconv.Atoi(os.Getenv("SQL_MANAGER_TEST_OPERATIONS"))
	if err != nil {
		t.Fatal(err)
	}

	manager, err := NewSQLManager(os.Getenv("SQL_MANAGER_TEST_DRIVER"), os.Getenv("SQL_MANAGER_TEST_URI"), int64(maxConnections))
	if err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	_, err = manager.Exec(`CREATE TABLE IF NOT EXISTS test (
		id bigserial PRIMARY KEY
	)`)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < operations; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := manager.Exec("INSERT INTO test DEFAULT VALUES")
			if err != nil {
				panic(err) // * Can't use t.Fatal outside of the tests goroutine
			}
		}()
	}

	wg.Wait()

	fmt.Println("Done")
}
