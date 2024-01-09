// go:build unit

package pangea_test

import (
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/stretchr/testify/assert"
)

func TestConvertTypes(t *testing.T) {
	pb := pangea.Bool(true)
	assert.NotNil(t, pb)
	assert.Equal(t, true, *pb)
	assert.Equal(t, true, pangea.BoolValue(pb))

	pi := pangea.Int(10)
	assert.NotNil(t, pi)
	assert.Equal(t, 10, *pi)
	assert.Equal(t, 10, pangea.IntValue(pi))

	ps := pangea.String("test")
	assert.NotNil(t, ps)
	assert.Equal(t, "test", *ps)
	assert.Equal(t, "test", pangea.StringValue(ps))

	pt := pangeautil.PangeaTimestamp(time.Now())
	ptp := pangea.PangeaTime(pt)
	assert.NotNil(t, ptp)

	ti := time.Now()
	tp := pangea.Time(ti)
	assert.NotNil(t, tp)
}

func TestStringify(t *testing.T) {
	o := struct {
		A string `json:"b"`
		B string `json:"a"`
	}{
		A: "some-string",
		B: "another-string",
	}
	s := pangea.Stringify(o)
	assert.NotNil(t, s)
	assert.NotEmpty(t, s)
}

func TestEncoding(t *testing.T) {
	m := "test"
	m64 := pangea.StrToB64(m)
	mr, err := pangea.B64ToStr(m64)
	assert.NoError(t, err)
	assert.Equal(t, m, string(mr))
}
