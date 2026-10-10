package main

import (
	"fmt"
	"regexp"
)

func main() {
	var quantityPattern = regexp.MustCompile(`^([1-9]{1}[0-9]{1})|^([0-9])$`)
	fmt.Println(quantityPattern.Match([]byte("1")))
}
