package audit_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/audit"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/log", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"event":{"message":"test"},"verbose":true}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result": {
					"envelope": {
						"event": {
							"message": "test"
						}
					},
					"hash": "9c9c3b5a627cce035d517c14c10779656e900532bf6e76a5d2c69148e45fdb8d"
				},
				"summary": "Logged 1 record(s)"
			}`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url))
	event := audit.Event{
		Message: "test",
	}
	ctx := context.Background()
	got, err := client.Log(ctx, event, true)

	assert.NoError(t, err)

	want := &audit.LogOutput{
		Hash: "9c9c3b5a627cce035d517c14c10779656e900532bf6e76a5d2c69148e45fdb8d",
		EventEnvelope: &audit.EventEnvelope{
			Event: &audit.Event{
				Message: "test",
			},
		},
		RawEnvelope: got.Result.RawEnvelope,
	}
	assert.Equal(t, want, got.Result)
}

func TestLog_FailHashVerification(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/log", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"event":{"message":"test"},"verbose":true}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result": {
					"envelope": {
						"event": {
							"message": "test"
						}
					},
					"hash": "notarealhash"
				},
				"summary": "Logged 1 record(s)"
			}`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url))
	event := audit.Event{
		Message: "test",
	}
	ctx := context.Background()
	got, err := client.Log(ctx, event, true)

	assert.Error(t, err)
	assert.NotNil(t, err)
	assert.Nil(t, got)
}

func TestLog_FailSigner(t *testing.T) {
	client, err := audit.New(pangeatesting.TestConfig("someurl"), audit.WithLogLocalSigning("notarealkey"))

	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, err.Error(), "audit: failed signer creation: signer: cannot read file notarealkey: open notarealkey: no such file or directory")
}

func TestDomainTrailingSlash(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/log", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"event":{"message":"test"},"verbose":true}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result": {
					"envelope": {
						"event": {
							"message": "test"
						}
					},
					"hash": "9c9c3b5a627cce035d517c14c10779656e900532bf6e76a5d2c69148e45fdb8d"
				},
				"summary": "Logged 1 record(s)"
			}`)
	})

	url = url + "/" // Add trailing slash to domain

	client, _ := audit.New(pangeatesting.TestConfig(url))
	event := audit.Event{
		Message: "test",
	}
	ctx := context.Background()
	got, err := client.Log(ctx, event, true)

	assert.NoError(t, err)

	want := &audit.LogOutput{
		Hash: "9c9c3b5a627cce035d517c14c10779656e900532bf6e76a5d2c69148e45fdb8d",
		EventEnvelope: &audit.EventEnvelope{
			Event: &audit.Event{
				Message: "test",
			},
		},
		RawEnvelope: got.Result.RawEnvelope,
	}
	assert.Equal(t, want, got.Result)
}

func TestSearch(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	t1 := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60))
	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result": {
					"count": 2,
					"events": [
						{
							"envelope": {
								"event": {
									"message": "test_2"
								}
							},
							"received_at": "%[1]v",
							"leaf_index": 2,
							"membership_proof": "some-proof"
						},
						{
							"envelope": {
								"event": {
									"message": "test_1"
								}
							},
							"received_at": "%[1]v",
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
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)

	want := &audit.SearchOutput{
		Count:     2,
		ExpiresAt: &t1,
		ID:        "some-id",
		Events: audit.SearchEvents{
			{
				EventEnvelope: &audit.EventEnvelope{
					Event: &audit.Event{
						Message: "test_2",
					},
					ReceivedAt: got.Result.Events[1].EventEnvelope.ReceivedAt,
				},
				LeafIndex:       pangea.Int(2),
				MembershipProof: pangea.String("some-proof"),
				RawEnvelope:     got.Result.Events[0].RawEnvelope,
			},
			{
				EventEnvelope: &audit.EventEnvelope{
					Event: &audit.Event{
						Message: "test_1",
					},
					ReceivedAt: got.Result.Events[1].EventEnvelope.ReceivedAt,
				},
				LeafIndex:       pangea.Int(3),
				MembershipProof: pangea.String("some-proof"),
				RawEnvelope:     got.Result.Events[1].RawEnvelope,
			},
		},
	}
	assert.Equal(t, want, got.Result)
}

func TestSearch_Verify(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "published": true,
					  "hash": "afa77464cad6e1b34e23d4847106081577f0b78f9c407ab501d16c09b23be202",
					  "leaf_index": 30,
					  "membership_proof": "l:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, audit.Success, got.Result.Events[0].MembershipVerification)
}

func TestSearch_InvalidEventHash(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "published": true,
					  "hash": "notarealhash",
					  "leaf_index": 30,
					  "membership_proof": "l:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled(), audit.DisableEventVerification())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, audit.Failed, got.Result.Events[0].MembershipVerification)
}

func TestSearch_InvalidSideInMembershipProof(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "published": true,
					  "hash": "afa77464cad6e1b34e23d4847106081577f0b78f9c407ab501d16c09b23be202",
					  "leaf_index": 30,
					  "membership_proof": "notvalidside:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, audit.Failed, got.Result.Events[0].MembershipVerification)
}

func TestSearch_InvalidHashInMembershipProof(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "published": true,
					  "hash": "afa77464cad6e1b34e23d4847106081577f0b78f9c407ab501d16c09b23be202",
					  "leaf_index": 30,
					  "membership_proof": "l:notavalidhash,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, len(got.Result.Events))
	assert.Equal(t, audit.Failed, got.Result.Events[0].MembershipVerification)
}

func TestSearch_InvalidRootHash(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "published": true,
					  "hash": "afa77464cad6e1b34e23d4847106081577f0b78f9c407ab501d16c09b23be202",
					  "leaf_index": 30,
					  "membership_proof": "l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "notavalidhash",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, len(got.Result.Events))
	assert.Equal(t, audit.Failed, got.Result.Events[0].MembershipVerification)
}

// There is a , after "action" in event
func TestSearch_FailedToUnmarshall(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "published": true,
					  "hash": "afa77464cad6e1b34e23d4847106081577f0b78f9c407ab501d16c09b23be202",
					  "leaf_index": 30,
					}
				  ],
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, got)

}

func TestSearch_VerifyFailSignature(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 1,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "notarealsignature"
					  },
					  "hash": "9ddc99bf74c65b345c442604f3ce84288218c4548499a761018bf13473d252d0",
					  "leaf_index": 30,
					  "membership_proof": "l:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, got.Result.Count)
	assert.Equal(t, 1, len(got.Result.Events))
	assert.Equal(t, audit.Failed, got.Result.Events[0].SignatureVerification)
}

func TestSearch_VerifyFailPublicKey(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 1,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "notarealpublickey",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "hash": "8c779d38f7ec4a88a4e09d064b6868ccc16da8696692de688d9367c28a3bdb08",
					  "leaf_index": 30,
					  "membership_proof": "l:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, got.Result.Count)
	assert.Equal(t, 1, len(got.Result.Events))
	assert.Equal(t, audit.Failed, got.Result.Events[0].SignatureVerification)
}

// deleted proof members
func TestSearch_VerifyFailMembershipProof(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 1,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "hash": "afa77464cad6e1b34e23d4847106081577f0b78f9c407ab501d16c09b23be202",
					  "leaf_index": 30,
					  "membership_proof": "l:notarealmembershipproof",
					  "published": true
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, got.Result.Count)
	assert.Equal(t, 1, len(got.Result.Events))
	assert.Equal(t, audit.Failed, got.Result.Events[0].MembershipVerification)
}

func TestSearch_VerifyFailHash(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 4,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "hash": "notarealhash",
					  "leaf_index": 30,
					  "membership_proof": "l:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "audit: cannot verify hash of record. Hash: [notarealhash]")
	assert.Nil(t, got)
}

func TestSearch_VerifyFailHashEmpty(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/search", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"query":"message:test","verbose":true}`)
		fmt.Fprintf(w,
			`{
				"request_id": "prq_pnlmbzvj4ytk7juvhlkwp5x4djeyiwov",
				"request_time": "2022-09-20T15:15:48.743Z",
				"response_time": "2022-09-20T15:15:49.772Z",
				"status": "Success",
				"summary": "Found 1 event(s)",
				"result": {
				  "id": "pit_q2zjhuymmbclgzsfg2dwi5bslswxbxd5",
				  "count": 1,
				  "expires_at": "2022-09-22T15:15:49.328006Z",
				  "events": [
					{
					  "envelope": {
						"event": {
						  "actor": "Actor",
						  "action": "Action",
						  "message": "sigtest100",
						  "new": "New",
						  "old": "Old",
						  "source": "Source",
						  "status": "Status",
						  "target": "Target"
						},
						"received_at": "2022-09-20T13:09:28.673562Z",
						"public_key": "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=",
						"signature": "dg7Wg+E8QzZzhECzQoH3v3pbjWObR8ve7SHREAyA9JlFOusKPHVb16t5D3rbscnv80ry/aWzfMTscRNSYJFzDA=="
					  },
					  "hash": "",
					  "leaf_index": 30,
					  "membership_proof": "l:c0bfb0fd1159f7f40c8b0e5f1ec28ebf3c7c7bbe41c8b9e62ee5f3238b1c51fa,l:edc77dec9297653dddf55e833ec9b415f2aa32d77a231408443a7d642504f9bb,l:17d7f0d7483acfdddadaef8941fe68af809d9be6c560a9277aad2c35fe958606,l:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce"
					}
				  ],
				  "root": {
					"url": "https://arweave.net/GHamEz43bRGY0oeGMT-3kB7K3U7WI4OY2g1y2RgUGcM",
					"published_at": "2022-09-20T13:30:33.280268Z",
					"size": 31,
					"root_hash": "58e83c3bed473694e34d714a5c71d78be3d2e6741fef6120c0108564a8c3519d",
					"consistency_proof": [
					  "x:a9e2809545d2e6a6a82ec636fd2c29bc84e3c063497f3f62356bf2c9fe7fcd2e,r:20bf7a8d010354fa3eacaf5d53d6b33a87ab23a7e6b4e1ac4cb712d5fca2a54a,l:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:cb0df5395b30583e928a2d779b101da997b8a25d2a162375ada3bdc8f6621f9c,r:7f50966e703039f135755e41afe8cb557941ff1431fa0c09d49d1ff2d7d906f3,l:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:26345b33ead978bf870990c8b4c2d116f4ed2c6de0802a4906d97c4504937824,r:51de5a887c3e09610693ac0514ae3ef53166bf7d1a774078dce1390c6228f940,l:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce",
					  "x:25cdb02dd291cf24068996c32b5848f6fa637327ecd94e7cc7f07562b0a997ce,r:029e187f6a9c5f51e6d44b40194de81e80c6288f08e3e00f59ea3b81fc092991"
					],
					"tree_name": "e4faf306ccb5e76f00430e203ef9ebb9dbf694f782fa17ca7d342c4802f031c7"
				  }
				}
			  }`)
	})

	client, _ := audit.New(pangeatesting.TestConfig(url), audit.WithLogProofVerificationEnabled())
	input := &audit.SearchInput{
		Query:   "message:test",
		Verbose: pangea.Bool(true),
	}
	ctx := context.Background()
	got, err := client.Search(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, got.Result.Count)
	assert.Equal(t, 1, len(got.Result.Events))
	assert.Equal(t, audit.NotVerified, got.Result.Events[0].MembershipVerification)
}

func TestSearchResults(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	t1 := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60))
	mux.HandleFunc("/v1/results", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"id":"some-id","limit":50}`)
		fmt.Fprintf(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result": {
					"count": 2,
					"events": [
						{
							"envelope": {
								"event": {
									"message": "test_2"
								}
							},
							"received_at": "%[1]v",
							"leaf_index": 2,
							"membership_proof": "some-proof"
						},
						{
							"envelope": {
								"event": {
									"message": "test_1"
								}
							},
							"received_at": "%[1]v",
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
		ID:    "some-id",
		Limit: 50,
	}
	ctx := context.Background()
	got, err := client.SearchResults(ctx, input)

	assert.NoError(t, err)

	want := &audit.SearchResultOutput{
		Count: 2,
		Events: audit.SearchEvents{
			{
				EventEnvelope: &audit.EventEnvelope{
					Event: &audit.Event{
						Message: "test_2",
					},
					ReceivedAt: got.Result.Events[1].EventEnvelope.ReceivedAt,
				},
				LeafIndex:       pangea.Int(2),
				MembershipProof: pangea.String("some-proof"),
				RawEnvelope:     got.Result.Events[0].RawEnvelope,
			},
			{
				EventEnvelope: &audit.EventEnvelope{
					Event: &audit.Event{
						Message: "test_1",
					},
					ReceivedAt: got.Result.Events[1].EventEnvelope.ReceivedAt,
				},
				LeafIndex:       pangea.Int(3),
				MembershipProof: pangea.String("some-proof"),
				RawEnvelope:     got.Result.Events[1].RawEnvelope,
			},
		},
		Root: &audit.Root{
			PublishedAt: &t1,
			RootHash:    "3a2563b40abe941f21c2ea929f2be92606fd2545762d3fa755ecec2942f5d452",
			Size:        11,
			ConsistencyProof: &[]string{
				"x:68810d719dc9dccee268d17a6c5baf3bf12d7ffad5673b763e06338121ed4e46,r:4a291c09b0bed8303d3e7f91315bd47da3df422151e642ded4208f46342a12f6",
				"x:82eba5aa211af097d22ecf215be386212c192d7068a02aeb4280905e4d200ff9,r:03c513b31ev80f4c871dbcd07b069fd369482529984f0770008f6c7777813026",
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
				"status": "Success",
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
		TreeSize: 11,
	}
	ctx := context.Background()
	got, err := client.Root(ctx, input)

	assert.NoError(t, err)

	want := &audit.RootOutput{
		Data: audit.Root{
			PublishedAt: &t1,
			RootHash:    "3a2563b40abe941f21c2ea929f2be92606fd2545762d3fa755ecec2942f5d452",
			Size:        11,
			ConsistencyProof: &[]string{
				"x:68810d719dc9dccee268d17a6c5baf3bf12d7ffad5673b763e06338121ed4e46,r:4a291c09b0bed8303d3e7f91315bd47da3df422151e642ded4208f46342a12f6",
				"x:82eba5aa211af097d22ecf215be386212c192d7068a02aeb4280905e4d200ff9,r:03c513b31ev80f4c871dbcd07b069fd369482529984f0770008f6c7777813026",
			},
		},
	}

	assert.Equal(t, want, got.Result)
}

func Test_BadDomain(t *testing.T) {
	_, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	cfg := pangeatesting.TestConfig(url)
	cfg.Domain = "fakedomain^^"
	client, _ := audit.New(cfg)
	event := audit.Event{
		Message: "test",
	}
	ctx := context.Background()
	got, err := client.Log(ctx, event, true)
	assert.Error(t, err)
	assert.Nil(t, got)

	got2, err := client.Search(ctx, &audit.SearchInput{})
	assert.Error(t, err)
	assert.Nil(t, got2)

	got3, err := client.SearchResults(ctx, &audit.SearchResultInput{})
	assert.Error(t, err)
	assert.Nil(t, got3)

	got4, err := client.Root(ctx, &audit.RootInput{})
	assert.Error(t, err)
	assert.Nil(t, got4)

}

func Test_NilRequest(t *testing.T) {
	_, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	cfg := pangeatesting.TestConfig(url)
	client, _ := audit.New(cfg)

	ctx := context.Background()

	got2, err := client.Search(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, got2)

	got3, err := client.SearchResults(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, got3)

	got4, err := client.Root(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, got4)

}

func TestLogError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.Log(context.Background(), audit.Event{}, true)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.Log", f)
}

func TestSearchError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.Search(context.Background(), &audit.SearchInput{})
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.Search", f)
}

func TestSearchResultsError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.SearchResults(context.Background(), &audit.SearchResultInput{})
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.SearchResults", f)
}

func TestRootError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := audit.New(cfg)
		_, err := client.Root(context.Background(), &audit.RootInput{})
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Audit.Root", f)
}

func TestFailedOptions(t *testing.T) {
	_, err := audit.New(
		pangeatesting.TestConfig("url"),
		audit.WithLogLocalSigning("bad file name"),
	)
	assert.Error(t, err)

	_, err = audit.New(
		pangeatesting.TestConfig("url"),
		audit.DisableEventVerification(),
	)
	assert.NoError(t, err)
}
