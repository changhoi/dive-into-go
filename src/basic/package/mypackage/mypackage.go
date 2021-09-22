package mypackage

import "strings"

func GetName() Name {
	return Name("dive into go")
}

const VERSION = "1.0"

type Name string

func (n Name) String() string {
	return strings.ToUpper(string(n))
}
