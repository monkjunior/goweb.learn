package main

import "fmt"

func main() {
	messages := make(chan string)
	go func() { messages <- "hello" }()
	go func() { messages <- "ping" }()
	msg := <-messages
	msg2 := <-messages
	fmt.Println(msg)
	fmt.Println(msg2)
}
