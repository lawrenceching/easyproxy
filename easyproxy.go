package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type proxy struct {
	destination string
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	DEBUG := os.Getenv("EASYPROXY_DEBUG")

	if DEBUG == "true" {
		log.Println("[" + req.Method + "] " + req.RemoteAddr + req.RequestURI + " -> " + p.destination + req.RequestURI)
	}

	client := &http.Client{}

	method := req.Method
	body := req.Body
	reqToProxy, err := http.NewRequest(method, p.destination+req.RequestURI, body)

	delHopHeaders(req.Header)

	for k, v := range req.Header {
		for _, h := range v {
			reqToProxy.Header.Add(k, h)
		}
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		appendHostToXForwardHeader(req.Header, clientIP)
	}

	resp, err := client.Do(reqToProxy)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	delHopHeaders(resp.Header)

	copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)

	buf := make([]byte, 1024)
	len := 1
	for len != 0 {
		l, _ := resp.Body.Read(buf)
		len = l
		wr.Write(buf[:len])
	}
}

func CreateProxy(from string, to string) http.Server {
	handler := &proxy{destination: to}
	server := http.Server{Addr: from, Handler: handler}
	return server
}

func StartProxy(from string, to string) {
	server := CreateProxy(from, to)
	server.ListenAndServe()
}

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("easyproxy\n" +
			"    --from <host-address> \n" +
			"    --to <target-address>")
		return
	}

	var from = flag.String("from", "127.0.0.1:8080", "The address listen to")
	var to = flag.String("to", "example.com", "The target proxy to")
	flag.Parse()

	log.Println("easyproxy is proxying request from", *from, "to", *to)

	StartProxy(*from, *to)
}
