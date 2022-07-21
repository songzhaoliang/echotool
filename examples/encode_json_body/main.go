package main

import (
	"fmt"
	"net/http"

	"github.com/songzhaoliang/echotool"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	user := &User{
		ID:   1,
		Name: "peter",
	}
	buffer, err := echotool.EncodeJSON(user)
	if err != nil {
		fmt.Printf("EncodeJSON error - %v\n", err)
		return
	}
	defer echotool.ReleaseBuffer(buffer)

	url := "http://localhost:1323/users"
	req, err := http.NewRequest(http.MethodPost, url, buffer)
	if err != nil {
		fmt.Printf("NewRequest error - %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Do error - %v\n", err)
		return
	}
	defer resp.Body.Close()

	// deal with resp.Body
}
