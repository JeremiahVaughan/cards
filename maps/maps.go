package maps

import "fmt"

func Maps() {
	colors := map[string]string{
		"red":   "#FF0000",
		"blue":  "#0000FF",
		"green": "#00FF00",
	}

	// not valid
	//var colors2 map[string]string
	//colors2["red"] = "#FF0000"
	//colors2["blue"] = "#0000FF"
	//colors2["green"] = "#00FF00"

	colors3 := make(map[string]string)
	colors3["red"] = "#FF0000"
	colors3["blue"] = "#0000FF"
	colors3["green"] = "#00FF00"

	fmt.Printf("colors:\n")
	printMap(colors)
	//fmt.Printf("colors2:\n")
	//printMap(colors2)
	fmt.Printf("colors3:\n")
	printMap(colors3)

	delete(colors, "green")
	fmt.Printf("colors: %+v\n", colors)

}

func printMap(c map[string]string) {
	for k, v := range c {
		fmt.Printf("Hex code for %s is %s\n", k, v)
	}
}
