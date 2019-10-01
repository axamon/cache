package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgraph-io/ristretto"
)

var cache *ristretto.Cache

var essencePath = "./essenze/"

func init() {
	initcache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	// Creates the path to store essences.
	err = os.MkdirAll(essencePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	cache = initcache
}

type Essence struct {
	URI           string
	LastModified  time.Time
	MaxAge        int
	ContentLength int64
	CachedBody    string
}

func retrieveContent(r *http.Request) *http.Response {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urlRequested := r.RequestURI

	var err error

	u, err := url.Parse(urlRequested)
	if err != nil {
		log.Fatal(err)
	}

	//uri := u.RequestURI()

	h := sha1.New()
	h.Write([]byte(u.RequestURI()))
	uriHash := hex.EncodeToString(h.Sum(nil))

	fmt.Println(uriHash)

	var response *http.Response

	// Look for key in cache
	essenceValues, found := cache.Get(uriHash)
	if found {

		if essenceValues.(Essence).LastModified.Before(time.Now().Add(-60 * time.Minute)) {
			fmt.Println("Essenza vecchia")
			goto RETRIEVE
		}
		response, err = clientCache(ctx, r, uriHash)
		if err != nil {
			log.Fatal(err)
		}
		return response
	}
RETRIEVE:
	//response = retrieveAndSave(ctx, r, uriHash)
	response, err = http.Get(r.RequestURI)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// chiude il body
	response.Body.Close()

	savedbody, err := os.Create(essencePath + uriHash)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//defer savedbody.Close()
	//defer response.Body.Close()

	_, err = savedbody.Write(bodyBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	savedbody.Close()

	// save response body into a file
	// _, err = io.Copy(savedbody, response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var essence = Essence{URI: r.RequestURI, ContentLength: r.ContentLength, LastModified: time.Now()}
	cache.Set(uriHash, essence, 1)

	// response, err = clientCache(ctx, r, uriHash)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(cache.Get(uriHash))
	return response
}

func retrieveAndSave(ctx context.Context, r *http.Request, uriHash string) *http.Response {

	var err error

	response, err := http.Get(r.RequestURI)

	if err != nil {
		fmt.Println(err)
	}

	savedbody, err := os.Create(essencePath + uriHash)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer savedbody.Close()

	// save response body into a file
	_, err = io.Copy(savedbody, response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var essence = Essence{URI: r.RequestURI, ContentLength: r.ContentLength, LastModified: time.Now()}
	cache.Set(uriHash, essence, 1)

	return response

}

func clientCache(ctx context.Context, r *http.Request, uriHash string) (*http.Response, error) {

	fmt.Println("Trovato")

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(essencePath)))
	c := &http.Client{Transport: t}
	response, err := c.Get("file:///" + uriHash)

	return response, err
}
