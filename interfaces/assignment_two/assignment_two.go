package assignment_two

import (
	"log"
	"os"
)

func AssignmentTwo(fileName string) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error, when reading file: %v", err)
	}

	log.Println(string(file))
}
