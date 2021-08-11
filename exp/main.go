package main

import (
	"fmt"

	"github.com/monkjunior/goweb.learn/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
