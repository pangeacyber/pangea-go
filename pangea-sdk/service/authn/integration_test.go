// go:build integration && !unit
package authn_test

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/authn"
	"github.com/stretchr/testify/assert"
)

// rand.Seed(time.Now().UnixNano())
var RANDOM_VALUE string
var USER_IDENTITY string
var EMAIL_TEST string
var EMAIL_DELETE string
var EMAIL_INVITE_DELETE string
var EMAIL_INVITE_KEEP string
var PROFILE_OLD = map[string]string{}
var PROFILE_NEW = map[string]string{}
var PASSWORD_OLD string
var PASSWORD_NEW string

const (
	testingEnvironment = pangeatesting.Develop
)

func authnIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func TestMain(m *testing.M) {
	// Write code here to run before tests
	rand.Seed(time.Now().UnixNano())
	RANDOM_VALUE = strconv.Itoa(rand.Intn(10000000))
	USER_IDENTITY = ""
	EMAIL_TEST = "andres.tournour+test" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_DELETE = "andres.tournour+delete" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_INVITE_DELETE = "andres.tournour+invite_del" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_INVITE_KEEP = "andres.tournour+invite_keep" + RANDOM_VALUE + "@pangea.cloud"
	PROFILE_OLD = map[string]string{
		"name":    "User name",
		"country": "Argentina",
	}
	PROFILE_NEW = map[string]string{"age": "18"}

	PASSWORD_OLD = "My1s+Password"
	PASSWORD_NEW = "My1s+Password_new"

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func Test_Integration_User_Create(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserCreateRequest{
		Email:         EMAIL_TEST,
		Authenticator: PASSWORD_OLD,
		IDProvider:    authn.IDPPassword,
	}
	out, err := client.User.Create(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Identity)
	USER_IDENTITY = out.Result.Identity

	input = authn.UserCreateRequest{
		Email:         EMAIL_DELETE,
		Authenticator: PASSWORD_OLD,
		IDProvider:    authn.IDPPassword,
	}
	out, err = client.User.Create(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)

}

func Test_Integration_User_Delete(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserDeleteRequest{
		Email: EMAIL_DELETE,
	}
	out, err := client.User.Delete(ctx, input)
	assert.NoError(t, err)
	assert.Empty(t, out.Result)
}

func Test_Integration_Password_Update(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.PasswordUpdateRequest{
		Email:     EMAIL_TEST,
		OldSecret: PASSWORD_OLD,
		NewSecret: PASSWORD_NEW,
	}
	out, err := client.Password.Update(ctx, input)
	assert.NoError(t, err)
	assert.Empty(t, out.Result)
}

func Test_Integration_User_Login(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserLoginRequest{
		Email:  EMAIL_TEST,
		Secret: PASSWORD_NEW,
	}
	resp, err := client.User.Login(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_IDENTITY, resp.Result.Identity)
}

func Test_Integration_User_Profile(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserProfileGetRequest{
		Email: pangea.String(EMAIL_TEST),
	}
	resp, err := client.User.Profile.Get(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_IDENTITY, resp.Result.Identity)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Empty(t, resp.Result.Profile)

	input = authn.UserProfileGetRequest{
		Identity: pangea.String(USER_IDENTITY),
	}
	resp, err = client.User.Profile.Get(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_IDENTITY, resp.Result.Identity)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Empty(t, resp.Result.Profile)

	input2 := authn.UserProfileUpdateRequest{
		Email:   pangea.String(EMAIL_TEST),
		Profile: PROFILE_OLD,
	}

	resp2, err := client.User.Profile.Update(ctx, input2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2.Result)
	assert.Equal(t, USER_IDENTITY, resp2.Result.Identity)
	assert.Equal(t, EMAIL_TEST, resp2.Result.Email)
	assert.Equal(t, PROFILE_OLD, resp2.Result.Profile)

	input3 := authn.UserProfileUpdateRequest{
		Email:   pangea.String(EMAIL_TEST),
		Profile: PROFILE_NEW,
	}

	resp3, err := client.User.Profile.Update(ctx, input3)
	assert.NoError(t, err)
	assert.NotNil(t, resp3.Result)

	var finalProfile = map[string]string{}
	for k, v := range PROFILE_OLD {
		finalProfile[k] = v
	}
	for k, v := range PROFILE_NEW {
		finalProfile[k] = v
	}

	assert.Equal(t, USER_IDENTITY, resp3.Result.Identity)
	assert.Equal(t, EMAIL_TEST, resp3.Result.Email)
	assert.Equal(t, finalProfile, resp3.Result.Profile)

}

func Test_Integration_User_Update(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserUpdateRequest{
		Email:      pangea.String(EMAIL_TEST),
		Disabled:   pangea.Bool(false),
		RequireMFA: pangea.Bool(false),
	}
	resp, err := client.User.Update(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_IDENTITY, resp.Result.Identity)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Equal(t, false, resp.Result.Disabled)
	assert.Equal(t, false, resp.Result.RequireMFA)
}

func Test_Integration_User_Invite(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	input := authn.UserInviteRequest{
		Inviter:  EMAIL_TEST,
		Email:    EMAIL_INVITE_KEEP,
		Callback: "https://someurl.com/callbacklink",
		State:    "Somestate",
	}
	resp, err := client.User.Invite(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, EMAIL_INVITE_KEEP, resp.Result.Email)
	assert.Equal(t, EMAIL_TEST, resp.Result.Inviter)

	input2 := authn.UserInviteRequest{
		Inviter:  EMAIL_TEST,
		Email:    EMAIL_INVITE_DELETE,
		Callback: "https://someurl.com/callbacklink",
		State:    "Somestate",
	}
	resp2, err := client.User.Invite(ctx, input2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2.Result)
	assert.Equal(t, EMAIL_INVITE_DELETE, resp2.Result.Email)
	assert.Equal(t, EMAIL_TEST, resp2.Result.Inviter)

	input3 := authn.UserInviteDeleteRequest{
		ID: resp2.Result.ID,
	}
	resp3, err := client.User.Invites.Delete(ctx, input3)
	assert.NoError(t, err)
	assert.Empty(t, resp3.Result)
}