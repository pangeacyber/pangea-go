//go:build unit

package pangea_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/stretchr/testify/assert"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeatesting"
)

func testClient(t *testing.T, url *url.URL) *pangea.Client {
	t.Helper()
	cfg := pangeatesting.TestConfig(url)
	headers := make(map[string]string, 0)
	headers["Key"] = "Value"
	cfg.AdditionalHeaders = headers
	return pangea.NewClient("service", cfg)
}

func TestConfigDefaults(t *testing.T) {
	config, err := pangea.NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, 2, config.MaxRetries)
}

func TestClientCustomUserAgent(t *testing.T) {
	cfg := pangeatesting.TestConfig(&url.URL{Host: "pangea.cloud"})
	cfg.CustomUserAgent = "Test"
	c := pangea.NewClient("service", cfg)
	assert.NotNil(t, c)
}

func TestClientPendingRequests(t *testing.T) {
	config, err := pangea.NewConfig()
	assert.NoError(t, err)

	client := pangea.NewClient("service", config)

	result := client.GetPendingRequestID()
	assert.Empty(t, result)
}

func TestDo_When_Nil_Context_Is_Given_It_Returns_Error(t *testing.T) {
	_, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	req, _ := client.NewRequest("GET", url, nil)
	_, err := client.Do(nil, req, nil, true)

	assert.Error(t, err)
	assert.Equal(t, "context must be non-nil", err.Error())
}

func TestDo_When_Server_Returns_400_It_Returns_Error(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "ValidationError",
			"result": {"errors": []},
			"summary": "bad request"
		}`)
	})

	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, make(map[string]any))
	_, err := client.Do(context.Background(), req, nil, true)

	assert.Error(t, err)

	pangeaErr, ok := err.(*pangea.APIError)
	assert.True(t, ok)
	assert.NotNil(t, pangeaErr.ResponseHeader)
	assert.Equal(t, "ValidationError", *pangeaErr.ResponseHeader.Status)
	assert.NotEmpty(t, pangeaErr.Error())
	assert.NotEmpty(t, pangeaErr.BaseError.Error())
}

func TestDo_When_Server_Returns_500_It_Returns_Error(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		id, err := gonanoid.New()
		if err != nil {
			t.Fatalf("expected no error got: %v", err)
		}
		w.Header().Set("x-request-id", id)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"request_id": "%s",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "InternalError",
			"result": null,
			"summary": "error"
		}`, id)
	})

	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, nil)
	_, err := client.Do(context.Background(), req, nil, true)
	assert.Error(t, err)
	pangeaErr, ok := err.(*pangea.APIError)
	assert.True(t, ok)
	assert.NotNil(t, pangeaErr.ResponseHeader)
	assert.NotNil(t, pangeaErr.ResponseHeader.Status)
	assert.Equal(t, "InternalError", *pangeaErr.ResponseHeader.Status)
}

func TestDo_When_Server_Returns_200_It_UnMarshals_Result_Into_Struct(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "Success",
			"result": {"key": "value"},
			"summary": "ok"
		}`)
	})

	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, nil)
	body := &struct {
		Key *string `json:"key"`
	}{}
	resp, err := client.Do(context.Background(), req, body, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Status)
	assert.Equal(t, "Success", *resp.Status)

	assert.NotNil(t, body.Key)
	assert.Equal(t, "value", *body.Key)
}

func TestDo_Request_With_Body_Sends_Request_With_Json_Body(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	type reqbody struct {
		pangea.BaseRequest
		Key *string `json:"key"`
	}

	reqBody := reqbody{Key: pangea.String("value")}
	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, &reqBody)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		body := &reqbody{}
		data, _ := io.ReadAll(r.Body)
		json.Unmarshal(data, body)
		if body.Key == nil && *body.Key != "value" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "Success",
			"result": null,
			"summary": "ok"
		}`)
	})

	resp, err := client.Do(context.Background(), req, nil, true)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "not initialized struct. Can't unmarshal result from response")
	assert.Nil(t, resp)
}

func TestDo_When_Client_Can_Not_UnMarshall_Response_Result_Into_Body_It_Returns_APIError(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, nil)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "Success",
			"summary": "ok"
		}`)
	})

	body := &struct {
		Key *string `json:"key"`
	}{}
	_, err := client.Do(context.Background(), req, body, true)

	var v *pangea.APIError
	assert.ErrorAs(t, err, &v)
}

func TestDo_With_Retries_Error(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()
	cfg := pangeatesting.TestConfig(url)
	cfg.MaxRetries = 1

	client := pangea.NewClient("service", cfg)

	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, nil)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.Do(context.Background(), req, nil, true)

	assert.Error(t, err)
}

func TestDo_When_Server_Returns_202_It_Returns_AcceptedError(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "Accepted",
			"result": null,
			"summary": "Accepted"
		}`)
	})

	url, _ = client.GetURL("/test")
	req, _ := client.NewRequest("POST", url, nil)
	_, err := client.Do(context.Background(), req, nil, true)

	if err == nil {
		t.Fatal("Expected error")
	}

	var v *pangea.AcceptedError
	ae := err.(*pangea.AcceptedError)
	assert.ErrorAs(t, err, &v)
	assert.NotNil(t, v.ResponseHeader.Status)
	assert.Equal(t, "Accepted", *v.ResponseHeader.Status)
	assert.NotEmpty(t, err.Error())
	assert.True(t, true, ae.Is(err))
	assert.NotEmpty(t, ae.ReqID())
}
