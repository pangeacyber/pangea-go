// go:build integration && !unit
package authn_test

import (
	"context"
	"fmt"
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
var CB_URI = "https://someurl.com/callbacklink"

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
		"first_name": "Name",
		"last_name":  "User",
	}
	PROFILE_NEW = map[string]string{"first_name": "NameUpdate"}

	PASSWORD_OLD = "My1s+Password"
	PASSWORD_NEW = "My1s+Password_new"

	// Run tests
	exitVal := m.Run()

	os.Exit(exitVal)
}

func flowHandlePasswordPhase(t *testing.T, ctx context.Context, client *authn.AuthN, flow_id, password string) *authn.FlowUpdateResult {
	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flow_id,
		Choice: authn.FCPassword,
		Data: authn.FlowUpdateDataPassword{
			Password: password,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	return resp.Result
}

func flowHandleProfilePhase(t *testing.T, ctx context.Context, client *authn.AuthN, flow_id string) *authn.FlowUpdateResult {
	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flow_id,
		Choice: authn.FCProfile,
		Data: authn.FlowUpdateDataProfile{
			Profile: PROFILE_OLD,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	return resp.Result
}

func flowHandleAgreementsPhase(t *testing.T, ctx context.Context, client *authn.AuthN, flow_id string, result *authn.FlowUpdateResult) *authn.FlowUpdateResult {
	// Iterate over flow_choices in response.result
	agreed := []string{}
	for _, flowChoice := range result.FlowChoices {
		// Check if the choice is AGREEMENTS
		if flowChoice.Choice == string(authn.FCAgreements) {
			// Assuming flowChoice.Data["agreements"] is a map[string]interface{}
			agreements, ok := flowChoice.Data["agreements"].(map[string]interface{})
			if ok {
				// Iterate over agreements and append the "id" values to agreed slice
				for _, v := range agreements {
					agreement, ok := v.(map[string]interface{})
					if ok {
						id, ok := agreement["id"].(string)
						if ok {
							agreed = append(agreed, id)
						}
					}
				}
			}
		}
	}

	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flow_id,
		Choice: authn.FCAgreements,
		Data: authn.FlowUpdateDataAgreements{
			Agreed: agreed,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	return resp.Result
}

func choiceIsAvailable(choices []authn.FlowChoiceItem, choice string) bool {
	for _, fc := range choices {
		if fc.Choice == choice {
			return true

		}
	}
	return false
}

func CreateAndLogin(t *testing.T, email, password string) *authn.FlowCompleteResult {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	fsresp, err := client.Flow.Start(ctx,
		authn.FlowStartRequest{
			Email:     email,
			FlowTypes: []authn.FlowType{authn.FTsignup, authn.FTsignin},
			CBURI:     CB_URI,
		})

	assert.NoError(t, err)
	assert.NotNil(t, fsresp)
	assert.NotNil(t, fsresp.Result)
	flowID := fsresp.Result.FlowID
	var result *authn.FlowUpdateResult = nil
	flowPhase := "initial"
	choices := fsresp.Result.FlowChoices

	for flowPhase != "phase_completed" {
		if choiceIsAvailable(choices, string(authn.FCPassword)) {
			result = flowHandlePasswordPhase(t, ctx, client, flowID, password)
		} else if choiceIsAvailable(choices, string(authn.FCProfile)) {
			result = flowHandleProfilePhase(t, ctx, client, flowID)
		} else if choiceIsAvailable(choices, string(authn.FCAgreements)) {
			result = flowHandleAgreementsPhase(t, ctx, client, flowID, result)
		} else {
			fmt.Printf("Phase %s not handled", result.FlowPhase)
		}
		flowPhase = result.FlowPhase
		choices = result.FlowChoices
	}

	fcresp, err := client.Flow.Complete(ctx,
		authn.FlowCompleteRequest{
			FlowID: flowID,
		})
	assert.NoError(t, err)
	assert.NotNil(t, fcresp)
	assert.NotNil(t, fcresp.Result)
	return fcresp.Result
}

func Login(t *testing.T, email, password string) *authn.FlowCompleteResult {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)
	fsresp, err := client.Flow.Start(ctx,
		authn.FlowStartRequest{
			Email:     EMAIL_TEST,
			FlowTypes: []authn.FlowType{authn.FTsignin},
			CBURI:     CB_URI,
		})
	assert.NoError(t, err)
	assert.NotNil(t, fsresp)
	assert.NotNil(t, fsresp.Result)
	flowID := fsresp.Result.FlowID

	furesp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flowID,
		Choice: authn.FCPassword,
		Data: authn.FlowUpdateDataPassword{
			Password: password,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, furesp)
	assert.NotNil(t, furesp.Result)

	fcresp, err := client.Flow.Complete(ctx,
		authn.FlowCompleteRequest{
			FlowID: flowID,
		})
	assert.NoError(t, err)
	return fcresp.Result
}

func Test_Integration_User_Create(t *testing.T) {
	CreateAndLogin(t, EMAIL_TEST, PASSWORD_OLD)
	CreateAndLogin(t, EMAIL_DELETE, PASSWORD_OLD)
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
	result := Login(t, EMAIL_TEST, PASSWORD_OLD)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ActiveToken)
	assert.NotEmpty(t, result.RefreshToken)

	// Change password
	input2 := authn.ClientPasswordChangeRequest{
		Token:       result.ActiveToken.Token,
		OldPassword: PASSWORD_OLD,
		NewPassword: PASSWORD_NEW,
	}
	resp2, err := client.Client.Password.Change(ctx, input2)
	assert.NoError(t, err)
	assert.Empty(t, resp2.Result)

	// Check token
	input6 := authn.ClientTokenCheckRequest{
		Token: result.ActiveToken.Token,
	}
	resp6, err := client.Client.Token.Check(ctx, input6)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp6.Result)

	// Refresh token
	input4 := authn.ClientSessionRefreshRequest{
		UserToken:    result.ActiveToken.Token,
		RefreshToken: result.RefreshToken.Token,
	}
	resp4, err := client.Client.Session.Refresh(ctx, input4)
	assert.NoError(t, err)
	assert.NotNil(t, resp4.Result)
	assert.NotEmpty(t, resp4.Result.ActiveToken)
	assert.NotEmpty(t, resp4.Result.RefreshToken)

	// Client session logout
	input7 := authn.ClientSessionLogoutRequest{
		Token: resp4.Result.ActiveToken.Token,
	}
	resp7, err := client.Client.Session.Logout(ctx, input7)
	assert.NoError(t, err)
	assert.NotNil(t, resp7)

	// Re-Login with password
	result = Login(t, EMAIL_TEST, PASSWORD_NEW)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ActiveToken)
	assert.NotEmpty(t, result.RefreshToken)

	// Get profile by email
	input := authn.UserProfileGetRequest{
		Email: pangea.String(EMAIL_TEST),
	}
	resp, err := client.User.Profile.Get(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	USER_ID = resp.Result.ID
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Equal(t, resp.Result.Profile, authn.ProfileData(PROFILE_OLD))

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

	// Get profile by user_id
	input := authn.UserProfileGetRequest{
		ID: pangea.String(USER_ID),
	}
	resp, err := client.User.Profile.Get(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_ID, resp.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Equal(t, resp.Result.Profile, authn.ProfileData(PROFILE_OLD))

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
		Email:    pangea.String(EMAIL_TEST),
		Disabled: pangea.Bool(false),
	}
	resp, err := client.User.Update(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, USER_ID, resp.Result.ID)
	assert.Equal(t, EMAIL_TEST, resp.Result.Email)
	assert.Equal(t, false, resp.Result.Disabled)
}

func Test_Integration_Client_Session(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	// Login with password
	result := Login(t, EMAIL_TEST, PASSWORD_NEW)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ActiveToken)
	assert.NotEmpty(t, result.RefreshToken)

	filter := authn.NewFilterSessionList()

	// Client session list
	input2 := authn.ClientSessionListRequest{
		Token:  result.ActiveToken.Token,
		Filter: filter.Filter(),
	}
	resp2, err := client.Client.Session.List(ctx, input2)
	assert.NoError(t, err)
	assert.Greater(t, len(resp2.Result.Sessions), 0)

	for _, s := range resp2.Result.Sessions {
		input3 := authn.ClientSessionInvalidateRequest{
			Token:     result.ActiveToken.Token,
			SessionID: s.ID,
		}
		client.Client.Session.Invalidate(ctx, input3)
	}
}

func Test_Integration_Session(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := authnIntegrationCfg(t)
	client := authn.New(cfg)

	// Login with password
	result := Login(t, EMAIL_TEST, PASSWORD_NEW)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ActiveToken)
	assert.NotEmpty(t, result.RefreshToken)

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
