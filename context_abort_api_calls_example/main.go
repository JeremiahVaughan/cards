package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

// In this example, we use a context.Context to cancel the other goroutines when an error occurs. We create a cancelable context with context.WithCancel(context.Background()). When an error is received in the errors channel, we call the cancel() function to signal the other goroutines to stop their execution. The makeAPICall function takes the context and listens for the cancel signal by using req = req.WithContext(ctx) when making the HTTP request. The http.DefaultClient.Do(req) function will return an error if the context is canceled, causing the goroutine to exit early.
func makeAPICall(ctx context.Context, url string, results chan<- string, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errors <- fmt.Errorf("Error creating request for %s: %v", url, err)
		return
	}

	req = req.WithContext(ctx)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		errors <- fmt.Errorf("Error calling %s: %v", url, err)
		return
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errors <- fmt.Errorf("Error reading %s response: %v", url, err)
		return
	}

	results <- fmt.Sprintf("%s response: %s", url, string(data))
}

func main() {
	api1 := "https://api.example1.com/data"
	api2 := "https://api.example2.com/data"

	results := make(chan string, 2)
	errors := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go makeAPICall(ctx, api1, results, errors, &wg)
	go makeAPICall(ctx, api2, results, errors, &wg)

	// Wait for either results or errors
	select {
	case result := <-results:
		fmt.Println(result)
	case err := <-errors:
		fmt.Fprintln(os.Stderr, err)
		cancel() // Cancel the other goroutines
	}

	wg.Wait()
	close(results)
	close(errors)

	// Process remaining results and errors
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
