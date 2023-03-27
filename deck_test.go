package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_newDeck(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Happy Path",
			want: 52,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(newDeck()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDeck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deal(t *testing.T) {
	type args struct {
		d        deck
		handSize int
	}
	tests := []struct {
		name  string
		args  args
		want  deck
		want1 deck
	}{
		{
			name: "Happy Path",
			args: args{
				d:        []string{"card1", "card2", "card3", "card4"},
				handSize: 2,
			},
			want:  []string{"card1", "card2"},
			want1: []string{"card3", "card4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := deal(tt.args.d, tt.args.handSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deal() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("deal() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_saveToDeckAndNewDeckFromFile(t *testing.T) {
	filename := "_decktesting"
	err := os.Remove(filename)
	if err != nil && err.Error() != "no such file or directory" {
		t.Errorf("error, during initial file cleanup: %v", err)
	}

	cards := newDeck()
	err = cards.saveToFile(filename)
	if err != nil {
		log.Fatalf("error, when attempting to save cards to the file: %v", err)
	}

	savedDeck, err := newDeckFromFile(filename)
	if err != nil {
		log.Fatalf("error, when attempting to new deck from file: %v", err)
	}

	expected := len(cards)
	got := len(savedDeck)
	if got != expected {
		t.Errorf("error, len of deck has changed since writing to a file. Expected %d. Got: %d", expected, got)
	}
}
