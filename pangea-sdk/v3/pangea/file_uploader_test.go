//go:build unit

package pangea_test

import (
	"os"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/stretchr/testify/assert"
)

// "crc32c": "754995fb",
// "sha256": "81655950d560e804a6315e09e74a7414e7b18ba99f722abe6122857e69a3aebd",
// "size": 10028,

var ZERO_BYTES_FILEPATH = "./testdata/zerobytes.txt"
var PDF_FILEPATH = "./testdata/testfile.pdf"

func TestGetFileParams_NonEmptyFile(t *testing.T) {
	file, err := os.Open(PDF_FILEPATH)
	assert.NoError(t, err)
	assert.NotNil(t, file)

	params, err := pangea.GetUploadFileParams(file)
	assert.NoError(t, err)
	assert.Equal(t, params.CRC32C, "754995fb")
	assert.Equal(t, params.SHA256, "81655950d560e804a6315e09e74a7414e7b18ba99f722abe6122857e69a3aebd")
	assert.Equal(t, params.Size, 10028)

}

func TestGetFileParams_EmptyFile(t *testing.T) {
	file, err := os.Open(ZERO_BYTES_FILEPATH)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	params, err := pangea.GetUploadFileParams(file)
	assert.NoError(t, err)
	assert.Equal(t, params.CRC32C, "00000000")
	assert.Equal(t, params.SHA256, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	assert.Equal(t, params.Size, 0)
}
