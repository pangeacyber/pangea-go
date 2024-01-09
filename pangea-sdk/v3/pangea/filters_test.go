// go:build unit

package pangea_test

import (
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/stretchr/testify/assert"
)

func TestFilterEqual(t *testing.T) {
	f := make(pangea.Filter)
	fe := pangea.NewFilterEqual[string]("name", &f)

	assert.Equal(t, 0, len(f))

	v := "value"
	fe.Set(&v)

	assert.Equal(t, 1, len(f))
	assert.Equal(t, v, *fe.Get())

	fe.Set(nil)
	assert.Equal(t, 0, len(f))
}

func TestFilterMatch(t *testing.T) {
	f := make(pangea.Filter)
	fm := pangea.NewFilterMatch[string]("name", &f)

	assert.Equal(t, 0, len(f))
	assert.Nil(t, fm.In())
	assert.Nil(t, fm.Contains())

	values := []string{"value1", "value2", "value3"}
	fm.SetIn(values)

	assert.Equal(t, 1, len(f))
	assert.Equal(t, values, fm.In())

	containsValues := []string{"contain1", "contain2", "contain3"}
	fm.SetContains(containsValues)
	assert.Equal(t, 2, len(f))
	assert.Equal(t, containsValues, fm.Contains())

	fm.SetContains(nil)
	fm.SetIn(nil)
	assert.Equal(t, 0, len(f))
}

func TestFilterRange(t *testing.T) {
	f := make(pangea.Filter)
	fr := pangea.NewFilterRange[int]("range", &f)

	assert.Equal(t, 0, len(f))

	// Ensure all parameters are initially nil
	assert.Nil(t, fr.LessThan())
	assert.Nil(t, fr.LessThanEqual())
	assert.Nil(t, fr.GreaterThan())
	assert.Nil(t, fr.GreaterThanEqual())

	lessThanValue := 10
	fr.SetLessThan(&lessThanValue)

	assert.Equal(t, 1, len(f))
	assert.Equal(t, &lessThanValue, fr.LessThan())

	lessThanEqualValue := 20
	fr.SetLessThanEqual(&lessThanEqualValue)

	assert.Equal(t, 2, len(f))
	assert.Equal(t, &lessThanEqualValue, fr.LessThanEqual())

	greaterThanValue := 30
	fr.SetGreaterThan(&greaterThanValue)

	assert.Equal(t, 3, len(f))
	assert.Equal(t, &greaterThanValue, fr.GreaterThan())

	greaterThanEqualValue := 40
	fr.SetGreaterThanEqual(&greaterThanEqualValue)

	assert.Equal(t, 4, len(f))
	assert.Equal(t, &greaterThanEqualValue, fr.GreaterThanEqual())

	// Test deleting values
	fr.SetLessThan(nil)
	assert.Equal(t, 3, len(f))
	assert.Nil(t, fr.LessThan())

	fr.SetLessThanEqual(nil)
	assert.Equal(t, 2, len(f))
	assert.Nil(t, fr.LessThanEqual())

	fr.SetGreaterThan(nil)
	assert.Equal(t, 1, len(f))
	assert.Nil(t, fr.GreaterThan())

	fr.SetGreaterThanEqual(nil)
	assert.Equal(t, 0, len(f))
	assert.Nil(t, fr.GreaterThanEqual())
}
