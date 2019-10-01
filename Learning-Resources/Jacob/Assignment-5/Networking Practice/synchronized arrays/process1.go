package main

import (
	"fmt"
	"os"
)

func main() {
	a := os.Getenv("PORT")
	fmt.Println("hello from process", a, ".go")
}
