package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {

	urlRequested := os.Args[1]

	u, err := url.Parse(urlRequested)

	// uri := u.RequestURI()

	h := sha1.New()
	h.Write([]byte(u.RequestURI()))
	urisha := hex.EncodeToString(h.Sum(nil))

	fmt.Println(urisha)

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/tmp/cache")))
	c := &http.Client{Transport: t}
	response, err := c.Get("file:///" + urisha)

	fmt.Println(err)

	if response.StatusCode >= 399 {

		log.Printf("non trovato in cache")

		response, err = http.Get(urlRequested)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer response.Body.Close()

		savedbody, err := os.Create(urisha)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer savedbody.Close()

		// save response body into a file
		io.Copy(savedbody, response.Body)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	fmt.Printf("Data saved into: %s\n", urisha)
}
