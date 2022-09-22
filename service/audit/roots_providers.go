package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pangeacyber/go-pangea/internal/arweave"
)

type RootsProvider interface {
	Roots(ctx context.Context, treeSizes []string) (map[int]Root, error)
}

type ArweaveRootsProvider struct {
	TreeName string
	Client   *arweave.Arweave
}

func NewArweaveRootsProvider(treeName string) *ArweaveRootsProvider {
	return &ArweaveRootsProvider{
		TreeName: treeName,
		Client:   arweave.New(),
	}
}

func (s *ArweaveRootsProvider) Roots(ctx context.Context, treeSizes []string) (map[int]Root, error) {
	tags := arweave.TagFilters{
		{
			Name:   "tree_size",
			Values: treeSizes,
		},
		{
			Name:   "tree_name",
			Values: []string{s.TreeName},
		},
	}
	resp, err := s.Client.TransactionConnectionByTags(ctx, tags)
	if err != nil {
		return nil, err
	}
	roots := make(map[int]Root)
	for _, edge := range resp.Data.Transactions.Edges {
		root := Root{}
		err = s.fetchTransaction(ctx, *edge.Node.ID, &root)
		if err != nil {
			return nil, err
		}
		size, err := treeSizeFromTransaction(edge.Node)
		if err != nil {
			return nil, err
		}
		roots[size] = root
	}
	return roots, nil
}

func treeSizeFromTransaction(tx *arweave.Transaction) (int, error) {
	if tx == nil || len(tx.Tags) == 0 {
		return 0, fmt.Errorf("audit: empty transaction")
	}
	for _, tag := range tx.Tags {
		if tag.Name != nil && *tag.Name == "tree_size" && tag.Value != nil {
			res, err := strconv.Atoi(*tag.Value)
			if err != nil {
				return 0, fmt.Errorf("audit: invalid tree size: %w", err)
			}
			return res, nil
		}
	}
	return 0, fmt.Errorf("audit: missing tree size in transaction tags")
}

func (s *ArweaveRootsProvider) fetchTransaction(ctx context.Context, txID string, target *Root) error {
	resp, err := s.Client.TransactionByID(ctx, txID)
	if err != nil {
		return err
	}
	if len(resp) > 0 {
		return json.Unmarshal(resp, target)
	} else {
		return nil
	}
}
