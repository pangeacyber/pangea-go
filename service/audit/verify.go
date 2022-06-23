package audit

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/pangea/hash"
)

type proofSide uint

const (
	Left proofSide = iota
	Right
)

type proofItem struct {
	Side proofSide
	Hash hash.Hash
}

type proof []proofItem

func DecodeProof(s string) (proof, error) {
	items := strings.Split(s, ",")
	if len(items) == 0 {
		return nil, nil
	}

	p := make(proof, 0, len(items))
	for _, item := range items {
		side, h, err := hashFromHashProof(item)
		if err != nil {
			return nil, err
		}
		pi := proofItem{
			Hash: h,
		}
		switch side {
		case "l":
			pi.Side = Left
		case "r":
			pi.Side = Right
		default:
			return nil, fmt.Errorf("audit: inconsistent proof item with no side declaration: %v", item)
		}
		p = append(p, pi)
	}
	return p, nil
}

func VerifyMembershipProof(root Root, auditOutput AuditRecord, required bool) (bool, error) {
	membershipProof := pangea.StringValue(auditOutput.MembershipProof)
	if membershipProof == "" {
		return !required, nil
	}
	targetHash, err := hash.Decode(pangea.StringValue(auditOutput.Hash))
	if err != nil {
		return false, err
	}
	rootHash, err := hash.Decode(pangea.StringValue(root.RootHash))
	if err != nil {
		return false, err
	}
	proof, err := DecodeProof(membershipProof)
	if err != nil {
		return false, err
	}
	return verifyLogProof(targetHash, rootHash, proof), nil
}

func verifyLogProof(target, root hash.Hash, p proof) bool {
	h := target
	for _, proofItem := range p {
		switch proofItem.Side {
		case Left:
			h = hash.Pair(proofItem.Hash).With(h)
		case Right:
			h = hash.Pair(h).With(proofItem.Hash)
		}
	}
	return root.Equal(h)
}

func VerifyConsistencyProof(publishedRoots map[int]Root, record AuditRecord, required bool) bool {
	if record.LeafIndex == nil {
		return !required
	}
	idx := *record.LeafIndex
	if idx <= 1 {
		return !required
	}
	current, ok := publishedRoots[idx]
	if !ok {
		return false
	}
	previous, ok := publishedRoots[idx-1]
	if !ok {
		return false
	}
	verified, err := verifyConsistencyProof(previous, current)
	if err != nil {
		return false
	}
	return verified
}

func verifyConsistencyProof(old, new Root) (bool, error) {
	oldHash, err := hash.Decode(pangea.StringValue(old.RootHash))
	if err != nil {
		return false, err
	}
	newHash, err := hash.Decode(pangea.StringValue(new.RootHash))
	if err != nil {
		return false, err
	}
	consistencyProof, err := decodeRootProof(new.ConsistencyProof)
	if err != nil {
		return false, err
	}
	if len(consistencyProof) == 0 {
		return false, fmt.Errorf("audit: consistency proof is empty")
	}
	rootHash := consistencyProof[0].Hash
	for i := 1; i < len(consistencyProof); i++ {
		rootHash = hash.Pair(consistencyProof[i].Hash).With(rootHash)
	}
	if !rootHash.Equal(oldHash) {
		return false, nil
	}
	for _, item := range consistencyProof {
		if !verifyLogProof(item.Hash, newHash, item.Proof) {
			return false, nil
		}
	}
	return true, nil
}

// leaf index no existe no se valida
// devolver una lista de los audits que no se validaron
// devolver un root
func VerifyAuditRecords(ctx context.Context, rp RootsProvider, root *Root, records AuditRecords, required bool) (bool, error) {
	if root == nil || len(records) == 0 {
		return false, fmt.Errorf("audit: empty root or records")
	}

	treeSizes := treeSizes(root, records)
	roots, err := rp.Roots(ctx, treeSizes)
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if !VerifyConsistencyProof(roots, *record, required) {
			return false, fmt.Errorf("audit: consistency proof failed for record %v", *record)
		}
	}
	return true, nil
}

func VerifyAuditRecordsWithArweave(ctx context.Context, searchOutput *SearchOutput, required bool) (bool, error) {
	arweavecli := NewArweaveRootsProvider(*searchOutput.Root.TreeName)
	return VerifyAuditRecords(ctx, arweavecli, searchOutput.Root, searchOutput.Audits, required)
}

func treeSizes(root *Root, records AuditRecords) []string {
	treeSizes := make(map[int]struct{}, 0)
	treeSizes[*root.Size] = struct{}{}
	for _, record := range records {
		leafIdx := *record.LeafIndex
		treeSizes[leafIdx] = struct{}{}
		if leafIdx > 1 {
			treeSizes[leafIdx-1] = struct{}{}
		}
	}
	sizes := make([]string, 0, len(treeSizes))
	for size := range treeSizes {
		sizes = append(sizes, strconv.Itoa(size))
	}
	return sizes
}
