package main

import "log"

func main() {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, number := range numbers {
		if number%2 == 0 {
			log.Printf("the number: %d is even", number)
		} else {
			log.Printf("the number: %d is odd", number)
		}
	}
}
