package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	uri := os.Args[1]

	data, err := http.Head(uri) //r.RequestURI)
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range data.Header {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}

	fmt.Println(data.Header["Etag"])

	//etag := data.Header["Etag"][0]

	req, err := http.NewRequest("Get", uri, nil)
	req.Header.Add("If-None-Match", "666-58fae1a4eb300")
	fmt.Println(req)

	c := http.Client{}

	resp, err := c.Do(req)

	fmt.Println(resp)
}
