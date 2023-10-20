package authn

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	v "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	di "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/domain_intel"
	ipi "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
)

type ClientUserinfoResult struct {
	RefreshToken LoginToken  `json:"refresh_token"`
	ActiveToken  *LoginToken `json:"active_token,omitempty"`
}

type ClientUserinfoRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Code string `json:"code"`
}

// @summary Get User (client token)
//
// @description Retrieve the logged in user's token and information.
//
// @operationId authn_post_v2_client_userinfo
//
// @example
//
//	input := authn.ClientUserinfoRequest{
//		Code: "pmc_d6chl6qulpn3it34oerwm3cqwsjd6dxw",
//	}
//
//	esp, err := authncli.Client.Userinfo(ctx, input)
func (a *Client) Userinfo(ctx context.Context, input ClientUserinfoRequest) (*pangea.PangeaResponse[ClientUserinfoResult], error) {
	return request.DoPost(ctx, a.client, "v2/client/userinfo", &input, &ClientUserinfoResult{})
}

type ClientJWKSResult struct {
	Keys []v.JWT `json:"keys"`
}

// @summary Get JWT verification keys
//
// @description Get JWT verification keys.
//
// @operationId authn_post_v2_client_jwks
//
// @example
//
//	resp, err := authncli.Client.JWKS(ctx)
func (a *Client) JWKS(ctx context.Context) (*pangea.PangeaResponse[ClientJWKSResult], error) {
	return request.DoPost(ctx, a.client, "v2/client/jwks", &pangea.BaseRequest{}, &ClientJWKSResult{})
}

type ClientTokenCheckRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Token string `json:"token"`
}

type ClientTokenCheckResult struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Life      int         `json:"life"`
	Expire    string      `json:"expire"`
	Identity  string      `json:"identity"`
	Email     string      `json:"email"`
	Scopes    Scopes      `json:"scopes"`
	Profile   ProfileData `json:"profile"`
	CreatedAt string      `json:"created_at"`
}

// @summary Check a token
//
// @description Look up a token and return its contents.
//
// @operationId authn_post_v2_client_token_check
//
// @example
//
//	input := authn.ClientTokenCheckRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authcli.Client.Token.Check(ctx, input)
func (a *ClientToken) Check(ctx context.Context, input ClientTokenCheckRequest) (*pangea.PangeaResponse[ClientTokenCheckResult], error) {
	return request.DoPost(ctx, a.Client, "v2/client/token/check", &input, &ClientTokenCheckResult{})
}

type ClientPasswordChangeRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Token       string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ClientPasswordChangeResult struct {
}

// @summary Change a user's password
//
// @description Change a user's password given the current password.
//
// @operationId authn_post_v2_client_password_change
//
// @example
//
//	input := authn.ClientPasswordChangeRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//		OldPassword: "hunter2",
//		NewPassword: "My2n+Password",
//	}
//
//	resp, err := authncli.Client.Password.Change(ctx, input)
func (a *ClientPassword) Change(ctx context.Context, input ClientPasswordChangeRequest) (*pangea.PangeaResponse[ClientPasswordChangeResult], error) {
	return request.DoPost(ctx, a.Client, "v2/client/password/change", &input, &ClientPasswordChangeResult{})
}

type IDProvider string

const (
	IDPFacebook        IDProvider = "facebook"
	IDPGithub                     = "github"
	IDPGoogle                     = "google"
	IDPMicrosoftOnline            = "microsoftonline"
	IDPPassword                   = "password"
)

type MFAProvider string

const (
	MFAPTOTP     MFAProvider = "totp"
	MFAPEmailOTP             = "email_otp"
	IDPSMSOTP                = "sms_otp"
)

type FlowType string

const (
	FTsignin FlowType = "signin"
	FTsignup          = "signup"
)

type ProfileData map[string]string
type Scopes []string

type UserCreateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Email   string      `json:"email"`
	Profile ProfileData `json:"profile"`
}

type UserCreateResult struct {
	User
}

// @summary Create User
//
// @description Create a user. Also allows creating the user's credentials.
//
// @operationId authn_post_v2_user_create
//
// @example
//
//	profile := &authn.ProfileData{
//		"first_name": "Joe",
//		"last_name": "User",
//	}
//
//	input := authn.UserCreateRequest{
//		Email: "joe.user@email.com",
//		Profile: profile,
//	}
//
//	resp, err := authncli.User.Create(ctx, input)
func (a *User) Create(ctx context.Context, input UserCreateRequest) (*pangea.PangeaResponse[UserCreateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/create", &input, &UserCreateResult{})
}

type UserDeleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Email string `json:"email,omitempty"`
	ID    string `json:"id,omitempty"`
}

type UserDeleteResult struct {
}

// @summary Delete User
//
// @description Delete a user.
//
// @operationId authn_post_v2_user_delete
//
// @example
//
//	input := UserDeleteRequest{
//		Email: "joe.user@email.com",
//	}
//
//	authncli.User.Delete(ctx, input)
func (a *User) Delete(ctx context.Context, input UserDeleteRequest) (*pangea.PangeaResponse[UserDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/delete", &input, &UserDeleteResult{})
}

type UserItem struct {
	ID                      string          `json:"id"`
	Email                   string          `json:"email"`
	Profile                 ProfileData     `json:"profile"`
	Verified                bool            `json:"verified"`
	Disabled                bool            `json:"disabled"`
	AcceptedEulaID          *string         `json:"accepted_eula_id,omitempty"`
	AcceptedPrivacyPolicyID *string         `json:"accepted_privacy_policy_id,omitempty"`
	LastLoginAt             *string         `json:"last_login_at,omitempty"`
	CreatedAt               string          `json:"created_at"`
	LoginCount              int             `json:"login_count"`
	LastLoginIP             *string         `json:"last_login_ip,omitempty"`
	LastLoginCity           *string         `json:"last_login_city,omitempty"`
	LastLoginCountry        *string         `json:"last_login_country,omitempty"`
	Authenticators          []Authenticator `json:"authenticators"`
}

type UserUpdateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID       *string `json:"id,omitempty"`
	Email    *string `json:"email,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
}

type UserUpdateResult struct {
	UserItem
}

// @summary Update user's settings
//
// @description Update user's settings.
//
// @operationId authn_post_v2_user_update
//
// @example
//
//	input := authn.UserUpdateRequest{
//		Email: pangea.String("joe.user@email.com"),
//		RequireMFA: pangea.Bool(true),
//	}
//
//	resp, err := authncli.User.Update(ctx, input)
func (a *User) Update(ctx context.Context, input UserUpdateRequest) (*pangea.PangeaResponse[UserUpdateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/update", &input, &UserUpdateResult{})
}

type UserInviteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Inviter  string `json:"inviter"`
	Email    string `json:"email"`
	Callback string `json:"callback"`
	State    string `json:"state"`
}

type UserInviteResult struct {
	ID         string `json:"id"`
	Inviter    string `json:"inviter"`
	InviteOrg  string `json:"invite_org"`
	Email      string `json:"email"`
	Callback   string `json:"callback"`
	State      string `json:"state"`
	RequireMFA bool   `json:"require_mfa"`
	CreatedAt  string `json:"created_at"`
	Expire     string `json:"expire"`
}

// @summary Invite User
//
// @description Send an invitation to a user.
//
// @operationId authn_post_v2_user_invite
//
// @example
//
//	input := authn.UserInviteRequest{
//		Inviter: "admin@email.com",
//		Email: "joe.user@email.com",
//		Callback: "https://www.myserver.com/callback",
//		State: "pcb_zurr3lkcwdp5keq73htsfpcii5k4zgm7",
//	}
//
//	resp, err := authncli.User.Invite(ctx, input)
func (a *User) Invite(ctx context.Context, input UserInviteRequest) (*pangea.PangeaResponse[UserInviteResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/invite", &input, &UserInviteResult{})
}

type UserListOrderBy string

const (
	ULOBid          UserListOrderBy = "id"
	ULOBcreatedAt                   = "created_at"
	ULOBlastLoginAt                 = "last_login_at"
	ULOBemail                       = "email"
)

type UserInviteListOrderBy string

const (
	UILOBid        UserListOrderBy = "id"
	UILOBcreatedAt                 = "created_at"
	UILOBtype                      = "type"
	UILOBexpire                    = "expire"
	UILOBcallback                  = "callback"
	UILOBstate                     = "state"
	UILOBemail                     = "email"
	UILOBinviter                   = "inviter"
	UILOBinviteOrg                 = "invite_org"
)

type FilterUserList struct {
	pangea.FilterBase
	acceptedEulaID   *pangea.FilterMatch[string]
	createdAt        *pangea.FilterRange[string]
	disabled         *pangea.FilterEqual[bool]
	email            *pangea.FilterMatch[string]
	id               *pangea.FilterMatch[string]
	lastLoginAt      *pangea.FilterMatch[string]
	lastLoginIP      *pangea.FilterMatch[string]
	lastLoginCity    *pangea.FilterMatch[string]
	lastLoginCountry *pangea.FilterMatch[string]
	loginCount       *pangea.FilterRange[int]
	requireMFA       *pangea.FilterEqual[bool]
	scopes           *pangea.FilterEqual[[]string]
	verified         *pangea.FilterEqual[bool]
}

func NewFilterUserList() *FilterUserList {
	filter := make(pangea.Filter)
	return &FilterUserList{
		FilterBase:       *pangea.NewFilterBase(filter),
		acceptedEulaID:   pangea.NewFilterMatch[string]("accepted_eula_id", &filter),
		createdAt:        pangea.NewFilterRange[string]("created_at", &filter),
		disabled:         pangea.NewFilterEqual[bool]("diabled", &filter),
		email:            pangea.NewFilterMatch[string]("email", &filter),
		id:               pangea.NewFilterMatch[string]("id", &filter),
		lastLoginAt:      pangea.NewFilterMatch[string]("last_login_at", &filter),
		lastLoginIP:      pangea.NewFilterMatch[string]("last_login_ip", &filter),
		lastLoginCity:    pangea.NewFilterMatch[string]("last_login_city", &filter),
		lastLoginCountry: pangea.NewFilterMatch[string]("last_login_country", &filter),
		loginCount:       pangea.NewFilterRange[int]("login_count", &filter),
		requireMFA:       pangea.NewFilterEqual[bool]("require_mfa", &filter),
		scopes:           pangea.NewFilterEqual[[]string]("scopes", &filter),
		verified:         pangea.NewFilterEqual[bool]("verified", &filter),
	}
}

func (fu *FilterUserList) AcceptedEulaID() *pangea.FilterMatch[string] {
	return fu.acceptedEulaID
}

func (fu *FilterUserList) CreatedAt() *pangea.FilterRange[string] {
	return fu.createdAt
}

func (fu *FilterUserList) Disabled() *pangea.FilterEqual[bool] {
	return fu.disabled
}

func (fu *FilterUserList) Email() *pangea.FilterMatch[string] {
	return fu.email
}

func (fu *FilterUserList) ID() *pangea.FilterMatch[string] {
	return fu.id
}

func (fu *FilterUserList) LastLoginAt() *pangea.FilterMatch[string] {
	return fu.lastLoginAt
}

func (fu *FilterUserList) LastLoginIP() *pangea.FilterMatch[string] {
	return fu.lastLoginIP
}

func (fu *FilterUserList) LastLoginCity() *pangea.FilterMatch[string] {
	return fu.lastLoginCity
}

func (fu *FilterUserList) LastLoginCountry() *pangea.FilterMatch[string] {
	return fu.lastLoginCountry
}

func (fu *FilterUserList) LoginCount() *pangea.FilterRange[int] {
	return fu.loginCount
}

func (fu *FilterUserList) RequireMFA() *pangea.FilterEqual[bool] {
	return fu.requireMFA
}

func (fu *FilterUserList) Scopes() *pangea.FilterEqual[[]string] {
	return fu.scopes
}

func (fu *FilterUserList) Verified() *pangea.FilterEqual[bool] {
	return fu.verified
}

type UserListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Should user FilterUserList object here
	Filter  pangea.Filter   `json:"filter,omitempty"`
	Last    string          `json:"last,omitempty"`
	Order   ItemOrder       `json:"order,omitempty"`
	OrderBy UserListOrderBy `json:"order_by,omitempty"`
	Size    int             `json:"size,omitempty"`
}

type UserListResult struct {
	Users []UserItem `json:"users"`
	Last  string     `json:"last,omitempty"`
	Count int        `json:"count"`
}

// @summary List Users
//
// @description Look up users by scopes.
//
// @operationId authn_post_v2_user_list
//
// @example
//
//	input := authn.UserListRequest{}
//	resp, err := authncli.User.List(ctx, input)
func (a *User) List(ctx context.Context, input UserListRequest) (*pangea.PangeaResponse[UserListResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/list", &input, &UserListResult{})
}

type IPIntelligence struct {
	IsBad       bool               `json:"is_bad"`
	IsVPN       bool               `json:"is_vpn"`
	IsProxy     bool               `json:"is_proxy"`
	Reputation  ipi.ReputationData `json:"reputation"`
	Geolocation ipi.GeolocateData  `json:"geolocation"`
}

type DomainIntelligence struct {
	IsBad      bool              `json:"is_bad"`
	Reputation di.ReputationData `json:"reputation"`
}

type Intelligence struct {
	Embargo     bool               `json:"embargo"`
	IPIntel     IPIntelligence     `json:"ip_intel"`
	DomainIntel DomainIntelligence `json:"domain_intel"`
	UserIntel   bool               `json:"user_intel"`
}

type LoginToken struct {
	SessionToken
	Token string `json:"token"`
}

type UserProfileGetRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID    *string `json:"id,omitempty"`
	Email *string `json:"email,omitempty"`
}

type UserProfileGetResult struct {
	UserItem
}

// @summary Get user
//
// @description Get user's information by identity or email.
//
// @operationId authn_post_v2_user_profile_get
//
// @example
//
//	input := authn.UserProfileGetRequest{
//		Email: pangea.String("joe.user@email.com"),
//	}
//
//	resp, err := authncli.User.Profile.Get(ctx, input)
func (a *UserProfile) Get(ctx context.Context, input UserProfileGetRequest) (*pangea.PangeaResponse[UserProfileGetResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/profile/get", &input, &UserProfileGetResult{})
}

type UserProfileUpdateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Profile ProfileData `json:"profile"`
	ID      *string     `json:"id,omitempty"`
	Email   *string     `json:"email,omitempty"`
}

type UserProfileUpdateResult struct {
	UserItem
}

// @summary Update user
//
// @description Update user's information by identity or email.
//
// @operationId authn_post_v2_user_profile_update
//
// @example
//
//	input := authn.UserProfileUpdateRequest{
//		Email: pangea.String("joe.user@email.com"),
//		Profile: authn.ProfileData{
//			"country": "Argentina",
//		},
//	}
//
//	resp, err := authncli.User.Profile.Update(ctx, input)
func (a *UserProfile) Update(ctx context.Context, input UserProfileUpdateRequest) (*pangea.PangeaResponse[UserProfileUpdateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/profile/update", &input, &UserProfileUpdateResult{})
}

type UserInviteData struct {
	ID         string `json:"id"`
	Inviter    string `json:"inviter"`
	InviteOrg  string `json:"invite_org"`
	Email      string `json:"email"`
	Callback   string `json:"callback"`
	State      string `json:"state"`
	RequireMFA bool   `json:"require_mfa"`
	CreatedAt  string `json:"created_at"`
	Expire     string `json:"expire"`
}

// FilterUserInviteList represents the filter criteria for user invites.
type FilterUserInviteList struct {
	pangea.FilterBase
	callback   *pangea.FilterMatch[string]
	email      *pangea.FilterMatch[string]
	id         *pangea.FilterMatch[string]
	inviteOrg  *pangea.FilterMatch[string]
	inviter    *pangea.FilterMatch[string]
	state      *pangea.FilterMatch[string]
	signup     *pangea.FilterEqual[bool]
	requireMFA *pangea.FilterEqual[bool]
	expire     *pangea.FilterRange[string]
	createdAt  *pangea.FilterRange[string]
}

func NewFilterUserInviteList() *FilterUserInviteList {
	filter := make(pangea.Filter)
	return &FilterUserInviteList{
		FilterBase: *pangea.NewFilterBase(filter),
		callback:   pangea.NewFilterMatch[string]("callback", &filter),
		email:      pangea.NewFilterMatch[string]("email", &filter),
		id:         pangea.NewFilterMatch[string]("id", &filter),
		inviteOrg:  pangea.NewFilterMatch[string]("invite_org", &filter),
		inviter:    pangea.NewFilterMatch[string]("inviter", &filter),
		state:      pangea.NewFilterMatch[string]("state", &filter),
		signup:     pangea.NewFilterEqual[bool]("signup", &filter),
		requireMFA: pangea.NewFilterEqual[bool]("require_mfa", &filter),
		expire:     pangea.NewFilterRange[string]("expire", &filter),
		createdAt:  pangea.NewFilterRange[string]("created_at", &filter),
	}
}

func (f *FilterUserInviteList) Callback() *pangea.FilterMatch[string] {
	return f.callback
}

func (f *FilterUserInviteList) Email() *pangea.FilterMatch[string] {
	return f.email
}

func (f *FilterUserInviteList) ID() *pangea.FilterMatch[string] {
	return f.id
}

func (f *FilterUserInviteList) InviteOrg() *pangea.FilterMatch[string] {
	return f.inviteOrg
}

func (f *FilterUserInviteList) Inviter() *pangea.FilterMatch[string] {
	return f.inviter
}

func (f *FilterUserInviteList) State() *pangea.FilterMatch[string] {
	return f.state
}

func (f *FilterUserInviteList) Signup() *pangea.FilterEqual[bool] {
	return f.signup
}

func (f *FilterUserInviteList) RequireMFA() *pangea.FilterEqual[bool] {
	return f.requireMFA
}

func (f *FilterUserInviteList) Expire() *pangea.FilterRange[string] {
	return f.expire
}

func (f *FilterUserInviteList) CreatedAt() *pangea.FilterRange[string] {
	return f.createdAt
}

type UserInviteListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Should use FilterUserInviteList object here
	Filter  pangea.Filter         `json:"filter,omitempty"`
	Last    string                `json:"last,omitempty"`
	Order   ItemOrder             `json:"order,omitempty"`
	OrderBy UserInviteListOrderBy `json:"order_by,omitempty"`
	Size    int                   `json:"size,omitempty"`
}

type UserInviteListResult struct {
	Invites []UserInviteData `json:"invites"`
	Last    string           `json:"last"`
	Count   int              `json:"count"`
}

// @summary List Invites
//
// @description Look up active invites for the userpool.
//
// @operationId authn_post_v2_user_invite_list
//
// @example
//
//	input := authn.UserInviteListRequest{}
//	resp, err := authncli.User.Invite.List(ctx, input)
func (a *UserInvite) List(ctx context.Context, input UserInviteListRequest) (*pangea.PangeaResponse[UserInviteListResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/invite/list", &input, &UserInviteListResult{})
}

type UserInviteDeleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID string `json:"id"`
}

type UserInviteDeleteResult struct {
}

// @summary Delete Invite
//
// @description Delete a user invitation.
//
// @operationId authn_post_v2_user_invite_delete
//
// @example
//
//	input := authn.UserInviteDeleteRequest{
//		ID: "pmc_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.User.Invite.Delete(ctx, input)
func (a *UserInvite) Delete(ctx context.Context, input UserInviteDeleteRequest) (*pangea.PangeaResponse[UserInviteDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/invite/delete", &input, &UserInviteDeleteResult{})
}

type FlowCompleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID string `json:"flow_id"`
}

type FlowCompleteResult struct {
	RefreshToken LoginToken `json:"refresh_token"`
	ActiveToken  LoginToken `json:"active_token"`
}

// @summary Complete sign-up/sign-in
//
// @description Complete a sign-up or sign-in flow.
//
// @operationId authn_post_v2_flow_complete
//
// @example
//
//	input := authn.FlowCompleteRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//	}
//
//	resp, err := authncli.Flow.Complete(ctx, input)
func (a *Flow) Complete(ctx context.Context, input FlowCompleteRequest) (*pangea.PangeaResponse[FlowCompleteResult], error) {
	return request.DoPost(ctx, a.Client, "v2/flow/complete", &input, &FlowCompleteResult{})
}

type FlowChoiceItem struct {
	Choice string         `json:"choice"`
	Data   map[string]any `json:"data"`
}

type CommonFlowResult struct {
	FlowID      string           `json:"flow_id"`
	FlowType    []string         `json:"flow_type"`
	Email       string           `json:"email"`
	Disclaimer  string           `json:"disclaimer"`
	FlowPhase   string           `json:"flow_phase"`
	FlowChoices []FlowChoiceItem `json:"flow_choices"`
}

type FlowStartRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	CBURI     string      `json:"cb_uri,omitempty"`
	Email     string      `json:"email,omitempty"`
	FlowTypes []FlowType  `json:"flow_types,omitempty"`
	Provider  *IDProvider `json:"provider,omitempty"`
}

type FlowStartResult struct {
	CommonFlowResult
}

// @summary Start a sign-up/sign-in flow
//
// @description Start a new sign-up or sign-in flow.
//
// @operationId authn_post_v2_flow_start
//
// @example
//
//	fts := []FlowType{authn.FTsignin,authn.FTsignup}
//	input := authn.FlowStartRequest{
//		CBURI: "https://www.myserver.com/callback",
//		Email: "joe.user@email.com",
//		FlowTypes: fts,
//		Provider: &authn.IDPPassword,
//	}
//
//	resp, err := authncli.Flow.Start(ctx, input)
func (a *Flow) Start(ctx context.Context, input FlowStartRequest) (*pangea.PangeaResponse[FlowStartResult], error) {
	return request.DoPost(ctx, a.Client, "v2/flow/start", &input, &FlowStartResult{})
}

type FlowUpdaterData interface {
	IsFlowUpdaterData() bool
}

type FlowUpdateData struct{}

func (fud FlowUpdateData) IsFlowUpdaterData() bool {
	return true
}

type FlowUpdateDataAgreements struct {
	FlowUpdateData
	Agreed []string `json:"agreed"`
}

type FlowUpdateDataCaptcha struct {
	FlowUpdateData
	Code string `json:"code"`
}

type FlowUpdateDataEmailOTP struct {
	FlowUpdateData
	Code string `json:"code"`
}

type FlowUpdateDataMagiclink struct {
	FlowUpdateData
	State string `json:"state"`
	Code  string `json:"code"`
}

type FlowUpdateDataPassword struct {
	FlowUpdateData
	Password string `json:"password"`
}

type FlowUpdateDataProfile struct {
	FlowUpdateData
	Profile ProfileData `json:"profile"`
}

type FlowUpdateDataProvisionalEnrollment struct {
	FlowUpdateData
	State string `json:"state"`
	Code  string `json:"code"`
}

type FlowUpdateDataResetPassword struct {
	FlowUpdateData
	State string `json:"state"`
	Code  string `json:"code"`
}

type FlowUpdateDataSetEmail struct {
	FlowUpdateData
	Email string `json:"email"`
}

type FlowUpdateDataSetPassword struct {
	FlowUpdateData
	Password string `json:"password"`
}

type FlowUpdateDataSMSOTP struct {
	FlowUpdateData
	Code string `json:"code"`
}

type FlowUpdateDataSocialProvider struct {
	FlowUpdateData
	SocialProvider string `json:"social_provider"`
	URI            string `json:"uri"`
}

type FlowUpdateDataTOTP struct {
	FlowUpdateData
	Code string `json:"code"`
}

type FlowUpdateDataVerifyEmail struct {
	FlowUpdateData
	State string `json:"state"`
	Code  string `json:"code"`
}

type FlowChoice string

const (
	FCAgreements            FlowChoice = "agreements"
	FCCaptcha                          = "captcha"
	FCEmailOTP                         = "email_otp"
	FCMagiclink                        = "magiclink"
	FCPassword                         = "password"
	FCProfile                          = "profile"
	FCProvisionalEnrollment            = "provisional_enrollment"
	FCResetPassword                    = "reset_password"
	FCSetEmail                         = "set_mail"
	FCSetPassword                      = "set_password"
	FCSMSOTP                           = "sms_otp"
	FCSocial                           = "social"
	FCTOTP                             = "totp"
	FCVerifyEmail                      = "verify_email"
)

type FlowUpdateRequest struct {
	pangea.BaseRequest
	FlowID string          `json:"flow_id"`
	Choice FlowChoice      `json:"choice"`
	Data   FlowUpdaterData `json:"data"`
}

type FlowUpdateResult struct {
	CommonFlowResult
}

// @summary Update a sign-up/sign-in flow
//
// @description Update a sign-up/sign-in flow.
//
// @operationId authn_post_v2_flow_update
//
// @example
//
//	input := authn.FlowUpdateRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Choice: authn.FCPassword,
//		Data: authn.FlowUpdateDataPassword{
//			Password: "somenewpassword",
//		}
//	}
//
//	resp, err := authncli.Flow.Update(ctx, input)
func (a *Flow) Update(ctx context.Context, input FlowUpdateRequest) (*pangea.PangeaResponse[FlowUpdateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/flow/update", &input, &FlowUpdateResult{})
}

type FlowRestartData struct{}

type FlowRestartDataSMSOTP struct {
	FlowRestartData
	Phone string `json:"phone"`
}

type FlowRestartRequest struct {
	pangea.BaseRequest
	FlowID string          `json:"flow_id"`
	Choice FlowChoice      `json:"choice"`
	Data   FlowRestartData `json:"data"`
}

type FlowRestartResult struct {
	CommonFlowResult
}

// @summary Restart a sign-up/sign-in flow
//
// @description Restart a signup-up/in flow choice.
//
// @operationId authn_post_v2_flow_restart
//
// @example
//
//	input := authn.FlowRestartRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Choice: authn.FCPassword,
//		Data: authn.FlowRestartData{},
//	}
//
//	resp, err := authncli.Flow.Restart(ctx, input)
func (a *Flow) Restart(ctx context.Context, input FlowRestartRequest) (*pangea.PangeaResponse[FlowRestartResult], error) {
	return request.DoPost(ctx, a.Client, "v2/flow/restart", &input, &FlowRestartResult{})
}

type UserAuthenticatorsDeleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID              *string `json:"id,omitempty"`
	Email           *string `json:"email,omitempty"`
	AuthenticatorID string  `json:"authenticator_id"`
}

type UserAuthenticatorsDeleteResult struct {
}

// @summary Delete user authenticator
//
// @description Delete a user's authenticator.
//
// @operationId authn_post_v2_user_authenticators_delete
//
// @example
//
//	input := authn.UserAuthenticatorsDeleteRequest{
//		ID: pangea.String("pui_zgp532cx6opljeavvllmbi3iwmq72f7f"),
//		AuthenticatorID: "pau_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.User.Authenticators.Delete(ctx, input)
func (a *UserAuthenticators) Delete(ctx context.Context, input UserAuthenticatorsDeleteRequest) (*pangea.PangeaResponse[UserAuthenticatorsDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/authenticators/delete", &input, &UserAuthenticatorsDeleteResult{})
}

type UserAuthenticatorsListRequest struct {
	pangea.BaseRequest

	ID    *string `json:"id,omitempty"`
	Email *string `json:"email,omitempty"`
}

type Authenticator struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Enable   bool    `json:"enable"`
	Provider *string `json:"provider,omitempty"`
	RPID     *string `json:"rpid,omitempty"`
	Phase    *string `json:"phase,omitempty"`
}

type UserAuthenticatorsListResult struct {
	Authenticators []Authenticator `json:"authenticators"`
}

// @summary Get user authenticators
//
// @description Get user's authenticators by identity or email.
//
// @operationId authn_post_v2_user_authenticators_list
//
// @example
//
//	input := authn.UserAuthenticatorsListRequest{
//		ID: pangea.String("pui_xpkhwpnz2cmegsws737xbsqnmnuwtvm5"),
//	}
//
//	resp, err := authncli.User.Authenticators.List(ctx, input)
func (a *UserAuthenticators) List(ctx context.Context, input UserAuthenticatorsListRequest) (*pangea.PangeaResponse[UserAuthenticatorsListResult], error) {
	return request.DoPost(ctx, a.Client, "v2/user/authenticators/list", &input, &UserAuthenticatorsListResult{})
}

type ClientSessionInvalidateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Token     string `json:"token"`
	SessionID string `json:"session_id"`
}

type ClientSessionInvalidateResult struct {
}

// @summary Invalidate Session | Client
//
// @description Invalidate a session by session ID using a client token.
//
// @operationId authn_post_v2_client_session_invalidate
//
// @example
//
//	input := authn.ClientSessionInvalidateRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//		SessionID: "pmt_zppkzrjguxyblaia6itbiesejn7jejnr",
//	}
//
//	resp, err := authncli.Client.Session.Invalidate(ctx, input)
func (a *ClientSession) Invalidate(ctx context.Context, input ClientSessionInvalidateRequest) (*pangea.PangeaResponse[ClientSessionInvalidateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/client/session/invalidate", &input, &ClientSessionInvalidateResult{})
}

type SessionListOrderBy string

const (
	SLOBid            SessionListOrderBy = "id"
	SLOBcreatedAt                        = "created_at"
	SLOBtype                             = "type"
	SLOBidentity                         = "identity"
	SLOBemail                            = "email"
	SLOBexpire                           = "expire"
	SLOBactiveTokenID                    = "active_token_id"
)

type ItemOrder string

const (
	IOasc  ItemOrder = "asc"
	IOdesc           = "desc"
)

type FilterSessionList struct {
	pangea.FilterBase
	id        *pangea.FilterMatch[string]
	typeStr   *pangea.FilterMatch[string]
	identity  *pangea.FilterMatch[string]
	email     *pangea.FilterMatch[string]
	createdAt *pangea.FilterRange[string]
	expire    *pangea.FilterRange[string]
	scopes    *pangea.FilterEqual[[]string]
}

func NewFilterSessionList() *FilterSessionList {
	filter := make(pangea.Filter)
	return &FilterSessionList{
		FilterBase: *pangea.NewFilterBase(filter),
		id:         pangea.NewFilterMatch[string]("id", &filter),
		typeStr:    pangea.NewFilterMatch[string]("type", &filter),
		identity:   pangea.NewFilterMatch[string]("identity", &filter),
		email:      pangea.NewFilterMatch[string]("email", &filter),
		createdAt:  pangea.NewFilterRange[string]("created_at", &filter),
		expire:     pangea.NewFilterRange[string]("expire", &filter),
		scopes:     pangea.NewFilterEqual[[]string]("scopes", &filter),
	}
}

func (f *FilterSessionList) ID() *pangea.FilterMatch[string] {
	return f.id
}

func (f *FilterSessionList) Type() *pangea.FilterMatch[string] {
	return f.typeStr
}

func (f *FilterSessionList) Identity() *pangea.FilterMatch[string] {
	return f.identity
}

func (f *FilterSessionList) Email() *pangea.FilterMatch[string] {
	return f.email
}

func (f *FilterSessionList) CreatedAt() *pangea.FilterRange[string] {
	return f.createdAt
}

func (f *FilterSessionList) Expire() *pangea.FilterRange[string] {
	return f.expire
}

func (f *FilterSessionList) Scopes() *pangea.FilterEqual[[]string] {
	return f.scopes
}

type ClientSessionListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Token string `json:"token"`

	// Should use FilterSessionList object here
	Filter  pangea.Filter      `json:"filter,omitempty"`
	Last    string             `json:"last,omitempty"`
	Order   ItemOrder          `json:"order,omitempty"`
	OrderBy SessionListOrderBy `json:"order_by,omitempty"`
}

type SessionToken struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Life         int           `json:"life"`
	Expire       string        `json:"expire"`
	Email        string        `json:"email"`
	Scopes       Scopes        `json:"scopes"`
	Profile      ProfileData   `json:"profile"`
	CreatedAt    string        `json:"created_at"`
	Intelligence *Intelligence `json:"intelligence,omitempty"`
}

type SessionItem struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Life        int           `json:"life"`
	Expire      string        `json:"expire"`
	Identity    string        `json:"identity"`
	Email       string        `json:"email"`
	Scopes      Scopes        `json:"scopes"`
	Profile     ProfileData   `json:"profile"`
	CreatedAt   string        `json:"created_at"`
	ActiveToken *SessionToken `json:"active_token,omitempty"`
}

type SessionListResult struct {
	Sessions []SessionItem `json:"sessions"`
	Last     string        `json:"last"`
}

// @summary List sessions (client token)
//
// @description List sessions using a client token.
//
// @operationId authn_post_v2_client_session_list
//
// @example
//
//	input := authn.ClientSessionListRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.Client.Session.List(ctx, input)
func (a *ClientSession) List(ctx context.Context, input ClientSessionListRequest) (*pangea.PangeaResponse[SessionListResult], error) {
	return request.DoPost(ctx, a.Client, "v2/client/session/list", &input, &SessionListResult{})
}

type ClientSessionLogoutRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Token string `json:"token"`
}

type ClientSessionLogoutResult struct {
}

// @summary Log out (client token)
//
// @description Log out the current user's session.
//
// @operationId authn_post_v2_client_session_logout
//
// @example
//
//	input := authn.ClientSessionLogoutRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.Client.Session.Logout(ctx, input)
func (a *ClientSession) Logout(ctx context.Context, input ClientSessionLogoutRequest) (*pangea.PangeaResponse[ClientSessionLogoutResult], error) {
	return request.DoPost(ctx, a.Client, "v2/client/session/logout", &input, &ClientSessionLogoutResult{})
}

type ClientSessionRefreshRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	RefreshToken string `json:"refresh_token"`
	UserToken    string `json:"user_token,omitempty"`
}

type ClientSessionRefreshResult struct {
	RefreshToken LoginToken `json:"refresh_token"`
	ActiveToken  LoginToken `json:"active_token"`
}

// @summary Refresh a Session
//
// @description Refresh a session token.
//
// @operationId authn_post_v2_client_session_refresh
//
// @example
//
//	input := authn.ClientSessionRefreshRequest{
//		RefreshToken: "ptr_xpkhwpnz2cmegsws737xbsqnmnuwtbm5",
//		UserToken: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.Client.Session.Refresh(ctx, input)
func (a *ClientSession) Refresh(ctx context.Context, input ClientSessionRefreshRequest) (*pangea.PangeaResponse[ClientSessionRefreshResult], error) {
	return request.DoPost(ctx, a.Client, "v2/client/session/refresh", &input, &ClientSessionRefreshResult{})
}

type SessionListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Should use FilterSessionList object here
	Filter  pangea.Filter      `json:"filter,omitempty"`
	Last    string             `json:"last,omitempty"`
	Order   ItemOrder          `json:"order,omitempty"`
	OrderBy SessionListOrderBy `json:"order_by,omitempty"`
}

// @summary List session (service token)
//
// @description List sessions.
//
// @operationId authn_post_v2_session_list
//
// @example
//
//	input := authn.SessionListRequest{}
//	resp, err := authn.Session.List(ctx, input)
func (a *Session) List(ctx context.Context, input SessionListRequest) (*pangea.PangeaResponse[SessionListResult], error) {
	return request.DoPost(ctx, a.Client, "v2/session/list", &input, &SessionListResult{})
}

type SessionInvalidateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	SessionID string `json:"session_id"`
}

type SessionInvalidateResult struct {
}

// @summary Invalidate Session
//
// @description Invalidate a session by session ID.
//
// @operationId authn_post_v2_session_invalidate
//
// @example
//
//	input := authn.SessionInvalidateRequest{
//		SessionID: "pmt_zppkzrjguxyblaia6itbiesejn7jejnr",
//	}
//
//	resp, err := authncli.Session.Invalidate(ctx, input)
func (a *Session) Invalidate(ctx context.Context, input SessionInvalidateRequest) (*pangea.PangeaResponse[SessionInvalidateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/session/invalidate", &input, &SessionInvalidateResult{})
}

type SessionLogoutRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID string `json:"user_id"`
}

type SessionLogoutResult struct {
}

// @summary Log out (service token)
//
// @description Invalidate all sessions belonging to a user.
//
// @operationId authn_post_v2_session_logout
//
// @example
//
//	input := authn.SessionLogoutRequest{
//		UserID: "pui_xpkhwpnz2cmegsws737xbsqnmnuwtvm5",
//	}
//
//	resp, err := authncli.Session.Logout(ctx, input)
func (a *Session) Logout(ctx context.Context, input SessionLogoutRequest) (*pangea.PangeaResponse[SessionLogoutResult], error) {
	return request.DoPost(ctx, a.Client, "v2/session/logout", &input, &SessionLogoutResult{})
}

type AgreementType string

const (
	ATeula          AgreementType = "eula"
	ATprivacyPolicy               = "privacy_policy"
)

type AgreementCreateRequest struct {
	pangea.BaseRequest

	Type   AgreementType `json:"type"`
	Name   string        `json:"name"`
	Text   string        `json:"text"`
	Active *bool         `json:"active,omitempty"`
}

type AgreementInfo struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	PublishedAt string `json:"published_at,omitempty"`
	Name        string `json:"name"`
	Text        string `json:"text"`
	Active      bool   `json:"active"`
}

type AgreementCreateResult AgreementInfo

// @summary Create an agreement
//
// @description Create an agreement.
//
// @operationId authn_post_v2_agreements_create
//
// @example
//
//	input := authn.AgreementCreateRequest{
//		Type: authn.ATeula,
//		Name: "EULA_V1",
//		Text: "You agree to behave yourself while logged in.",
//	}
//
//	resp, err := authncli.Agreements.Create(ctx, input)
func (a *Agreements) Create(ctx context.Context, input AgreementCreateRequest) (*pangea.PangeaResponse[AgreementCreateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/agreements/create", &input, &AgreementCreateResult{})
}

type AgreementDeleteRequest struct {
	pangea.BaseRequest

	Type AgreementType `json:"type"`
	ID   string        `json:"id"`
}

type AgreementDeleteResult struct{}

// @summary Delete an agreement
//
// @description Delete an agreement.
//
// @operationId authn_post_v2_agreements_delete
//
// @example
//
//	input := authn.AgreementDeleteRequest{
//		Type: authn.ATeula,
//		ID: "peu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.Agreements.Delete(ctx, input)
func (a *Agreements) Delete(ctx context.Context, input AgreementDeleteRequest) (*pangea.PangeaResponse[AgreementDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v2/agreements/delete", &input, &AgreementDeleteResult{})
}

type AgreementListOrderBy string

const (
	ALOBid        AgreementListOrderBy = "id"
	ALOBcreatedAt                      = "created_at"
	ALOBname                           = "name"
	ALOBtext                           = "text"
)

type FilterAgreementList struct {
	pangea.FilterBase
	active       *pangea.FilterEqual[bool]
	created_at   *pangea.FilterRange[string]
	published_at *pangea.FilterRange[string]
	typeStr      *pangea.FilterMatch[string]
	id           *pangea.FilterMatch[string]
	name         *pangea.FilterMatch[string]
	text         *pangea.FilterMatch[string]
}

func NewFilterAgreementList() *FilterAgreementList {
	filter := make(pangea.Filter)
	return &FilterAgreementList{
		FilterBase:   *pangea.NewFilterBase(filter),
		active:       pangea.NewFilterEqual[bool]("active", &filter),
		created_at:   pangea.NewFilterRange[string]("created_at", &filter),
		published_at: pangea.NewFilterRange[string]("published_at", &filter),
		typeStr:      pangea.NewFilterMatch[string]("type", &filter),
		id:           pangea.NewFilterMatch[string]("id", &filter),
		name:         pangea.NewFilterMatch[string]("name", &filter),
		text:         pangea.NewFilterMatch[string]("text", &filter),
	}
}

func (f *FilterAgreementList) Active() *pangea.FilterEqual[bool] {
	return f.active
}

func (f *FilterAgreementList) CreatedAt() *pangea.FilterRange[string] {
	return f.created_at
}

func (f *FilterAgreementList) PublishedAt() *pangea.FilterRange[string] {
	return f.published_at
}

func (f *FilterAgreementList) Type() *pangea.FilterMatch[string] {
	return f.typeStr
}

func (f *FilterAgreementList) ID() *pangea.FilterMatch[string] {
	return f.id
}

func (f *FilterAgreementList) Name() *pangea.FilterMatch[string] {
	return f.name
}

func (f *FilterAgreementList) Text() *pangea.FilterMatch[string] {
	return f.text
}

type AgreementListRequest struct {
	pangea.BaseRequest

	// Should use FilterAgreementList object here
	Filter  map[string]any       `json:"filter,omitempty"`
	Last    string               `json:"last,omitempty"`
	Order   ItemOrder            `json:"order,omitempty"`
	OrderBy AgreementListOrderBy `json:"order_by,omitempty"`
	Size    int                  `json:"size,omitempty"`
}

type AgreementListResult struct {
	Agreements []AgreementInfo `json:"agreements"`
	Count      int             `json:"count"`
	Last       string          `json:"last,omitempty"`
}

// @summary List agreements
//
// @description List agreements.
//
// @operationId authn_post_v2_agreements_list
//
// @example
//
//	input := authn.AgreementListRequest{}
//
//	resp, err := authncli.Agreements.List(ctx, input)
func (a *Agreements) List(ctx context.Context, input AgreementListRequest) (*pangea.PangeaResponse[AgreementListResult], error) {
	return request.DoPost(ctx, a.Client, "v2/agreements/list", &input, &AgreementListResult{})
}

type AgreementUpdateRequest struct {
	pangea.BaseRequest

	Type   AgreementType `json:"type"`
	ID     string        `json:"id"`
	Name   *string       `json:"name,omitempty"`
	Text   *string       `json:"text,omitempty"`
	Active *bool         `json:"active,omitempty"`
}

type AgreementUpdateResult AgreementInfo

// @summary Update agreement
//
// @description Update agreement.
//
// @operationId authn_post_v2_agreements_update
//
// @example
//
//	input := authn.AgreementUpdateRequest{
//		Type: authn.ATeula,
//		ID: pangea.String("peu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a"),
//		Text: pangea.String("You agree to behave yourself while logged in. Don't be evil."),
//	}
//
//	resp, err := authncli.Agreements.Update(ctx, input)
func (a *Agreements) Update(ctx context.Context, input AgreementUpdateRequest) (*pangea.PangeaResponse[AgreementUpdateResult], error) {
	return request.DoPost(ctx, a.Client, "v2/agreements/update", &input, &AgreementUpdateResult{})
}
