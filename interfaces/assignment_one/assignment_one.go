package assignment_one

import "log"

type Triangle struct {
	Length float64
	Width  float64
}

type Square struct {
	Length float64
	Width  float64
}

type Shape interface {
	Area() float64
}

func (s Square) Area() float64 {
	return s.Length * s.Width
}

func (t Triangle) Area() float64 {
	return (t.Length * t.Width) / 2
}

func AssignmentOne() {
	t := Triangle{
		Length: 5,
		Width:  5,
	}

	printIt(t)

	s := Square{
		Length: 5,
		Width:  5,
	}

	printIt(s)
}

func printIt(s Shape) {
	log.Printf("Area of the shape is: %v", s.Area())
}
