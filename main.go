package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

// set to false, if you don't want the body to be cached
const dumpBody = true

var filesToSave = regexp.MustCompile(`.*(ico|gif|pdf|jpg|css|js|mp4f)$`)

func main() {

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(`.*gif$`))).HandleConnect(goproxy.AlwaysMitm).

	proxy.OnRequest(goproxy.UrlMatches(filesToSave)).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// r.Header.Set("X-GoProxy", "yxorPoG-X")

			return r, retrieveContent(r)
		})

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
