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
	// We are closing channels afterwards for two reasons:
	//    1. Indicate that no more values will be sent on the channel: When you close a channel, you signal to the receivers that there won't be any more values sent on the channel. This can be useful when you have multiple goroutines sending values to a channel and want to signal the receiver that all the goroutines have finished their work.
	//    2. Allow for iteration using the range keyword: When you close a channel, you can use the range keyword to iterate over the remaining values in the channel until the channel is empty. If the channel is not closed, the range loop will block, waiting for more values, and the loop will not terminate.
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
