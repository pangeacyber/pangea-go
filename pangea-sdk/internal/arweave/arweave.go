package arweave

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/defaults"
)

const baseURL = "https://arweave.net"

type Arweave struct {
	*http.Client
}

func New() *Arweave {
	return &Arweave{
		Client: defaults.HTTPClient(),
	}
}

// TransactionByID returns a transaction by its ID already decoded.
func (a *Arweave) TransactionByID(ctx context.Context, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/tx/%s/data", baseURL, id)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("arweave: failed GET %v: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arweave: GET %v with status code %v", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("arweave: GET %v failed to read response body: %w", url, err)
	}

	dbuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(body)))
	n, err := base64.RawURLEncoding.Decode(dbuf, body)
	if err != nil {
		return nil, fmt.Errorf("arweave: GET %v failed to decode response body: %w", url, err)
	}
	return dbuf[:n], nil
}

func (a *Arweave) TransactionConnectionByTags(ctx context.Context, tags TagFilters) (*TransactionConnectionResponse, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("arweave: missing tags")
	}
	query := fmt.Sprintf(`{
		transactions(tags: %v) {
			edges {
				node {
					id
					tags {
						name
						value
					}
				}
			}
		}
	}`, tags.GraphqlInput())

	url := fmt.Sprintf("%s/graphql", baseURL)
	q, _ := json.Marshal(GraphQLRequest{Query: query})
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(q))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("arweave: failed POST %v: %w", url, err)
	}
	defer resp.Body.Close()

	// A graphql API should always return a 200 status code
	// but we'll check for a non-200 status code anyway
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arweave: POST %v with status code %v", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("arweave: POST %v failed to read response body: %w", url, err)
	}
	var response TransactionConnectionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("arweave: POST %v failed to unmarshal response: %w", url, err)
	}
	if response.Err != nil {
		return nil, fmt.Errorf("arweave: POST %v failed with error: %w", url, response.Err)
	}
	return &response, nil
}
