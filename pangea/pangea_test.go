package pangea_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pangeacyber/go-pangea/pangea"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
)

func testClient(t *testing.T, url string) *pangea.Client {
	t.Helper()
	return pangea.NewClient("service", &pangea.Config{Token: "TestToken", Endpoint: url})
}

func TestDo_When_Nil_Context_Is_Given_It_Returns_Error(t *testing.T) {
	_, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(nil, req, nil)

	if err == nil {
		t.Errorf("Expected error")
	}

	if err.Error() != "context must be non-nil" {
		t.Errorf("Expected error message to be 'context must be non-nil', got %s", err.Error())
	}
}

func TestDo_When_Server_Returns_400_It_Returns_Error(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	client := testClient(t, url)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"request_id": "some-id",
			"request_time": "1970-01-01T00:00:00Z",
			"response_time": "1970-01-01T00:00:10Z",
			"status_code": 400,
			"status": "error",
			"result": null,
			"summary": "bad request"
		}`)
	})

	req, _ := client.NewRequest("POST", "test", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}

	pangeaErr, ok := err.(*pangea.APIError)
	if !ok {
		t.Fatalf("Expected pangea.ErrorResponse, got %T", err)
	}
	if pangeaErr.ResponseMetadata == nil {
		t.Fatal("Expected ResponseMetadata to be non-nil")
	}
	if pangeaErr.ResponseMetadata.StatusCode == nil {
		t.Fatal("Expected non-nil status code")
	}
	if *pangeaErr.ResponseMetadata.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, *pangeaErr.ResponseMetadata.StatusCode)
	}
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
			"status_code": 500,
			"status": "error",
			"result": null,
			"summary": "error"
		}`)
	})

	req, _ := client.NewRequest("POST", "test", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 500 error, got no error.")
	}

	pangeaErr, ok := err.(*pangea.APIError)
	if !ok {
		t.Fatalf("Expected pangea.ErrorResponse, got %v", err)
	}
	if pangeaErr.ResponseMetadata == nil {
		t.Fatal("Expected ResponseMetadata to be non-nil")
	}
	if pangeaErr.ResponseMetadata.StatusCode == nil {
		t.Fatal("Expected non-nil status code")
	}
	if *pangeaErr.ResponseMetadata.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, *pangeaErr.ResponseMetadata.StatusCode)
	}
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
			"status_code": 200,
			"status": "ok",
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

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.StatusCode == nil || *resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status code 200")
	}

	if body.Key == nil {
		t.Fatal("Expected body.Key to be non-nil")
	}

	if *body.Key != "value" {
		t.Errorf("Expected body.Key to be 'value', got %v", *body.Key)
	}
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
			"status_code": 200,
			"status": "ok",
			"result": null,
			"summary": "ok"
		}`)
	})

	resp, err := client.Do(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.StatusCode == nil || *resp.StatusCode != http.StatusOK {
		t.Error("Expected status code 200")
	}
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

	if err == nil {
		t.Fatal("Expected an error")
	}

	_, ok := err.(*pangea.UnMarshalError)
	if !ok {
		t.Errorf("Expected pangea.UnMarshalError, got %T", err)
	}
}

func TestDo_When_Client_Can_Not_UnMarshall_Response_Result_Into_Body_It_Returns_UnMarshalError(t *testing.T) {
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
			"status_code": 200,
			"status": "ok",
			"summary": "ok"
		}`)
	})

	body := &struct {
		Key *string `json:"key"`
	}{}
	_, err := client.Do(context.Background(), req, body)

	if err == nil {
		t.Fatal("Expected an error")
	}

	_, ok := err.(*pangea.UnMarshalError)
	if !ok {
		t.Errorf("Expected pangea.UnMarshalError, got %T", err)
	}
}
