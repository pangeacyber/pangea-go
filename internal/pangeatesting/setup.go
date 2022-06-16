package pangeatesting

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

const baseURLPath = "/api/v1"

// SetupServer sets up a test HTTP server
//
// Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func SetupServer() (mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()

	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	url, _ := url.Parse(server.URL + baseURLPath + "/")

	return mux, url.String(), server.Close
}
