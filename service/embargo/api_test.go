package embargo_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
	"github.com/stretchr/testify/assert"
)

func TestISOCheck(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/iso/check", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"iso_code":"CU"}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result":{
					"sanctions": [
						{
							"list_name": "ITAR",
							"embargoed_country_name": "Cuba",
							"embargoed_country_iso_code": "CU",
							"issuing_country": "US",
							"annotations": {
								"reference": {
									"paragraph": "d1",
									"regulation": "CFR 126.1"
								},
								"restriction_name": "ITAR"
							}
						}
					],
					"count": 1
				},
				"summary": "success"
			}`)
	})

	client := embargo.New(pangeatesting.TestConfig(url))
	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("CU"),
	}
	ctx := context.Background()
	got, err := client.ISOCheck(ctx, input)

	assert.NoError(t, err)

	sanction := embargo.Sanction{
		ListName:                "ITAR",
		EmbargoedCountryName:    "Cuba",
		EmbargoedCountryISOCode: "CU",
		IssuingCountry:          "US",
		Annotations: map[string]interface{}{
			"reference": map[string]interface{}{
				"paragraph":  "d1",
				"regulation": "CFR 126.1",
			},
			"restriction_name": "ITAR",
		},
	}
	want := &embargo.CheckOutput{
		Sanctions: []embargo.Sanction{sanction},
		Count:     1,
	}
	assert.Equal(t, want, got.Result)
}

func TestIPCheck(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/ip/check", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"ip":"200.0.16.2"}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result":{
					"sanctions": [
						{
							"list_name": "ITAR",
							"embargoed_country_name": "Cuba",
							"embargoed_country_iso_code": "CU",
							"issuing_country": "US",
							"annotations": {
								"reference": {
									"paragraph": "d1",
									"regulation": "CFR 126.1"
								},
								"restriction_name": "ITAR"
							}
						}
					],
					"count": 1
				},
				"summary": "success"
			}`)
	})

	client := embargo.New(pangeatesting.TestConfig(url))
	input := &embargo.IPCheckInput{
		IP: pangea.String("200.0.16.2"),
	}
	ctx := context.Background()
	got, err := client.IPCheck(ctx, input)

	assert.NoError(t, err)

	sanction := embargo.Sanction{
		ListName:                "ITAR",
		EmbargoedCountryName:    "Cuba",
		EmbargoedCountryISOCode: "CU",
		IssuingCountry:          "US",
		Annotations: map[string]interface{}{
			"reference": map[string]interface{}{
				"paragraph":  "d1",
				"regulation": "CFR 126.1",
			},
			"restriction_name": "ITAR",
		},
	}
	want := &embargo.CheckOutput{
		Sanctions: []embargo.Sanction{sanction},
		Count:     1,
	}
	assert.Equal(t, want, got.Result)
}

func TestCheckISOError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client := embargo.New(cfg)
		_, err := client.ISOCheck(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Embargo.Check", f)
}

func TestCheckIPError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client := embargo.New(cfg)
		_, err := client.IPCheck(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Embargo.Check", f)
}
