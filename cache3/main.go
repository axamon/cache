package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elazarl/goproxy"

	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/redis"
)

// set to false, if you don't want the body to be cached
const dumpBody = true

func main() {


    ringOpt := &redis.RingOptions{
        Addrs: map[string]string{
            "server": ":6379",
        },
    }
    cacheClient , err := cache.NewClient(
        cache.ClientWithAdapter(redis.NewAdapter(ringOpt)),
        cache.ClientWithTTL(10 * time.Minute),
        cache.ClientWithRefreshKey("opn"),
    )
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(`.*gif$`))).HandleConnect(goproxy.AlwaysMitm).

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// r.Header.Set("X-GoProxy", "yxorPoG-X")
			return r, nil
		})

	log.Fatal(http.ListenAndServe(":8080", cacheClient.Middleware(proxy)))
}
