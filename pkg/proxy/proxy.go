package proxy

import (
	"io"
	"net/http"

	"github.com/snasphysicist/ferp/v2/pkg/functional"
	"github.com/snasphysicist/ferp/v2/pkg/log"
	"github.com/snasphysicist/ferp/v2/pkg/url"
)

// Remapper constructs the downstream URL from the incoming path
type Remapper struct {
	url.BaseURL
	Mapper url.PathRewriter
}

// ForwardRequest forwards the incoming request to the configured downstream
// and writes out the received reponse to the outgoing response
func (r Remapper) ForwardRequest(w http.ResponseWriter, req *http.Request) {
	url := url.Rewrite(*req.URL, r.BaseURL, r.Mapper)
	dReq, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		log.Errorf("Failed to construct downstream request: %s", err)
		sendInternalErrorResponse(w)
		return
	}
	transferRequestHeaders(req, dReq)
	res, err := (&http.Client{}).Do(dReq)
	if err != nil {
		log.Errorf("Failed to send downstream request: %s", err)
		sendInternalErrorResponse(w)
		return
	}
	defer func() { _ = res.Body.Close() }()
	transferResponseHeaders(res, w)
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Errorf("Failed to forward response body: %s", err)
		return
	}
	log.Infof("Successfully proxied request from %s to %#v",
		req.URL.String(), dReq.URL.String())
}

// sendInternal sends an error response when something goes wrong in the proxy itself
func sendInternalErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(internalErrorMessage))
	if err != nil {
		log.Errorf("Failed to write error response body: %s", err)
	}
}

// transferRequestHeaders copies all headers from "from" to "to"
// except for content-length, which must be autogenerated
func transferRequestHeaders(from *http.Request, to *http.Request) {
	for k, values := range from.Header {
		if functional.Contains(doNotTransfer(), k) {
			break
		}
		for _, value := range values {
			to.Header.Add(k, value)
		}
	}
}

// transferRequestHeaders copies all headers from "from" to "to"
// except for content-length, which must be autogenerated
func transferResponseHeaders(from *http.Response, to http.ResponseWriter) {
	for k, values := range from.Header {
		if functional.Contains(doNotTransfer(), k) {
			break
		}
		for _, value := range values {
			to.Header().Add(k, value)
		}
	}
}

// internalErrorMessage is a standard non-implementation detail leaking error message
const internalErrorMessage = "500: something went wrong"

// doNotTransfer headers which should not be transferred through the proxy
func doNotTransfer() []string {
	return []string{"Content-Length", "Connection", "Close", "Keep-Alive"}
}
