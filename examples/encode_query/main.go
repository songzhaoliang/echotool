package main

import (
	"fmt"
	"net/http"

	"github.com/songzhaoliang/echotool"
)

type User struct {
	ID   int    `url:"id"`
	Name string `url:"name"`
}

func main() {
	user := &User{
		ID:   1,
		Name: "peter",
	}
	data, err := echotool.EncodeValues(user)
	if err != nil {
		fmt.Printf("EncodeValues error - %v\n", err)
		return
	}

	url := fmt.Sprintf("http://localhost:1323/users?%s", data.Encode())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Get error - %v\n", err)
		return
	}
	defer resp.Body.Close()

	// deal with resp.Body
}
