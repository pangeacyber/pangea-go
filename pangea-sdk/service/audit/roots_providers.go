package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/arweave"
)

type RootsProvider interface {
	UpdateRoots(ctx context.Context, treeSizes []string) map[int]Root
	OverrideRoots(roots map[int]Root) map[int]Root
}

type ArweaveRootsProvider struct {
	TreeName string
	Client   *arweave.Arweave
	Roots    map[int]Root
}

func NewArweaveRootsProvider(treeName string) *ArweaveRootsProvider {
	return &ArweaveRootsProvider{
		TreeName: treeName,
		Client:   arweave.New(),
		Roots:    make(map[int]Root),
	}
}

func (rp *ArweaveRootsProvider) UpdateRoots(ctx context.Context, treeSizes []string) map[int]Root {
	newTreeSizes := make([]string, 0)
	// Request only new tree sizes roots
	for _, s := range treeSizes {
		n, err := strconv.Atoi(s)
		if err == nil {
			if _, ok := rp.Roots[n]; !ok {
				// If not present in map, add to newTreeSizes to make request
				newTreeSizes = append(newTreeSizes, s)
			}
		}
	}

	if len(newTreeSizes) > 0 {
		tags := arweave.TagFilters{
			arweave.TagFilter{
				Name:   "tree_size",
				Values: newTreeSizes,
			},
			arweave.TagFilter{
				Name:   "tree_name",
				Values: []string{rp.TreeName},
			},
		}
		resp, err := rp.Client.TransactionConnectionByTags(ctx, tags)
		if err != nil {
			return rp.Roots
		}

		if resp == nil || resp.Data == nil || resp.Data.Transactions == nil || resp.Data.Transactions.Edges == nil {
			return rp.Roots
		}

		for _, edge := range resp.Data.Transactions.Edges {
			if edge == nil {
				continue
			}
			root := Root{}
			err = rp.fetchTransaction(ctx, *edge.Node.ID, &root)
			if err != nil {
				continue
			}
			size, err := treeSizeFromTransaction(edge.Node)
			if err != nil {
				continue
			}
			rp.Roots[size] = root
		}
	}
	return rp.Roots
}

func (rp *ArweaveRootsProvider) OverrideRoots(roots map[int]Root) map[int]Root {
	for treeSize, root := range roots {
		rp.Roots[treeSize] = root
	}
	return rp.Roots
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
