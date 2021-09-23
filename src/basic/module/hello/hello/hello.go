package hello

import (
	"fmt"

	"example.com/world"
)

func Hello() string {
	txt := world.GetWorld()
	return fmt.Sprint("안녕,", txt)
}
