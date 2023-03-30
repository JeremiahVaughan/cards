package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// In this example, we use the QueryContext function from the database/sql package, which takes a context as its first argument. This allows us to cancel the database query if an error occurs in one of the other goroutines, similar to the example with API calls. Make sure to replace the connStr variable with your actual database connection string.
func queryDB(ctx context.Context, db *sql.DB, query string, results chan<- string, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		errors <- fmt.Errorf("Error executing query (%s): %v", query, err)
		return
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			errors <- fmt.Errorf("Error scanning result for query (%s): %v", query, err)
			return
		}
		results <- fmt.Sprintf("Query result (%s): %s", query, result)
	}

	if err := rows.Err(); err != nil {
		errors <- fmt.Errorf("Error in rows for query (%s): %v", query, err)
		return
	}
}

func main() {
	connStr := "user=dbuser dbname=mydb password=mypassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	query1 := "SELECT data FROM table1 LIMIT 1;"
	query2 := "SELECT data FROM table2 LIMIT 1;"

	results := make(chan string, 2)
	errors := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go queryDB(ctx, db, query1, results, errors, &wg)
	go queryDB(ctx, db, query2, results, errors, &wg)

	select {
	case result := <-results:
		fmt.Println(result)
	case err := <-errors:
		fmt.Fprintln(os.Stderr, err)
		cancel()
	}

	wg.Wait()
	close(results)
	close(errors)

	for {
		select {
		case result, ok := <-results:
			if !ok {
				results = nil
			} else {
				fmt.Println(result)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		if results == nil && errors == nil {
			break
		}
	}
}
