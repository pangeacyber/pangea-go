package pangeatesting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/audit"
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

	url, _ := url.Parse(server.URL + baseURLPath)

	return mux, url.String(), server.Close
}

func TestConfig(url string) *pangea.Config {
	// Clean scheme. It will be adden after decide if it should be secure o insecure
	// It only happens on testing because of local server
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	}

	return &pangea.Config{
		Token:      "TestToken",
		Domain:     url,
		Insecure:   true,
		Enviroment: "local",
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
		t.Fatalf("Error reading request body: %v", err)
	}
	if got := strings.Trim(string(b), "\n"); got != want {
		t.Fatalf("request Body is %s, want %s", got, want)
	}
}

func TestNewRequestAndDoFailure(t *testing.T, method string, f func(cfg *pangea.Config) error) {
	t.Helper()

	emptyDomainCfg := &pangea.Config{Domain: ""}
	doErr := f(emptyDomainCfg)
	if doErr == nil {
		t.Fatalf("call to method %v with empty Enpoint got nil err, want error", method)
	}

	badUrlCfg := &pangea.Config{Domain: "htt://   "}
	newRequestErr := f(badUrlCfg)
	if newRequestErr == nil {
		t.Fatalf("call to method %v with bad Domain got nil err, want error", method)
	}
}

func CreateFile(t *testing.T, contents []byte) *os.File {
	t.Helper()
	tmpdir := t.TempDir()
	file, err := ioutil.TempFile(tmpdir, "*")
	if err != nil {
		t.Fatal("failed to creat temp file")
	}
	file.Write(contents)
	return file
}

func GetEnvVarOrSkip(t *testing.T, varname string) string {
	t.Helper()
	envVar := os.Getenv(varname)
	if envVar == "" {
		t.Skipf("set %v env variable to run this test", varname)
	}
	return envVar
}

type TestEnvironment string

const (
	Live    TestEnvironment = "LVE"
	Develop                 = "DEV"
	Staging                 = "STG"
)

func GetTestDomain(t *testing.T, env TestEnvironment) string {
	t.Helper()
	varname := "PANGEA_INTEGRATION_DOMAIN_" + string(env)
	return GetEnvVarOrSkip(t, varname)
}

func GetTestToken(t *testing.T, env TestEnvironment) string {
	t.Helper()
	varname := "PANGEA_INTEGRATION_TOKEN_" + string(env)
	return GetEnvVarOrSkip(t, varname)
}

func GetVaultSignatureTestToken(t *testing.T, env TestEnvironment) string {
	t.Helper()
	varname := "PANGEA_INTEGRATION_VAULT_TOKEN_" + string(env)
	return GetEnvVarOrSkip(t, varname)
}

func GetCustomSchemaTestToken(t *testing.T, env TestEnvironment) string {
	t.Helper()
	varname := "PANGEA_INTEGRATION_CUSTOM_SCHEMA_TOKEN_" + string(env)
	return GetEnvVarOrSkip(t, varname)
}

func GetMultiConfigTestToken(t *testing.T, env TestEnvironment) string {
	t.Helper()
	varname := "PANGEA_INTEGRATION_MULTI_CONFIG_TOKEN_" + string(env)
	return GetEnvVarOrSkip(t, varname)
}

func GetConfigID(t *testing.T, env TestEnvironment, service string, configNumber int) string {
	t.Helper()
	varname := fmt.Sprintf("PANGEA_%s_CONFIG_ID_%d_%s", strings.ToUpper(service), configNumber, string(env))
	return GetEnvVarOrSkip(t, varname)
}

type CustomSchemaEvent struct {
	Message       string              `json:"message"`
	FieldInt      int                 `json:"field_int,omitempty"`
	FieldBool     bool                `json:"field_bool,omitempty"`
	FieldStrShort string              `json:"field_str_short,omitempty"`
	FieldStrLong  string              `json:"field_str_long,omitempty"`
	FieldTime     *pu.PangeaTimestamp `json:"field_time,omitempty"`

	// TenantID field
	TenantID string `json:"tenant_id,omitempty"`
}

func (_ *CustomSchemaEvent) NewFromJSON(b []byte) (audit.IEvent, error) {
	var e CustomSchemaEvent

	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (e *CustomSchemaEvent) GetTenantID() string {
	return e.TenantID
}

func (e *CustomSchemaEvent) SetTenantID(tid string) {
	e.TenantID = tid
}
