package hello

import (
	"fmt"
	"test/world"
)

func init() {
	fmt.Println("HELLO INIT")
}

func Hello() string {
	return "Hello, " + world.World()
}
