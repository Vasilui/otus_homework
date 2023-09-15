package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	out := Reverse("Hello, OTUS!")
	fmt.Println(out)
}

func Reverse(in string) string {
	return reverse.String(in)
}
