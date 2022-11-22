package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Remapper struct {
	Protocol string
	Host     string
	Port     uint16
	Base     string
	Mapper   pathRewriter
}

type pathRewriter func(string) string

func (r Remapper) ForwardRequest(w http.ResponseWriter, req *http.Request) {
	mapped := r.Mapper(req.URL.Path)
	url := r.makeURL(mapped)
	dReq, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		log.Printf("Failed to construct downstream request: %s", err)
		sendHTTPErrorResponse(w, 500, internalErrorMessage)
		return
	}
	transferRequestHeaders(req, dReq)
	res, err := (&http.Client{}).Do(dReq)
	if err != nil {
		log.Printf("Failed to send downstream request: %s", err)
		sendHTTPErrorResponse(w, 500, internalErrorMessage)
		return
	}
	transferResponseHeaders(res, w)
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Printf("Failed to forward response body: %s", err)
		return
	}
	log.Printf("Successfully proxied request from %s to %#v",
		req.URL.String(), dReq.URL.String())
}

func (r Remapper) makeURL(suffix string) string {
	base := strings.TrimLeft(r.Base, "/")
	if suffix == "" {
		url := fmt.Sprintf("%s://%s:%d/%s", r.Protocol, r.Host, r.Port, base)
		log.Printf("No suffix, returning %s (base '%s')", url, base)
		return url
	}
	suffix = strings.TrimLeft(suffix, "/")
	base = strings.TrimRight(base, "/")
	if base != "" {
		base = fmt.Sprintf("%s/", base)
	}
	url := fmt.Sprintf("%s://%s:%d/%s%s", r.Protocol, r.Host, r.Port, base, suffix)
	log.Printf("Had suffix, returning %s (base '%s', suffix '%s')", url, base, suffix)
	return url
}

func sendHTTPErrorResponse(w http.ResponseWriter, code int, message string) {
	b, err := json.Marshal(struct {
		E string `json:"error"`
	}{E: message})
	if err != nil {
		code = 500
		b = []byte("something went wrong")
	}
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("Failed to write error response body: %s", err)
	}
}

func transferRequestHeaders(from *http.Request, to *http.Request) {
	for k, values := range from.Header {
		if k == "content-length" {
			break
		}
		for _, value := range values {
			to.Header.Add(k, value)
		}
	}
}

func transferResponseHeaders(from *http.Response, to http.ResponseWriter) {
	for k, values := range from.Header {
		if k == "content-length" {
			break
		}
		for _, value := range values {
			to.Header().Add(k, value)
		}
	}
}

// internalErrorMessage is a standard non-implementation detail leaking error message
const internalErrorMessage = "Something went wrong"
