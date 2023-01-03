// go:build integration && !unit
package authn_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/authn"
	"github.com/stretchr/testify/assert"
)

var RANDOM_VALUE = strconv.Itoa(rand.Intn(10000000))
var USER_IDENTITY = ""
var EMAIL_TEST = "andres.tournour+test" + RANDOM_VALUE + "@pangea.cloud"
var EMAIL_DELETE = "andres.tournour+delete" + RANDOM_VALUE + "@pangea.cloud"
var EMAIL_INVITE_DELETE = "andres.tournour+invite_del" + RANDOM_VALUE + "@pangea.cloud"
var EMAIL_INVITE_KEEP = "andres.tournour+invite_keep" + RANDOM_VALUE + "@pangea.cloud"

const (
	testingEnvironment = pangeatesting.Develop
	PASSWORD_OLD       = "My1s+Password"
	PASSWORD_NEW       = "My1s+Password_new"
	// PROFILE_OLD = {"name": "User name", "country": "Argentina"}
	// PROFILE_NEW = {"age": "18"}
)

func authnIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func Test_Integration___(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserCreateRequest{
		Email: EMAIL_TEST,
	}
	out, err := client.User.Create(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
}
