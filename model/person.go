package model

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	ContactInfo
}

func (p *Person) UpdateName(newFirstName string) {
	p.FirstName = newFirstName
}

func (p *Person) Print() {
	fmt.Printf("%+v", p)
}
