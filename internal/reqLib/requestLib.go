package reqLib

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

// RequestLib exposes a shared HTTP client configured for the local Riot endpoints.
type RequestLib struct {
	mu     sync.RWMutex
	client *http.Client
}

var (
	requestLibOnce      sync.Once
	requestLibSingleton *RequestLib
)

// Instance returns the singleton RequestLib instance, constructing it on first use.
func Instance() *RequestLib {
	requestLibOnce.Do(func() {
		requestLibSingleton = &RequestLib{
			client: defaultHTTPClient(),
		}
	})

	return requestLibSingleton
}

// SetClient swaps the underlying HTTP client used by the singleton.
func (r *RequestLib) SetClient(client *http.Client) {
	if client == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.client = client
}

// Client returns the currently configured HTTP client.
func (r *RequestLib) Client() *http.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.client
}

// Do executes the provided request using the shared HTTP client.
func (r *RequestLib) Do(req *http.Request) (*http.Response, error) {
	return r.Client().Do(req)
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
