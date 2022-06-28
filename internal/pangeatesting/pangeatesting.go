package pangeatesting

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/pangeacyber/go-pangea/pangea"
)

const baseURLPath = "/api"

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

func TestConfig(url string) *pangea.Config {
	return &pangea.Config{
		Token:    "TestToken",
		Endpoint: url,
	}
}

func TestMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := strings.Trim(string(b), "\n"); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}
