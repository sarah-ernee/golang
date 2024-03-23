package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [function]")
		return
	}

	runSpecificFunc(os.Args[1])
}

func runSpecificFunc(function string) {
	switch function {
	case "main":
		variableSet()	
	case "second":
		library()
	default:
		fmt.Println("Unknown function:", function)
	}
}

func variableSet() {
	// shorthand inside function for var declaration
	// var := value

	var name string = "Mario"
	age := 20

	// formatted string %_ = format specifier
	fmt.Printf("%v is %v years old \n", name, age)
	fmt.Printf("His name in quotes is %q \n", name) 
	fmt.Printf("Age variable type is %T \n", age)
	fmt.Printf("His credit score in two decimal points is %0.2f \n", 255.5242433)

	saved_str := fmt.Sprintf("%v is %v years old \n", name, age)
	fmt.Println(saved_str)

	// arrays and slicing
	scores := []int{100, 75, 60}
	scores[2] = 25

	// unlike in PY JS, arr append does not overwrite OG array but returns a new one 
	// initialize scores array to append() 
	scores = append(scores,  85)
	fmt.Println(scores, len(scores))

	first_range := scores[0:1] 
	second_range := scores[1:]
	fmt.Printf("Range restricted items from scores is %v \n", first_range)
	fmt.Printf("Item 2 onwards till end of array is %v \n", second_range)

	third_range := scores[:3]
	fmt.Printf("Start of array up to till exclusive of item 4 of the array is %v", third_range)
}

func library() {
	// string package
	greeting := "hello there friends!"
	fmt.Println(strings.Contains(greeting, "hello!"))
	fmt.Println("Equivalent of REGEX: ", strings.ReplaceAll(greeting, "hello", "hi"))
	fmt.Println("Greeting in CAPS:", strings.ToUpper(greeting))
	fmt.Println("Position where 'll' occurs:", strings.Index(greeting, "ll"))
	fmt.Println("Split string by space", strings.Split(greeting, ""))

	fmt.Println("Original greeting is still unchanged:", greeting)

	// sort package 
	ages := []int{45, 29, 28, 18, 15, 30}
	sort.Ints(ages)
	index := sort.SearchInts(ages, 30)
	fmt.Println("Index for item 30 after sorting is", index)




}

