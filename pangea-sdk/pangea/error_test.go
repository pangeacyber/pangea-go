//go:build unit

package pangea_test

import (
	"errors"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	ue := pangea.NewUnmarshalError(errors.New("Error test"), make([]byte, 10), nil)
	assert.NotEmpty(t, ue.Error())

}
