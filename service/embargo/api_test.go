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

func TestCheck(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/check", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"iso_code":"CU"}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
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

	client, _ := embargo.New(pangeatesting.TestConfig(url))
	input := &embargo.CheckInput{
		ISOCode: pangea.String("CU"),
	}
	ctx := context.Background()
	got, _, err := client.Check(ctx, input)

	assert.NoError(t, err)

	sanction := &embargo.Sanction{
		ListName:                pangea.String("ITAR"),
		EmbargoedCountryName:    pangea.String("Cuba"),
		EmbargoedCountryISOCode: pangea.String("CU"),
		IssuingCountry:          pangea.String("US"),
		Annotations: map[string]interface{}{
			"reference": map[string]interface{}{
				"paragraph":  "d1",
				"regulation": "CFR 126.1",
			},
			"restriction_name": "ITAR",
		},
	}
	want := &embargo.CheckOutput{
		Sanctions: []*embargo.Sanction{sanction},
		Count:     pangea.Int(1),
	}
	assert.Equal(t, want, got)
}

func TestCheckError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := embargo.New(cfg)
		_, _, err := client.Check(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Embargo.Check", f)
}
