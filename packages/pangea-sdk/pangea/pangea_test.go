package pangea_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/pangea"
	"github.com/stretchr/testify/assert"

	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/internal/pangeatesting"
)

func testClient(t *testing.T, url string) *pangea.Client {
	t.Helper()
	return pangea.NewClient("service", pangeatesting.TestConfig(url))
}

func TestDo_When_Nil_Context_Is_Given_It_Returns_Error(t *testing.T) {
	_, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(nil, req, nil)

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

	req, _ := client.NewRequest("POST", "test", nil)
	_, err := client.Do(context.Background(), req, nil)

	assert.Error(t, err)

	pangeaErr, ok := err.(*pangea.APIError)
	assert.True(t, ok)
	assert.NotNil(t, pangeaErr.ResponseHeader)
	assert.Equal(t, "ValidationError", *pangeaErr.ResponseHeader.Status)
}

func TestDo_When_Server_Returns_500_It_Returns_Error(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status": "InternalError",
			"result": null,
			"summary": "error"
		}`)
	})

	req, _ := client.NewRequest("POST", "test", nil)
	_, err := client.Do(context.Background(), req, nil)
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

	req, _ := client.NewRequest("POST", "test", nil)
	body := &struct {
		Key *string `json:"key"`
	}{}
	resp, err := client.Do(context.Background(), req, body)
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
		Key *string `json:"key"`
	}

	reqBody := reqbody{Key: pangea.String("value")}
	req, _ := client.NewRequest("POST", "test", reqBody)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		body := &reqbody{}
		data, _ := ioutil.ReadAll(r.Body)
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

	resp, err := client.Do(context.Background(), req, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Not initialized struct. Can't unmarshal result from response")
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Status)
	assert.Equal(t, *resp.Status, "Success")
}

func TestDo_When_Client_Can_Not_UnMarshall_Response_It_Returns_UnMarshalError(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	req, _ := client.NewRequest("POST", "test", nil)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `ERROR`)
	})

	_, err := client.Do(context.Background(), req, nil)

	var v *pangea.UnmarshalError
	assert.ErrorAs(t, err, &v)
}

func TestDo_When_Client_Can_Not_UnMarshall_Response_Result_Into_Body_It_Returns_APIError(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	req, _ := client.NewRequest("POST", "test", nil)

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
	_, err := client.Do(context.Background(), req, body)

	var v *pangea.APIError
	assert.ErrorAs(t, err, &v)
}

func TestDo_With_Retries_Success(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()
	cfg := pangeatesting.TestConfig(url)
	cfg.Retry = true
	cfg.RetryConfig = &pangea.RetryConfig{
		RetryMax: 1,
	}

	client := pangea.NewClient("service", cfg)
	req, _ := client.NewRequest("POST", "test", nil)

	handler := func() func(w http.ResponseWriter, r *http.Request) {
		var reqCount int
		return func(w http.ResponseWriter, r *http.Request) {
			if reqCount == 1 {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{
					"request_id": "some-id",
					"request_time": "1970-01-01T00:00:00Z",
					"response_time": "1970-01-01T00:00:10Z",
					"status": "Success",
					"summary": "ok"
				}`)
			} else {
				reqCount++
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
	mux.HandleFunc("/test", handler())

	resp, err := client.Do(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Not initialized struct. Can't unmarshal result from response")
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Status)
	assert.Equal(t, "Success", *resp.Status)
}

func TestDo_With_Retries_Error(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()
	cfg := pangeatesting.TestConfig(url)
	cfg.Retry = true
	cfg.RetryConfig = &pangea.RetryConfig{
		RetryMax: 1,
	}

	client := pangea.NewClient("service", cfg)

	req, _ := client.NewRequest("POST", "test", nil)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.Do(context.Background(), req, nil)

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

	req, _ := client.NewRequest("POST", "test", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected error")
	}

	var v *pangea.AcceptedError
	assert.ErrorAs(t, err, &v)
	assert.NotNil(t, v.ResponseHeader.Status)
	assert.Equal(t, "Accepted", *v.ResponseHeader.Status)
}
