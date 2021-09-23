package main

import (
	"fmt"

	"test/hello"
	"test/variable"
)

var overVar = variable.OverVar()

func init() {
	fmt.Println("MAIN INIT")
}

func main() {
	fmt.Print("MAIN FUNC / ")
	fmt.Println(hello.Hello())
}

var underVar = variable.UnderVar()
