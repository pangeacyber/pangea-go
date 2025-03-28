//go:build unit

package pangea_test

import (
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	m := "test"
	h := pangea.HashSHA256(m)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", h)
	assert.Equal(t, "9f86d", pangea.GetHashPrefix(h, 5))
}
