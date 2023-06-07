package pangeautil_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeautil"
	"github.com/stretchr/testify/assert"
)

func TestCanonicalizeJSONMarshall_Given_Unsorted_Struct_Fields_Returns_Json_With_Keys_SortedBy_Json_Tags(t *testing.T) {
	input := struct {
		A string `json:"b"`
		B string `json:"a"`
	}{
		A: "some-string",
		B: "another-string",
	}

	b := pangeautil.CanonicalizeJSONMarshall(input)
	assert.Equal(t, `{"a":"another-string","b":"some-string"}`, string(b))
}

func TestCanonicalizeJSONMarshall_Given_Unsorted_Struct_Fields_Returns_Json_With_Keys_SortedBy_Json_Tags_2(t *testing.T) {
	input := struct {
		A string  `json:"b"`
		B *string `json:"a"`
	}{
		A: "some-string",
		B: nil,
	}

	b := pangeautil.CanonicalizeJSONMarshall(input)
	assert.Equal(t, `{"b":"some-string"}`, string(b))
}

func TestCanonicalizeJSONMarshall_Given_StrPtr_With_Value_It_Returns_Value(t *testing.T) {
	str := "some-string"
	input := struct {
		B *string `json:"a"`
	}{
		B: &str,
	}

	b := pangeautil.CanonicalizeJSONMarshall(input)
	assert.Equal(t, `{"a":"some-string"}`, string(b))
}

func TestCanonicalizeJSONMarshall_Given_NilStrPtr_With_Value_It_Returns_Value(t *testing.T) {
	input := struct {
		B *string `json:"a"`
	}{
		B: nil,
	}

	b := pangeautil.CanonicalizeJSONMarshall(input)
	assert.Equal(t, `{}`, string(b))
}

func TestCanonicalizeJSONMarshall_Given_EmptyStruct_It_Returns_Empty_Json(t *testing.T) {
	input := struct{}{}
	b := pangeautil.CanonicalizeJSONMarshall(input)
	assert.Equal(t, `{}`, string(b))
}

func TestCanonicalizeJSONMarshall_Given_PtrEmptyStruct_It_Returns_Empty_Json(t *testing.T) {
	input := struct{}{}
	b := pangeautil.CanonicalizeJSONMarshall(&input)
	assert.Equal(t, `{}`, string(b))
}

func TestCanonicalizeJSONMarshall_Given_UnTagged_Struct_Fields_It_Returns_Empty_Json(t *testing.T) {
	input := struct {
		B string
	}{
		B: "other-string",
	}
	b := pangeautil.CanonicalizeJSONMarshall(&input)
	assert.Equal(t, `{}`, string(b))
}

func TestCanonicalizeStruct(t *testing.T) {
	input := struct {
		A string `json:"b"`
		B string `json:"a"`
	}{
		A: "some-string",
		B: "another-string",
	}

	b, err := pangeautil.CanonicalizeStruct(input)
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestPangeaTimestamp(t *testing.T) {
	pt := pangeautil.PangeaTimestamp(time.Now())
	ptb, err := pt.MarshalJSON()
	assert.NoError(t, err)

	ptu := pangeautil.PangeaTimestamp{}
	err = json.Unmarshal(ptb, &ptu)
	assert.NoError(t, err)
	ptt := time.Time(pt)
	ptut := time.Time(ptu)
	assert.True(t, ptt.Equal(ptut))
}
