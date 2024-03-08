package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

const INPUT = "Hello, OTUS!"

func main() {
	str := reverse.String(INPUT)
	fmt.Print(str)
}
