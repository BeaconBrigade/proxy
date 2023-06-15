package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var proxyAddr string = "https://stackoverflow.com"

const localAddr string = "localhost:3000"

func main() {
	args := os.Args
	if len(args) > 1 {
		proxyAddr = args[1]
	}
	log.Printf("proxying to: %s", proxyAddr)
	http.HandleFunc("/", root)
	http.ListenAndServe(localAddr, nil)
}

func root(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// send request body through proxy
	bod, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("couldn't read request body: %v", err)
	}
	new, err := http.NewRequest(req.Method, fmt.Sprintf("%s%s", proxyAddr, path), bytes.NewReader(bod))
	if err != nil {
		log.Fatalf("could not build req: %v", err)
	}

	// send headers through proxy
	for name, headers := range req.Header {
		for _, val := range headers {
			new.Header.Add(name, val)
		}
	}

	defer req.Body.Close()

	// send request
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Printf("redirected to: %s", req.URL)
			return nil
		},
	}
	res, err := client.Do(new)
	if err != nil {
		log.Printf("could not send request: %v", err)
	}

	defer res.Body.Close()

	// copy headers to return
	for name, headers := range res.Header {
		for _, val := range headers {
			w.Header().Add(name, val)
		}
	}

	// copy body to return
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("couldn't read body: %v", err)
		return
	}

	// fix links
	localWithoutHttp := regexp.MustCompile(`https?:\/\/`).ReplaceAll([]byte(localAddr), nil)
	regex := regexp.MustCompile(fmt.Sprintf(`^(?:http(?:s?):\/\/?)?([A-Za-z0-9_:.-]\.?%s+)\/?`, localWithoutHttp))
	if strings.HasPrefix(res.Header.Get("Content-Type"), "text/") {
		b = regex.ReplaceAll(b, []byte(proxyAddr))
	}

	w.Write(b)
}
