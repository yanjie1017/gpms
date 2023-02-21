package main

import (
	"fmt"
	service "src/services"
)

func main() {
	fmt.Println("Hello World")
	service.GeneratePassword("a", "b", "c", "d", 1)
}
