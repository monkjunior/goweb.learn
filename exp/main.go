package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	toHash := []byte("this is my string to hash")
	h := hmac.New(sha256.New, []byte("my-secret-key"))
	_, err := h.Write(toHash)
	if err != nil {
		panic(err)
	}
	b := h.Sum(nil)
	fmt.Println(b)
	h.Reset()
	h = hmac.New(sha256.New, []byte("new-secret-key"))
	_, err = h.Write(toHash)
	if err != nil {
		panic(err)
	}
	b = h.Sum(nil)
	fmt.Println(b)
}
