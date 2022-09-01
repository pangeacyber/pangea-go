package audit_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/audit"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/log", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"event":{"message":"test"},"return_hash":true,"verbose":true}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result": {
					"canonical_event_base64": "eyJtZXNzYWdlIjoicHJ1ZWJhXzQ1NiIsInJlY2VpdmVkX2F0IjoiMjAyMi0wNi0yOFQfadDowMjowNS40ODAyNjdaIn0=",
					"event": {
						"message": "test"
					},
					"hash": "b0e7b01c733ed4983e4c706206a8e6a77a00503ffadb13a3ab27f37ae1dd8484"
				},
				"summary": "Logged 1 record(s)"
			}`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url))
	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String("test"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Log(ctx, input)

	assert.NoError(t, err)

	want := &audit.LogOutput{
		CanonicalEventBase64: pangea.String("eyJtZXNzYWdlIjoicHJ1ZWJhXzQ1NiIsInJlY2VpdmVkX2F0IjoiMjAyMi0wNi0yOFQfadDowMjowNS40ODAyNjdaIn0="),
		Hash:                 pangea.String("b0e7b01c733ed4983e4c706206a8e6a77a00503ffadb13a3ab27f37ae1dd8484"),
		Event: &audit.LogEventOutput{
			Event: audit.Event{
				Message: pangea.String("test"),
			},
		},
	}
	assert.Equal(t, want, got.Result)
}

func TestDomainTrailingSlash(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/log", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"event":{"message":"test"},"return_hash":true,"verbose":true}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result": {
					"canonical_event_base64": "eyJtZXNzYWdlIjoicHJ1ZWJhXzQ1NiIsInJlY2VpdmVkX2F0IjoiMjAyMi0wNi0yOFQfadDowMjowNS40ODAyNjdaIn0=",
					"event": {
						"message": "test"
					},
					"hash": "b0e7b01c733ed4983e4c706206a8e6a77a00503ffadb13a3ab27f37ae1dd8484"
				},
				"summary": "Logged 1 record(s)"
			}`)
	})

	url = url + "/" // Add trailing slash to domain

	client, _ := audit.New(pangeatesting.TestConfig(url))
	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String("test"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Log(ctx, input)

	assert.NoError(t, err)

	want := &audit.LogOutput{
		CanonicalEventBase64: pangea.String("eyJtZXNzYWdlIjoicHJ1ZWJhXzQ1NiIsInJlY2VpdmVkX2F0IjoiMjAyMi0wNi0yOFQfadDowMjowNS40ODAyNjdaIn0="),
		Hash:                 pangea.String("b0e7b01c733ed4983e4c706206a8e6a77a00503ffadb13a3ab27f37ae1dd8484"),
		Event: &audit.LogEventOutput{
			Event: audit.Event{
				Message: pangea.String("test"),
			},
		},
	}
	assert.Equal(t, want, got.Result)
}

func TestSearch(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	t1 := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60))
	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","include_membership_proof":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result": {
					"count": 2,
					"events": [
						{
							"event": {
								"message": "test_2",
								"received_at": "%[1]v"
							},
							"leaf_index": 2,
							"membership_proof": "some-proof"
						},
						{
							"event": {
								"message": "test_1",
								"received_at": "%[1]v"
							},
							"leaf_index": 3,
							"membership_proof": "some-proof"
						}
					],
					"expires_at": "%[1]v",
					"id": "some-id"
				},
				"summary": "Found 13 event(s)"
			}`, t1.Format(time.RFC3339))
	})

	client, _ := audit.New(pangeatesting.TestConfig(url))
	input := &audit.SearchInput{
		Query:                  pangea.String("message:test"),
		IncludeMembershipProof: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)

	want := &audit.SearchOutput{
		Count:     pangea.Int(2),
		ExpiresAt: &t1,
		ID:        pangea.String("some-id"),
		Events: audit.Events{
			{
				Event: &audit.Event{
					Message: pangea.String("test_2"),
				},
				LeafIndex:       pangea.Int(2),
				MembershipProof: pangea.String("some-proof"),
			},
			{
				Event: &audit.Event{
					Message: pangea.String("test_1"),
				},
				LeafIndex:       pangea.Int(3),
				MembershipProof: pangea.String("some-proof"),
			},
		},
	}
	assert.Equal(t, want, got.Result)
}

func TestSearchResults(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	t1 := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60))
	mux.HandleFunc("/v1/results", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"id":"some-id","include_membership_proof":true,"limit":50}`)
		fmt.Fprintf(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result": {
					"count": 2,
					"events": [
						{
							"event": {
								"message": "test_2",
								"received_at": "%[1]v"
							},
							"leaf_index": 2,
							"membership_proof": "some-proof"
						},
						{
							"event": {
								"message": "test_1",
								"received_at": "%[1]v"
							},
							"leaf_index": 3,
							"membership_proof": "some-proof"
						}
					],
					"root": {
						"published_at": "%[1]v",
						"root_hash": "3a2563b40abe941f21c2ea929f2be92606fd2545762d3fa755ecec2942f5d452",
						"size": 11,
						"consistency_proof": [
							"x:68810d719dc9dccee268d17a6c5baf3bf12d7ffad5673b763e06338121ed4e46,r:4a291c09b0bed8303d3e7f91315bd47da3df422151e642ded4208f46342a12f6",
							"x:82eba5aa211af097d22ecf215be386212c192d7068a02aeb4280905e4d200ff9,r:03c513b31ev80f4c871dbcd07b069fd369482529984f0770008f6c7777813026"
						]
					}
				},
				"summary": "Found 13 event(s)"
			}`, t1.Format(time.RFC3339))
	})

	client, _ := audit.New(pangeatesting.TestConfig(url))
	input := &audit.SearchResultInput{
		ID:                     pangea.String("some-id"),
		IncludeMembershipProof: pangea.Bool(true),
		Limit:                  pangea.Int(50),
	}
	ctx := context.Background()
	got, err := client.SearchResults(ctx, input)

	assert.NoError(t, err)

	want := &audit.SearchResultOutput{
		Count: pangea.Int(2),
		Events: audit.Events{
			{
				Event: &audit.Event{
					Message: pangea.String("test_2"),
				},
				LeafIndex:       pangea.Int(2),
				MembershipProof: pangea.String("some-proof"),
			},
			{
				Event: &audit.Event{
					Message: pangea.String("test_1"),
				},
				LeafIndex:       pangea.Int(3),
				MembershipProof: pangea.String("some-proof"),
			},
		},
		Root: &audit.Root{
			PublishedAt: &t1,
			RootHash:    pangea.String("3a2563b40abe941f21c2ea929f2be92606fd2545762d3fa755ecec2942f5d452"),
			Size:        pangea.Int(11),
			ConsistencyProof: []*string{
				pangea.String("x:68810d719dc9dccee268d17a6c5baf3bf12d7ffad5673b763e06338121ed4e46,r:4a291c09b0bed8303d3e7f91315bd47da3df422151e642ded4208f46342a12f6"),
				pangea.String("x:82eba5aa211af097d22ecf215be386212c192d7068a02aeb4280905e4d200ff9,r:03c513b31ev80f4c871dbcd07b069fd369482529984f0770008f6c7777813026"),
			},
		},
	}
	assert.Equal(t, want, got.Result)
}

func TestRoot(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	t1 := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60))
	mux.HandleFunc("/v1/root", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"tree_size":11}`)
		fmt.Fprintf(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result": {
					"data":  {
						"published_at": "%v",
						"root_hash": "3a2563b40abe941f21c2ea929f2be92606fd2545762d3fa755ecec2942f5d452",
						"size": 11,
						"consistency_proof": [
							"x:68810d719dc9dccee268d17a6c5baf3bf12d7ffad5673b763e06338121ed4e46,r:4a291c09b0bed8303d3e7f91315bd47da3df422151e642ded4208f46342a12f6",
							"x:82eba5aa211af097d22ecf215be386212c192d7068a02aeb4280905e4d200ff9,r:03c513b31ev80f4c871dbcd07b069fd369482529984f0770008f6c7777813026"
						]
					}
				},
				"summary": "success"
			}`, t1.Format(time.RFC3339))
	})

	client, _ := audit.New(pangeatesting.TestConfig(url))
	input := &audit.RootInput{
		TreeSize: pangea.Int(11),
	}
	ctx := context.Background()
	got, err := client.Root(ctx, input)

	assert.NoError(t, err)

	want := &audit.RootOutput{
		Data: &audit.Root{
			PublishedAt: &t1,
			RootHash:    pangea.String("3a2563b40abe941f21c2ea929f2be92606fd2545762d3fa755ecec2942f5d452"),
			Size:        pangea.Int(11),
			ConsistencyProof: []*string{
				pangea.String("x:68810d719dc9dccee268d17a6c5baf3bf12d7ffad5673b763e06338121ed4e46,r:4a291c09b0bed8303d3e7f91315bd47da3df422151e642ded4208f46342a12f6"),
				pangea.String("x:82eba5aa211af097d22ecf215be386212c192d7068a02aeb4280905e4d200ff9,r:03c513b31ev80f4c871dbcd07b069fd369482529984f0770008f6c7777813026"),
			},
		},
	}

	assert.Equal(t, want, got.Result)
}

func TestLogError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.Log(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.Log", f)
}

func TestSearchError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.Search(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.Search", f)
}

func TestSearchResultsError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.SearchResults(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.SearchResults", f)
}

func TestRootError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.Root(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.Root", f)
}

func TestFailedOptions(t *testing.T) {
	_, err := audit.New(
		pangeatesting.TestConfig("url"),
		audit.WithLogSigningEnabled("bad file name"),
	)
	assert.Error(t, err)

	_, err = audit.New(
		pangeatesting.TestConfig("url"),
		audit.WithLogSignatureVerificationEnabled(),
	)
	assert.NoError(t, err)
}
