package main

import "fmt"

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
	return result
}

func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) toString() string {
	[]string(d)
}
