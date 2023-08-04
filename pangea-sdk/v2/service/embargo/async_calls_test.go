// go:build integration
package embargo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/embargo"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironmentAsync = pangeatesting.Live
)

func embargoIntegrationCfgAsync(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironmentAsync)
}

func Test_Integration_Async_CallAndWait(t *testing.T) {
	cfg := embargoIntegrationCfgAsync(t)
	client := embargo.New(cfg)
	defer client.Close()

	w := pangea.NewWorker(5)
	defer w.Close()

	// Call 1
	ctx1, cancelFn1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn1()
	input1 := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	p1 := pangea.NewPromise(client.ISOCheck, ctx1, input1)
	assert.NotNil(t, p1)
	w.Run(p1.Execute)

	p1.Wait()

	res, err := p1.Get()
	assert.NoError(t, err)
	assert.NotNil(t, res.Result)
	assert.NotZero(t, res.Result.Count)
	assert.GreaterOrEqual(t, len(res.Result.Sanctions), 1)
	assert.Equal(t, res.Result.Sanctions[0].EmbargoedCountryName, "Cuba")
}

func Test_Integration_Async_CancelCall(t *testing.T) {
	cfg := embargoIntegrationCfgAsync(t)
	client := embargo.New(cfg)
	defer client.Close()

	w := pangea.NewWorker(5)
	defer w.Close()

	// Call 1
	ctx1, cancelFn1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn1()
	input1 := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	p1 := pangea.NewPromise(client.ISOCheck, ctx1, input1)
	assert.NotNil(t, p1)
	w.Run(p1.Execute)

	time.Sleep(300 * time.Millisecond) // Wait for the call to be sent
	p1.Cancel()                        // Cancel the call

	res, err := p1.Get()
	assert.Error(t, err) // Should have an error due to cancellation
	fmt.Println(err)
	assert.Nil(t, res)
}

func Test_Integration_Async_MultipleCalls(t *testing.T) {
	cfg := embargoIntegrationCfgAsync(t)
	client := embargo.New(cfg)
	defer client.Close()

	w := pangea.NewWorker(5)
	defer w.Close()

	// Call 1
	ctx1, cancelFn1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn1()
	input1 := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	p1 := pangea.NewPromise(client.ISOCheck, ctx1, input1)
	assert.NotNil(t, p1)
	w.Run(p1.Execute)

	// Call 2
	ctx2, cancelFn2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn2()
	input2 := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	p2 := pangea.NewPromise(client.ISOCheck, ctx2, input2)
	assert.NotNil(t, p2)
	w.Run(p2.Execute)

	// Call 3
	ctx3, cancelFn3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn3()
	input3 := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	p3 := pangea.NewPromise(client.ISOCheck, ctx3, input3)
	assert.NotNil(t, p3)
	w.Run(p3.Execute)

	time.Sleep(100 * time.Millisecond) // Wait for the calls to be sent
	assert.Equal(t, uint(3), w.ThreadsRunning())
	assert.Equal(t, uint(0), w.ThreadsPending())
	assert.Equal(t, 3, client.GetNumRequestsInProgress())

	// Wait for the calls to finish
	client.WaitGroup()

	// Check the results p1
	res, err := p1.Get()
	assert.NoError(t, err)
	assert.NotNil(t, res.Result)
	assert.NotZero(t, res.Result.Count)
	assert.GreaterOrEqual(t, len(res.Result.Sanctions), 1)
	assert.Equal(t, res.Result.Sanctions[0].EmbargoedCountryName, "Cuba")

	// Check the results p2
	res, err = p2.Get()
	assert.NoError(t, err)
	assert.NotNil(t, res.Result)
	assert.NotZero(t, res.Result.Count)
	assert.GreaterOrEqual(t, len(res.Result.Sanctions), 1)
	assert.Equal(t, res.Result.Sanctions[0].EmbargoedCountryName, "Cuba")

	// Check the results p3
	res, err = p3.Get()
	assert.NoError(t, err)
	assert.NotNil(t, res.Result)
	assert.NotZero(t, res.Result.Count)
	assert.GreaterOrEqual(t, len(res.Result.Sanctions), 1)
	assert.Equal(t, res.Result.Sanctions[0].EmbargoedCountryName, "Cuba")
}
