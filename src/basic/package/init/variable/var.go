package variable

import "fmt"

func init() {
	fmt.Println("VAR INIT1")
}

func OverVar() string {
	fmt.Println("OverVar Call")
	return "OverVar"
}

func UnderVar() string {
	fmt.Println("UnderVar Call")
	return "UnderVar"
}

func init() {
	fmt.Println("VAR INIT2")
}
