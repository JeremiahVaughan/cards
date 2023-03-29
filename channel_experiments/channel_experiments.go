package channel_experiments

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func ChannelExperiments() {
	//testDoSomething()
	testBufferCapacity()
}

func testBufferCapacity() {
	waitGroup := sync.WaitGroup{}
	errChan := make(chan error, 1)
	bufferChan := make(chan string, 2)
	paths := []string{"test1", "test2", "test3", "test4", "test5", "test6"}
	waitGroup.Add(len(paths))
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Generating UI artifacts...")
	for _, pathI := range paths {
		go func(pathY string) {
			defer waitGroup.Done()
			e := doBufferCapacity(cancelCtx, bufferChan, pathY)
			if e != nil {
				// Ensuring we only receive one error in the buffer at a time so that if context is canceled nothing else is stuck in the buffer, thus errors stuck in buffer can still flow to a receiving operator even though the context is canceled.
				select {
				case errChan <- e:
					cancel()
					log.Printf("context canceled due to error: %v", e)
				default:
				}

			}
		}(pathI)
	}

	time.Sleep(time.Second)
	log.Printf("work being collected: %v", <-bufferChan)
	log.Printf("work being collected 2: %v", <-bufferChan)
	time.Sleep(4 * time.Second)
	log.Printf("work being collected a short time later: %v", <-bufferChan)
	time.Sleep(11 * time.Second)
	log.Printf("work being collected much later, but still in the buffer therefore made it past the context cancelation: %v", <-bufferChan)
	waitGroup.Wait()
}

func testDoSomething() {
	waitGroup := sync.WaitGroup{}
	errChan := make(chan error, 1)
	paths := []string{"test1", "test2", "test3"}
	waitGroup.Add(len(paths))
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Generating UI artifacts...")
	for _, pathI := range paths {
		go func(pathY string) {
			defer waitGroup.Done()
			e := doSomething(cancelCtx, pathY)
			if e != nil {
				// Ensuring we only cancel the context 1 time since buffer has a size of 1. However, there is no harm in canceling more than once
				select {
				case errChan <- e:
					cancel()
				default:
				}

			}
		}(pathI)
	}
	waitGroup.Wait()
}

func doSomething(ctx context.Context, path string) error {
	if path == "test2" {
		log.Printf("%s is having trouble ...", path)
		time.Sleep(3 * time.Second)
		return fmt.Errorf("someError")
	}
	if path == "test3" {
		time.Sleep(7 * time.Second)
	}
	log.Printf("path is doing work: %s", path)
	select {
	case <-ctx.Done():
		log.Printf("%s has just been told to abandon his work", path)
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("error, timeout waiting for channel to unblock")
	}
}

func doBufferCapacity(ctx context.Context, c chan string, path string) error {
	//if path == "test2" {
	//	log.Printf("%s is taking a while ...", path)
	//	time.Sleep(3 * time.Second)
	//	return fmt.Errorf("someError")
	//}
	log.Printf("path is doing good work: %s", path)
	select {
	case c <- path:
		log.Printf("added %s to buffer", path)
	case <-ctx.Done():
		log.Printf("%s has just been told to abandon his work", path)
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return fmt.Errorf("error, timeout waiting for channel to unblock")
	}
	return nil
}
