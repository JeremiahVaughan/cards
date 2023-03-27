package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type deck []string

func (d deck) print() {
	for _, card := range d {
		println(card)
	}
}

func newDeck() deck {
	cardSuites := []string{"Spades", "Hearts", "Diamonds", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	var result deck
	for _, suite := range cardSuites {
		for _, value := range cardValues {
			result = append(
				result,
				fmt.Sprintf("%s of %s", value, suite),
			)
		}
	}
	result.shuffle()
	return result
}

func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) toString() string {
	return strings.Join(d, ",")
}

func (d deck) saveToFile(filename string) error {
	err := os.WriteFile(
		filename,
		[]byte(d.toString()),
		0644,
	)
	if err != nil {
		return fmt.Errorf("error, when attempting to save the deck to the file: %v", err)
	}
	return nil
}

func newDeckFromFile(filename string) (deck, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error, when attempting to read the file: %v", err)
	}
	return strings.Split(
		string(file),
		",",
	), nil
}

func (d deck) shuffle() {
	for i := range d {
		newPosition := rand.Intn(len(d) - 1)
		d[i], d[newPosition] = d[newPosition], d[i]
	}
}
