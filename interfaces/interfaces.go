package interfaces

import "fmt"

type bot interface {
	getGreeting() string
}

type englishBot struct{}
type spanishBot struct{}

func Interfaces() {
	eb := englishBot{}
	sb := spanishBot{}

	printGreeting(eb)
	printGreeting(sb)
}

// Pretending these functions getGreeting() are very specific custom logic. Even though the name is the same for each struct, the implementations are very concrete, not generic.
func (englishBot) getGreeting() string {
	return "hello there"
}

func (sb spanishBot) getGreeting() string {
	return "hola amigo"
}

// Pretending that both printGreeting() functions are very generic implementations, not concrete.
func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}
