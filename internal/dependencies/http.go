package dependencies

import (
	"net/http"
)

//counterfeiter:generate . HTTPClient

// HTTPClient is an interface to interact with an HTTP Client
type HTTPClient interface {
	// Get(url string) (*http.Response, error)
	// Head(url string) (*http.Response, error)
	// Post(url string, contentType string, body io.Reader) (*http.Response, error)
	// PostForm(url string, values url.Values) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}
