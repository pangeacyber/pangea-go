// go:build integration && !unit
package authn_test

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/authn"
	"github.com/stretchr/testify/assert"
)

var RANDOM_VALUE string
var USER_ID string
var EMAIL_TEST string
var EMAIL_DELETE string
var EMAIL_INVITE_DELETE string
var EMAIL_INVITE_KEEP string
var PROFILE_OLD = map[string]string{}
var PROFILE_NEW = map[string]string{}
var PASSWORD_OLD string
var PASSWORD_NEW string
var timeNow = time.Now()
var timeStr = timeNow.Format("yyyyMMdd_HHmmss")

const (
	testingEnvironment = pangeatesting.Live
)

func authnIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	RANDOM_VALUE = strconv.Itoa(rand.Intn(10000000))
	USER_ID = ""
	EMAIL_TEST = "user+test" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_DELETE = "user+delete" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_INVITE_DELETE = "user+invite_del" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_INVITE_KEEP = "user+invite_keep" + RANDOM_VALUE + "@pangea.cloud"
	PROFILE_OLD = authn.ProfileData{
		"name":    "User name",
		"country": "Argentina",
	}
	PROFILE_NEW = map[string]string{"age": "18"}

	PASSWORD_OLD = "My1s+Password"
	PASSWORD_NEW = "My1s+Password_new"

	// Run tests
	exitVal := m.Run()

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
	assert.NotEmpty(t, out.Result.ID)
	USER_ID = out.Result.ID

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

func Test_Integration_User_Login_And_User_Stuff(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	// Login with password
	input := authn.UserLoginPasswordRequest{
		Email:    EMAIL_TEST,
		Password: PASSWORD_OLD,
	}
	resp, err := client.User.Login.Password(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.ActiveToken)
	assert.NotEmpty(t, resp.Result.RefreshToken)

	// Change password
	input2 := authn.ClientPasswordChangeRequest{
		Token:       resp.Result.ActiveToken.Token,
		OldPassword: PASSWORD_OLD,
		NewPassword: PASSWORD_NEW,
	}
	resp2, err := client.Client.Password.Change(ctx, input2)
	assert.NoError(t, err)
	assert.Empty(t, resp2.Result)

	// Verify
	input3 := authn.UserVerifyRequest{
		IDProvider:    authn.IDPPassword,
		Email:         EMAIL_TEST,
		Authenticator: PASSWORD_NEW,
	}
	resp3, err := client.User.Verify(ctx, input3)
	assert.NoError(t, err)
	assert.Equal(t, USER_ID, resp3.Result.ID)

	// Check token
	input6 := authn.ClientTokenCheckRequest{
		Token: resp.Result.ActiveToken.Token,
	}
	resp6, err := client.Client.Token.Check(ctx, input6)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp6.Result)

	// Refresh token
	input4 := authn.ClientSessionRefreshRequest{
		UserToken:    resp.Result.ActiveToken.Token,
		RefreshToken: resp.Result.RefreshToken.Token,
	}
	resp4, err := client.Client.Session.Refresh(ctx, input4)
	assert.NoError(t, err)
	assert.NotNil(t, resp4.Result)
	assert.NotEmpty(t, resp4.Result.ActiveToken)
	assert.NotEmpty(t, resp4.Result.RefreshToken)

	// Reset password
	input5 := authn.UserPasswordResetRequest{
		UserID:      USER_ID,
		NewPassword: PASSWORD_NEW,
	}
	resp5, err := client.User.Password.Reset(ctx, input5)
	assert.NoError(t, err)
	assert.NotNil(t, resp5)

	// Client session logout
	input7 := authn.ClientSessionLogoutRequest{
		Token: resp4.Result.ActiveToken.Token,
	}
	resp7, err := client.Client.Session.Logout(ctx, input7)
	assert.NoError(t, err)
	assert.NotNil(t, resp7)

	// Re-Login with password
	input = authn.UserLoginPasswordRequest{
		Email:    EMAIL_TEST,
		Password: PASSWORD_NEW,
	}
	resp, err = client.User.Login.Password(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.ActiveToken)
	assert.NotEmpty(t, resp.Result.RefreshToken)

	// Session logout
	input8 := authn.SessionLogoutRequest{
		UserID: USER_ID,
	}
	resp8, err := client.Session.Logout(ctx, input8)
	assert.NoError(t, err)
	assert.NotNil(t, resp8)
}

func Test_Integration_User_Profile(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	// Get profile by email
	input := authn.UserProfileGetRequest{
		Email: pangea.String(EMAIL_TEST),
	}
	resp, err := client.User.Profile.Get(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_ID, resp.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Empty(t, resp.Result.Profile)

	// Get profile by user_id
	input = authn.UserProfileGetRequest{
		ID: pangea.String(USER_ID),
	}
	resp, err = client.User.Profile.Get(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_ID, resp.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Empty(t, resp.Result.Profile)

	// Update request by email
	input2 := authn.UserProfileUpdateRequest{
		Email:   pangea.String(EMAIL_TEST),
		Profile: PROFILE_OLD,
	}
	resp2, err := client.User.Profile.Update(ctx, input2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2.Result)
	assert.Equal(t, USER_ID, resp2.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp2.Result.Email)
	assert.Equal(t, authn.ProfileData(PROFILE_OLD), resp2.Result.Profile)

	// Update request by email
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

	assert.Equal(t, USER_ID, resp3.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp3.Result.Email)
	assert.Equal(t, authn.ProfileData(finalProfile), resp3.Result.Profile)

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
	assert.Equal(t, USER_ID, resp.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Equal(t, false, resp.Result.Disabled)
	assert.Equal(t, false, resp.Result.RequireMFA)
}

func Test_Integration_Client_Session(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	// Login with password
	input := authn.UserLoginPasswordRequest{
		Email:    EMAIL_TEST,
		Password: PASSWORD_NEW,
	}
	resp, err := client.User.Login.Password(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.ActiveToken)
	assert.NotEmpty(t, resp.Result.RefreshToken)

	filter := authn.NewFilterSessionList()

	// Client session list
	input2 := authn.ClientSessionListRequest{
		Token:  resp.Result.ActiveToken.Token,
		Filter: filter.Filter(),
	}
	resp2, err := client.Client.Session.List(ctx, input2)
	assert.NoError(t, err)
	assert.Greater(t, len(resp2.Result.Sessions), 0)

	for _, s := range resp2.Result.Sessions {
		input3 := authn.ClientSessionInvalidateRequest{
			Token:     resp.Result.ActiveToken.Token,
			SessionID: s.ID,
		}
		resp3, err := client.Client.Session.Invalidate(ctx, input3)
		assert.NoError(t, err)
		assert.NotNil(t, resp3)
	}
}

func Test_Integration_Session(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	// Login with password
	input := authn.UserLoginPasswordRequest{
		Email:    EMAIL_TEST,
		Password: PASSWORD_NEW,
	}
	resp, err := client.User.Login.Password(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.ActiveToken)
	assert.NotEmpty(t, resp.Result.RefreshToken)

	// Client session list
	filter := authn.NewFilterSessionList()
	input2 := authn.SessionListRequest{
		Filter: filter.Filter(),
	}
	resp2, err := client.Session.List(ctx, input2)
	assert.NoError(t, err)
	assert.Greater(t, len(resp2.Result.Sessions), 0)

	for _, s := range resp2.Result.Sessions {
		input3 := authn.SessionInvalidateRequest{
			SessionID: s.ID,
		}
		resp3, err := client.Session.Invalidate(ctx, input3)
		assert.NoError(t, err)
		assert.NotNil(t, resp3)
	}
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

func Test_Integration_User_Invite_List(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	filter := authn.NewFilterUserInviteList()
	resp, err := client.User.Invites.List(ctx, authn.UserInviteListRequest{
		Filter: filter.Filter(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Greater(t, len(resp.Result.Invites), 0)
}

func Test_Integration_User_List(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	filter := authn.NewFilterUserList()
	input := authn.UserListRequest{
		Filter: filter.Filter(),
	}
	resp, err := client.User.List(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Greater(t, len(resp.Result.Users), 0)

	for _, u := range resp.Result.Users {
		input := authn.UserDeleteRequest{
			ID: u.ID,
		}
		rDel, err := client.User.Delete(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, rDel)
	}
}

func agreementsCycle(t *testing.T, client *authn.AuthN, ctx context.Context, at authn.AgreementType) {
	name := string(at) + timeStr
	text := "This is agreement text"
	active := false

	// Create
	cr, err := client.Agreements.Create(ctx, authn.AgreementCreateRequest{
		Type:   at,
		Name:   name,
		Text:   text,
		Active: pangea.Bool(active),
	})
	assert.NoError(t, err)
	assert.NotNil(t, cr)
	assert.NotNil(t, cr.Result)
	assert.Equal(t, name, cr.Result.Name)
	assert.Equal(t, text, cr.Result.Text)
	assert.Equal(t, active, cr.Result.Active)
	assert.NotEmpty(t, cr.Result.ID)
	id := cr.Result.ID

	// Update agreement
	newName := name + "v2"
	newText := text + "v2"

	ur, err := client.Agreements.Update(ctx, authn.AgreementUpdateRequest{
		ID:     id,
		Type:   at,
		Text:   &newText,
		Name:   &newName,
		Active: pangea.Bool(active),
	})
	assert.NoError(t, err)
	assert.NotNil(t, ur)
	assert.NotNil(t, ur.Result)
	assert.Equal(t, newName, ur.Result.Name)
	assert.Equal(t, newText, ur.Result.Text)
	assert.Equal(t, active, ur.Result.Active)

	filter := authn.NewFilterAgreementList()

	// List
	lr, err := client.Agreements.List(ctx, authn.AgreementListRequest{
		Filter: filter.Filter(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, lr)
	assert.NotNil(t, lr.Result)
	assert.Greater(t, lr.Result.Count, 0)
	assert.Greater(t, len(lr.Result.Agreements), 0)
	count := lr.Result.Count

	// delete
	dr, err := client.Agreements.Delete(ctx, authn.AgreementDeleteRequest{
		Type: at,
		ID:   id,
	})
	assert.NoError(t, err)
	assert.NotNil(t, dr)
	assert.NotNil(t, dr.Result)

	// List again
	lr2, err := client.Agreements.List(ctx, authn.AgreementListRequest{
		Filter: filter.Filter(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, lr2)
	assert.NotNil(t, lr2.Result)
	assert.Equal(t, count-1, lr2.Result.Count)
}

func Test_Integration_AgreementsCycleEULA(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	agreementsCycle(t, client, ctx, authn.ATeula)
}

func Test_Integration_AgreementsCyclePP(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	agreementsCycle(t, client, ctx, authn.ATprivacyPolicy)
}
