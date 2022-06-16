package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Root struct {
	TreeName         *string   `json:"tree_name"`
	Size             *uint     `json:"size"`
	RootHash         *string   `json:"root_hash"`
	URL              *string   `json:"url"`
	PublishedAt      *string   `json:"published_at"`
	ConsistencyProof []*string `json:"consistency_proof"`
}

type RootOutput struct {
	Data *Root `json:"data"`
}

func (a *Audit) Root(ctx context.Context) (*RootOutput, *pangea.Response, error) {
	req, err := a.Client.NewRequest("POST", "audit", "v1/root", struct{}{})
	if err != nil {
		return nil, nil, err
	}
	var out RootOutput
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}
