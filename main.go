package main

import (
	"concurrency/workers"
	"fmt"
)

func main() {
	fmt.Println("Hello, Kris")
	workers.InitWork("tasks.JSON")
}
