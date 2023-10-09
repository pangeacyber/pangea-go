package authn

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	v "github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
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
// @operationId authn_post_v1_client_userinfo
//
// @example
//
//	input := authn.ClientUserinfoRequest{
//		Code: "pmc_d6chl6qulpn3it34oerwm3cqwsjd6dxw",
//	}
//
//	esp, err := authncli.Client.Userinfo(ctx, input)
func (a *Client) Userinfo(ctx context.Context, input ClientUserinfoRequest) (*pangea.PangeaResponse[ClientUserinfoResult], error) {
	return request.DoPost(ctx, a.client, "v1/client/userinfo", &input, &ClientUserinfoResult{})
}

type ClientJWKSResult struct {
	Keys []v.JWT `json:"keys"`
}

// @summary Get JWT verification keys
//
// @description Get JWT verification keys.
//
// @operationId authn_post_v1_client_jwks
//
// @example
//
//	resp, err := authncli.Client.JWKS(ctx)
func (a *Client) JWKS(ctx context.Context) (*pangea.PangeaResponse[ClientJWKSResult], error) {
	return request.DoPost(ctx, a.client, "v1/client/jwks", &pangea.BaseRequest{}, &ClientJWKSResult{})
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
// @operationId authn_post_v1_client_token_check
//
// @example
//
//	input := authn.ClientTokenCheckRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authcli.Client.Token.Check(ctx, input)
func (a *ClientToken) Check(ctx context.Context, input ClientTokenCheckRequest) (*pangea.PangeaResponse[ClientTokenCheckResult], error) {
	return request.DoPost(ctx, a.Client, "v1/client/token/check", &input, &ClientTokenCheckResult{})
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
// @operationId authn_post_v1_client_password_change
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
	return request.DoPost(ctx, a.Client, "v1/client/password/change", &input, &ClientPasswordChangeResult{})
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

	Email         string       `json:"email"`
	Authenticator string       `json:"authenticator"`
	IDProvider    IDProvider   `json:"id_provider"`
	Verified      *bool        `json:"verified,omitempty"`
	RequireMFA    *bool        `json:"require_mfa,omitempty"`
	Profile       *ProfileData `json:"profile,omitempty"`
	Scopes        *Scopes      `json:"scopes,omitempty"`
}

type UserCreateResult struct {
	ID           string        `json:"id"`
	Email        string        `json:"email"`
	Profile      ProfileData   `json:"profile"`
	IDProviders  []string      `json:"id_providers"`
	RequireMFA   bool          `json:"require_mfa"`
	Verified     bool          `json:"verified"`
	LastLoginAt  *string       `json:"last_login_at,omitempty"`
	Disabled     bool          `json:"disabled"`
	MFAProviders []MFAProvider `json:"mfa_providers,omitempty"`
	CreatedAt    string        `json:"created_at"`
}

// @summary Create User
//
// @description Create a user.
//
// @operationId authn_post_v1_user_create
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
//		Authenticator: "My1s+Password",
//		IDProvider: authn.IDPPassword,
//		Profile: profile,
//	}
//
//	resp, err := authncli.User.Create(ctx, input)
func (a *User) Create(ctx context.Context, input UserCreateRequest) (*pangea.PangeaResponse[UserCreateResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/create", &input, &UserCreateResult{})
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
// @description Delete a user by email address.
//
// @operationId authn_post_v1_user_delete
//
// @example
//
//	input := UserDeleteRequest{
//		Email: "joe.user@email.com",
//	}
//
//	authncli.User.Delete(ctx, input)
func (a *User) Delete(ctx context.Context, input UserDeleteRequest) (*pangea.PangeaResponse[UserDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/delete", &input, &UserDeleteResult{})
}

type UserUpdateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID            *string `json:"id,omitempty"`
	Email         *string `json:"email,omitempty"`
	Authenticator *string `json:"authenticator,omitempty"`
	Disabled      *bool   `json:"disabled,omitempty"`
	RequireMFA    *bool   `json:"require_mfa,omitempty"`
}

type UserUpdateResult struct {
	ID           string        `json:"id"`
	Email        string        `json:"email"`
	Profile      ProfileData   `json:"profile"`
	Scopes       *Scopes       `json:"scopes,omitempty"`
	IDProviders  []string      `json:"id_providers"`
	MFAProviders []MFAProvider `json:"mfa_providers,omitempty"`
	RequireMFA   bool          `json:"require_mfa"`
	Verified     bool          `json:"verified"`
	Disabled     bool          `json:"disabled"`
	LastLoginAt  string        `json:"last_login_at,omitempty"`
	CreatedAt    string        `json:"created_at"`
}

// @summary Update user's settings
//
// @description Update user's settings.
//
// @operationId authn_post_v1_user_update
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
	return request.DoPost(ctx, a.Client, "v1/user/update", &input, &UserUpdateResult{})
}

type UserInviteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Inviter    string `json:"inviter"`
	Email      string `json:"email"`
	Callback   string `json:"callback"`
	State      string `json:"state"`
	RequireMFA *bool  `json:"require_mfa,omitempty"`
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
// @operationId authn_post_v1_user_invite
//
// @example
//
//	input := authn.UserInviteRequest{
//		Inviter: "admin@email.com",
//		Email: "joe.user@email.com",
//		Callback: "/callback",
//		State: "pcb_zurr3lkcwdp5keq73htsfpcii5k4zgm7",
//	}
//
//	resp, err := authncli.User.Invite(ctx, input)
func (a *User) Invite(ctx context.Context, input UserInviteRequest) (*pangea.PangeaResponse[UserInviteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/invite", &input, &UserInviteResult{})
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

type UserInfo struct {
	Profile      ProfileData `json:"profile"`
	ID           string      `json:"id"`
	Email        string      `json:"email"`
	Scopes       Scopes      `json:"scopes"`
	IDProviders  []string    `json:"id_providers"`
	MFAProviders []string    `json:"mfa_providers"`
	RequireMFA   bool        `json:"require_mfa"`
	Verified     bool        `json:"verified"`
	Disabled     bool        `json:"disabled"`
	LastLoginAt  *string     `json:"last_login_at,omitempty"`
	CreatedAt    string      `json:"created_at"`
}

type UserListResult struct {
	Users []UserInfo `json:"users"`
	Last  string     `json:"last,omitempty"`
	Count int        `json:"count"`
}

// @summary List Users
//
// @description Look up users by scopes.
//
// @operationId authn_post_v1_user_list
//
// @example
//
//	input := authn.UserListRequest{}
//	resp, err := authncli.User.List(ctx, input)
func (a *User) List(ctx context.Context, input UserListRequest) (*pangea.PangeaResponse[UserListResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/list", &input, &UserListResult{})
}

type LoginToken struct {
	Token     string      `json:"token"`
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Life      int         `json:"life"`
	Expire    string      `json:"expire"`
	Identity  string      `json:"identity"`
	Email     string      `json:"email"`
	Profile   ProfileData `json:"profile"`
	Scopes    Scopes      `json:"scopes"`
	CreatedAt string      `json:"created_at"`
}

type UserLoginResult struct {
	RefreshToken LoginToken  `json:"refresh_token"`
	ActiveToken  *LoginToken `json:"active_token,omitempty"`
}

type UserLoginPasswordRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Email        string       `json:"email"`
	Password     string       `json:"password"`
	ExtraProfile *ProfileData `json:"extra_profile,omitempty"`
}

type UserLoginSocialRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Email        string       `json:"email"`
	Provider     IDProvider   `json:"provider"`
	SocialID     string       `json:"social_id"`
	ExtraProfile *ProfileData `json:"extra_profile,omitempty"`
}

// @summary Login with a password
//
// @description Login a user with a password and return the user's token and information.
//
// @operationId authn_post_v1_user_login_password
//
// @example
//
//	input := authn.UserLoginPasswordRequest{
//		Email: "joe.user@email.com",
//		Password: "My1s+Password",
//		ExtraProfile: &authn.ProfileData{
//			"country": "Argentina",
//		},
//	}
//
//	resp, err := authncli.User.Login.Password(ctx, input)
func (a *UserLogin) Password(ctx context.Context, input UserLoginPasswordRequest) (*pangea.PangeaResponse[UserLoginResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/login/password", &input, &UserLoginResult{})
}

// @summary Login with a social provider
//
// @description Login a user by their social ID and return the user's token and information.
//
// @operationId authn_post_v1_user_login_social
//
// @example
//
//	input := authn.UserLoginSocialRequest{
//		Email: "joe.user@email.com",
//		Provider: authn.IDPGoogle,
//		SocialID: "My1s+Password",
//	}
//
//	resp, err := authncli.User.Login.Social(ctx, input)
func (a *UserLogin) Social(ctx context.Context, input UserLoginSocialRequest) (*pangea.PangeaResponse[UserLoginResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/login/social", &input, &UserLoginResult{})
}

type UserProfileGetRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID    *string `json:"id,omitempty"`
	Email *string `json:"email,omitempty"`
}

type UserProfileGetResult struct {
	ID           string        `json:"id"`
	Email        string        `json:"email"`
	Profile      ProfileData   `json:"profile"`
	IDProviders  []string      `json:"id_providers"`
	MFAProviders []MFAProvider `json:"mfa_providers,omitempty"`
	RequireMFA   bool          `json:"require_mfa"`
	Verified     bool          `json:"verified"`
	LastLoginAt  *string       `json:"last_login_at,omitempty"`
	Disabled     bool          `json:"disabled"`
	CreatedAt    string        `json:"created_at"`
}

// @summary Get user
//
// @description Get user's information.
//
// @operationId authn_post_v1_user_profile_get
//
// @example
//
//	input := authn.UserProfileGetRequest{
//		Email: pangea.String("joe.user@email.com"),
//	}
//
//	resp, err := authncli.User.Profile.Get(ctx, input)
func (a *UserProfile) Get(ctx context.Context, input UserProfileGetRequest) (*pangea.PangeaResponse[UserProfileGetResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/profile/get", &input, &UserProfileGetResult{})
}

type UserProfileUpdateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Profile ProfileData `json:"profile"`
	ID      *string     `json:"id,omitempty"`
	Email   *string     `json:"email,omitempty"`
}

type UserProfileUpdateResult struct {
	ID           string        `json:"id"`
	Email        string        `json:"email"`
	Profile      ProfileData   `json:"profile"`
	IDProviders  []string      `json:"id_providers"`
	MFAProviders []MFAProvider `json:"mfa_providers"`
	RequireMFA   bool          `json:"require_mfa"`
	Verified     bool          `json:"verified"`
	LastLoginAt  *string       `json:"last_login_at,omitempty"`
	Disabled     bool          `json:"disabled"`
	CreatedAt    string        `json:"created_at"`
}

// @summary Update user
//
// @description Update user's information by identity or email.
//
// @operationId authn_post_v1_user_profile_update
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
	return request.DoPost(ctx, a.Client, "v1/user/profile/update", &input, &UserProfileUpdateResult{})
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
// @operationId authn_post_v1_user_invite_list
//
// @example
//
//	input := authn.UserInviteListRequest{}
//	resp, err := authncli.User.Invite.List(ctx, input)
func (a *UserInvite) List(ctx context.Context, input UserInviteListRequest) (*pangea.PangeaResponse[UserInviteListResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/invite/list", &input, &UserInviteListResult{})
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
// @operationId authn_post_v1_user_invite_delete
//
// @example
//
//	input := authn.UserInviteDeleteRequest{
//		ID: "pmc_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.User.Invite.Delete(ctx, input)
func (a *UserInvite) Delete(ctx context.Context, input UserInviteDeleteRequest) (*pangea.PangeaResponse[UserInviteDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/invite/delete", &input, &UserInviteDeleteResult{})
}

type UserPasswordResetRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string `json:"user_id"`
	NewPassword string `json:"new_password"`
}

type UserPasswordResetResult struct {
}

// @summary User Password Reset
//
// @description Manually reset a user's password.
//
// @operationId authn_post_v1_user_password_reset
//
// @example
//
//	input := authn.UserPasswordResetRequest{
//		UserID: "pui_xpkhwpnz2cmegsws737xbsqnmnuwtvm5",
//		NewPassword: "My2n+Password",
//	}
//
//	resp, err := authncli.User.Password.Reset(ctx, input)
func (a *UserPassword) Reset(ctx context.Context, input UserPasswordResetRequest) (*pangea.PangeaResponse[UserPasswordResetResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/password/reset", &input, &UserPasswordResetResult{})
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

// @summary Complete Sign-up/in
//
// @description Complete a login or signup flow.
//
// @operationId authn_post_v1_flow_complete
//
// @example
//
//	input := authn.FlowCompleteRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//	}
//
//	resp, err := authncli.Flow.Complete(ctx, input)
func (a *Flow) Complete(ctx context.Context, input FlowCompleteRequest) (*pangea.PangeaResponse[FlowCompleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/complete", &input, &FlowCompleteResult{})
}

type FlowEnrollMFACompleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID string `json:"flow_id"`
	Code   string `json:"code"`
	Cancel *bool  `json:"cancel,omitempty"`
}

type EnrollMFAStartData struct {
	MFAProviders *[]MFAProvider `json:"mfa_providers,omitempty"`
}

type TOTPSecretData struct {
	QRImage string `json:"qr_image"`
	Secret  string `json:"secret"`
}

type EnrollMFACompleteData struct {
	TOTPSecret TOTPSecretData `json:"totp_secret"`
}

type SocialSignupData struct {
	RedirectURI map[string]string `json:"redirect_uri"`
}

type PasswordSignupData struct {
	PasswordCharsMin  int `json:"password_chars_min"`
	PasswordCharsMax  int `json:"password_chars_max"`
	PasswordLowerMin  int `json:"password_lower_min"`
	PasswordUpperMin  int `json:"passwrod_upper_min"`
	PasswordPunctMin  int `json:"password_punct_min"`
	PasswordNumberMin int `json:"password_number_min"`
}

type VerifyCaptchaData struct {
	SikeKey string `json:"site_key"`
}

type VerifyMFAStartData struct {
	MFAProviders *[]MFAProvider `json:"mfa_providers,omitempty"`
}

type VerifyPasswordData struct {
	PasswordCharsMin  int `json:"password_chars_min"`
	PasswordCharsMax  int `json:"password_chars_max"`
	PasswordLowerMin  int `json:"password_lower_min"`
	PasswordUpperMin  int `json:"passwrod_upper_min"`
	PasswordPunctMin  int `json:"password_punct_min"`
	PasswordNumberMin int `json:"password_number_min"`
}

type SignupData struct {
	SocialSignup   SocialSignupData   `json:"social_signup"`
	PasswordSignup PasswordSignupData `json:"password_signup"`
}

type VerifySocialData struct {
	RedirectURI string `json:"redirect_uri"`
}

type CommonFlowResult struct {
	FlowID            string                 `json:"flow_id,omitempty"`
	NextStep          string                 `json:"next_step"`
	Error             *string                `json:"error,omitempty"`
	Complete          map[string]any         `json:"complete,omitempty"`
	EnrollMFAstart    *EnrollMFAStartData    `json:"enroll_mfa_start,omitempty"`
	EnrollMFAComplete *EnrollMFACompleteData `json:"enroll_mfa_complete,omitempty"`
	Signup            *SignupData            `json:"signup,omitempty"`
	VerifyCaptcha     *VerifyCaptchaData     `json:"verify_captcha,omitempty"`
	VerifyEmail       map[string]any         `json:"verify_email,omitempty"`
	VerifyMFAStart    *VerifyMFAStartData    `json:"verify_mfa_start,omitempty"`
	VerifyMFAComplete map[string]any         `json:"verify_mfa_complete,omitempty"`
	VerifyPassword    *VerifyPasswordData    `json:"verify_password,omitempty"`
	VerifySocial      *VerifySocialData      `json:"verify_social,omitempty"`
	ResetPassword     *VerifyPasswordData    `json:"reset_password,omitempty"`
}

type FlowEnrollMFACompleteResult struct {
	CommonFlowResult
}

// @summary Complete MFA Enrollment
//
// @description Complete MFA enrollment by verifying a trial MFA code.
//
// @operationId authn_post_v1_flow_enroll_mfa_complete
//
// @example
//
//	input := authn.FlowEnrollMFACompleteRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Code: "391423",
//	}
//
//	resp, err := authncli.Flow.Enroll.MFA.Complete(ctx, input)
func (a *FlowEnrollMFA) Complete(ctx context.Context, input FlowEnrollMFACompleteRequest) (*pangea.PangeaResponse[FlowEnrollMFACompleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/enroll/mfa/complete", &input, &FlowEnrollMFACompleteResult{})
}

type FlowResetPasswordRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID   string `json:"flow_id"`
	Password string `json:"password"`
	Cancel   *bool  `json:"cancel,omitempty"`
	CBState  string `json:"cb_state,omitempty"`
	CBCode   string `json:"cb_code,omitempty"`
}

type FlowResetPasswordResult struct {
	CommonFlowResult
}

// @summary Password Reset
//
// @description Reset password during sign-in.
//
// @operationId authn_post_v1_flow_reset_password
//
// @example
//
//	input := authn.FlowResetPasswordRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Password: "My1s+Password",
//	}
//
//	resp, err := authncli.Flow.Reset.Password(ctx, input)
func (a *FlowReset) Password(ctx context.Context, input FlowResetPasswordRequest) (*pangea.PangeaResponse[FlowResetPasswordResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/reset/password", &input, &FlowResetPasswordResult{})
}

type FlowEnrollMFAStartRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID      string      `json:"flow_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Phone       string      `json:"phone,omitempty"`
}

type FlowEnrollMFAStartResult struct {
	CommonFlowResult
}

// @summary Start MFA Enrollment
//
// @description Start the process of enrolling an MFA.
//
// @operationId authn_post_v1_flow_enroll_mfa_start
//
// @example
//
//	input := authn.FlowEnrollMFAStartRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		MFAProvider: authn.IDPSMSOTP,
//		Phone: "1-808-555-0173",
//	}
//
//	resp, err := authncli.Flow.Enroll.MFA.Start(ctx, input)
func (a *FlowEnrollMFA) Start(ctx context.Context, input FlowEnrollMFAStartRequest) (*pangea.PangeaResponse[FlowEnrollMFAStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/enroll/mfa/start", &input, &FlowEnrollMFAStartResult{})
}

type FlowSignupPasswordRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID    string `json:"flow_id"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type FlowSignupPasswordResult struct {
	CommonFlowResult
}

// @summary Password Sign-up
//
// @description Signup a new account using a password.
//
// @operationId authn_post_v1_flow_signup_password
//
// @example
//
//	input := authn.FlowSignupPasswordRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Password: "My1s+Password",
//		FirstName: "Joe",
//		LastName: "User",
//	}
//
//	resp, err := authncli.Flow.Signup.Password(ctx, input)
func (a *FlowSignup) Password(ctx context.Context, input FlowSignupPasswordRequest) (*pangea.PangeaResponse[FlowSignupPasswordResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/signup/password", &input, &FlowSignupPasswordResult{})
}

type FlowSignupSocialRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID  string `json:"flow_id"`
	CBState string `json:"cb_state"`
	CBCode  string `json:"cb_code"`
}

type FlowSignupSocialResult struct {
	CommonFlowResult
}

// @summary Social Sign-up
//
// @description Signup a new account using a social provider.
//
// @operationId authn_post_v1_flow_signup_social
//
// @example
//
//	input := authn.FlowSignupSocialRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		CBState: "pcb_zurr3lkcwdp5keq73htsfpcii5k4zgm7",
//		CBCode: "poc_fwg3ul4db1jpivexru3wyj354u9ej5e2",
//	}
//
//	resp, err := authncli.Flow.Signup.Social(ctx, input)
func (a *FlowSignup) Social(ctx context.Context, input FlowSignupSocialRequest) (*pangea.PangeaResponse[FlowSignupSocialResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/signup/social", &input, &FlowSignupSocialResult{})
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

// @summary Start a sign-up/in
//
// @description Start a new signup or signin flow.
//
// @operationId authn_post_v1_flow_start
//
// @example
//
//	fts := []FlowType{FTsignin,FTsignup}
//	input := authn.FlowStartRequest{
//		CBURI: "https://www.myserver.com/callback",
//		Email: "joe.user@email.com",
//		FlowTypes: fts,
//		Provider: &authn.IDPPassword,
//	}
//
//	resp, cli := authncli.Flow.Start(ctx, input)
func (a *Flow) Start(ctx context.Context, input FlowStartRequest) (*pangea.PangeaResponse[FlowStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/start", &input, &FlowStartResult{})
}

type FlowVerifyCaptchaRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID string `json:"flow_id"`
	Code   string `json:"code"`
}

type FlowVerifyCaptchaResult struct {
	CommonFlowResult
}

// @summary Verify Captcha
//
// @description Verify a CAPTCHA during a signup or signin flow.
//
// @operationId authn_post_v1_flow_verify_captcha
//
// @example
//
//	input := authn.FlowVerifyCaptchaRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Code: "SOMEREALLYLONGANDOPAQUESTRINGFROMCAPTCHAVERIFICATION",
//	}
//
//	resp, err := authncli.Flow.Verify.Captcha(ctx, input)
func (a *FlowVerify) Captcha(ctx context.Context, input FlowVerifyCaptchaRequest) (*pangea.PangeaResponse[FlowVerifyCaptchaResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/captcha", &input, &FlowVerifyCaptchaResult{})
}

type FlowVerifyEmailRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID  string `json:"flow_id"`
	CBState string `json:"cb_state,omitempty"`
	CBCode  string `json:"cb_code,omitempty"`
}

type FlowVerifyEmailResult struct {
	CommonFlowResult
}

// @summary Verify Email Address
//
// @description Verify an email address during a signup or signin flow.
//
// @operationId authn_post_v1_flow_verify_email
//
// @example
//
//	input := authn.FlowVerifyEmailRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		CBState: "pcb_zurr3lkcwdp5keq73htsfpcii5k4zgm7",
//		CBCode: "poc_fwg3ul4db1jpivexru3wyj354u9ej5e2",
//	}
//
//	resp, err := authncli.Flow.Verify.Email(ctx, input)
func (a *FlowVerify) Email(ctx context.Context, input FlowVerifyEmailRequest) (*pangea.PangeaResponse[FlowVerifyEmailResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/email", &input, &FlowVerifyEmailResult{})
}

type FlowVerifyMFACompleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID string  `json:"flow_id"`
	Code   *string `json:"code,omitempty"`
	Cancel *bool   `json:"cancel,omitempty"`
}

type FlowVerifyMFACompleteResult struct {
	CommonFlowResult
}

// @summary Complete MFA Verification
//
// @description Complete MFA verification.
//
// @operationId authn_post_v1_flow_verify_mfa_complete
//
// @example
//
//	input := authn.FlowVerifyMFACompleteRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Code: pangea.String("999999"),
//	}
//
//	resp, err := authncli.Flow.Verify.MFA.Complete(ctx, input)
func (a *FlowVerifyMFA) Complete(ctx context.Context, input FlowVerifyMFACompleteRequest) (*pangea.PangeaResponse[FlowVerifyMFACompleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/mfa/complete", &input, &FlowVerifyMFACompleteResult{})
}

type FlowVerifyMFAStartRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID      string      `json:"flow_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
}

type FlowVerifyMFAStartResult struct {
	CommonFlowResult
}

// @summary Start MFA Verification
//
// @description Start the process of MFA verification.
//
// @operationId authn_post_v1_flow_verify_mfa_start
//
// @example
//
//	input := authn.FlowVerifyMFAStartRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		MFAProvider: authn.MFAPTOTP,
//	}
//
//	resp, err := authncli.Flow.Verify.MFA.Start(ctx, input)
func (a *FlowVerifyMFA) Start(ctx context.Context, input FlowVerifyMFAStartRequest) (*pangea.PangeaResponse[FlowVerifyMFAStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/mfa/start", &input, &FlowVerifyMFAStartResult{})
}

type FlowVerifyPasswordRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID   string  `json:"flow_id"`
	Password *string `json:"password,omitempty"`
	Reset    *bool   `json:"reset,omitempty"`
}

type FlowVerifyPasswordResult struct {
	CommonFlowResult
}

// @summary Password Sign-in
//
// @description Sign in with a password.
//
// @operationId authn_post_v1_flow_verify_password
//
// @example
//
//	input := authn.FlowVerifyPasswordRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		Password: pangea.String("My1s+Password"),
//	}
//
//	resp, err := authncli.Flow.Verify.Password(ctx, input)
func (a *FlowVerify) Password(ctx context.Context, input FlowVerifyPasswordRequest) (*pangea.PangeaResponse[FlowVerifyPasswordResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/password", &input, &FlowVerifyPasswordResult{})
}

type FlowVerifySocialRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID  string `json:"flow_id"`
	CBState string `json:"cb_state"`
	CBCode  string `json:"cb_code"`
}

type FlowVerifySocialResult struct {
	CommonFlowResult
}

// @summary Social Sign-in
//
// @description Signin with a social provider.
//
// @operationId authn_post_v1_flow_verify_social
//
// @example
//
//	input := authn.FlowVerifySocialRequest{
//		FlowID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		CBState: "pcb_zurr3lkcwdp5keq73htsfpcii5k4zgm7",
//		CBCode: "poc_fwg3ul4db1jpivexru3wyj354u9ej5e2",
//	}
//
//	resp, err := authncli.Flow.Verify.Social(ctx, input)
func (a *FlowVerify) Social(ctx context.Context, input FlowVerifySocialRequest) (*pangea.PangeaResponse[FlowVerifySocialResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/social", &input, &FlowVerifySocialResult{})
}

type UserMFADeleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
}

type UserMFADeleteResult struct {
}

// @summary Delete MFA Enrollment
//
// @description Delete MFA enrollment for a user.
//
// @operationId authn_post_v1_user_mfa_delete
//
// @example
//
//	input := authn.UserMFADeleteRequest{
//		UserID: "pui_zgp532cx6opljeavvllmbi3iwmq72f7f",
//		MFAProvider: authn.MFAPTOTP,
//	}
//
//	resp, err := authncli.User.MFA.Delete(ctx, input)
func (a *UserMFA) Delete(ctx context.Context, input UserMFADeleteRequest) (*pangea.PangeaResponse[UserMFADeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/delete", &input, &UserMFADeleteResult{})
}

type UserMFAEnrollRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Code        string      `json:"code"`
}

type UserMFAEnrollResult struct {
}

// @summary Enroll In MFA
//
// @description Enroll in MFA for a user by proving the user has access to an MFA verification code.
//
// @operationId authn_post_v1_user_mfa_enroll
//
// @example
//
//	input := authn.UserMFAEnrollRequest{
//		UserID: "pui_zgp532cx6opljeavvllmbi3iwmq72f7f",
//		MFAProvider: authn.MFAPTOTP,
//		Code: "999999",
//	}
//
//	resp, err := authncli.User.MFA.Enroll(ctx, input)
func (a *UserMFA) Enroll(ctx context.Context, input UserMFAEnrollRequest) (*pangea.PangeaResponse[UserMFAEnrollResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/enroll", &input, &UserMFAEnrollResult{})
}

type UserMFAStartRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Enroll      *bool       `json:"enroll,omitempty"`
	Phone       *string     `json:"phone,omitempty"`
}

type UserMFAStartTOTPSecret struct {
	QRImage string `json:"qr_image"`
	Secret  string `json:"secret"`
}

type UserMFAStartResult struct {
	TOTPSecret *UserMFAStartTOTPSecret `json:"totp_secret,omitempty"`
}

// @summary Start MFA Verification
//
// @description Start MFA verification for a user, generating a new one-time code, and sending it if necessary. When enrolling TOTP, this returns the TOTP secret.
//
// @operationId authn_post_v1_user_mfa_start
//
// @example
//
//	input := authn.UserMFAStartRequest{
//		UserID: "pfl_dxiqyuq7ndc5ycjwdgmguwuodizcaqhh",
//		MFAProvider: authn.MFAPTOTP,
//	}
//
//	resp, err := authncli.User.MFA.Start(ctx, input)
func (a *UserMFA) Start(ctx context.Context, input UserMFAStartRequest) (*pangea.PangeaResponse[UserMFAStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/start", &input, &UserMFAStartResult{})
}

type UserMFAVerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Code        string      `json:"code"`
}

type UserMFAVerifyResult struct {
}

// @summary Verify An MFA Code
//
// @description Verify that the user has access to an MFA verification code.
//
// @operationId authn_post_v1_user_mfa_verify
//
// @example
//
//	input := authn.UserMFAVerifyRequest{
//		UserID: "pui_zgp532cx6opljeavvllmbi3iwmq72f7f",
//		MFAProvider: authn.MFAPTOTP,
//		Code: "999999",
//	}
//
//	resp, err := authncli.User.MFA.Verify(ctx, input)
func (a *UserMFA) Verify(ctx context.Context, input UserMFAVerifyRequest) (*pangea.PangeaResponse[UserMFAVerifyResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/verify", &input, &UserMFAVerifyResult{})
}

type UserVerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	IDProvider    IDProvider `json:"id_provider"`
	Email         string     `json:"email"`
	Authenticator string     `json:"authenticator"`
}

type UserVerifyResult struct {
	ID           string      `json:"id"`
	Email        string      `json:"email"`
	Profile      ProfileData `json:"profile"`
	Scopes       Scopes      `json:"scopes"`
	IDProviders  []string    `json:"id_providers"`
	MFAProviders []string    `json:"mfa_providers"`
	RequireMFA   bool        `json:"require_mfa"`
	Verified     bool        `json:"verified"`
	Disable      bool        `json:"disable"`
	LastLoginAt  *string     `json:"last_login_at,omitempty"`
	CreatedAt    string      `json:"created_at"`
}

// @summary Verify User
//
// @description Verify a user's primary authentication.
//
// @operationId authn_post_v1_user_verify
//
// @example
//
//	input := authn.UserVerifyRequest{
//		IDProvider: authn.IDPPassword,
//		Email: "joe.user@email.com",
//		Authenticator: "My1s+Password",
//	}
//
//	resp, err := authncli.User.Verify(ctx, input)
func (a *User) Verify(ctx context.Context, input UserVerifyRequest) (*pangea.PangeaResponse[UserVerifyResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/verify", &input, &UserVerifyResult{})
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
// @operationId authn_post_v1_client_session_invalidate
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
	return request.DoPost(ctx, a.Client, "v1/client/session/invalidate", &input, &ClientSessionInvalidateResult{})
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
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Life      int         `json:"life"`
	Expire    string      `json:"expire"`
	Email     string      `json:"email"`
	Scopes    Scopes      `json:"scopes"`
	Profile   ProfileData `json:"profile"`
	CreatedAt string      `json:"created_at"`
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
// @operationId authn_post_v1_client_session_list
//
// @example
//
//	input := authn.ClientSessionListRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.Client.Session.List(ctx, input)
func (a *ClientSession) List(ctx context.Context, input ClientSessionListRequest) (*pangea.PangeaResponse[SessionListResult], error) {
	return request.DoPost(ctx, a.Client, "v1/client/session/list", &input, &SessionListResult{})
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
// @operationId authn_post_v1_client_session_logout
//
// @example
//
//	input := authn.ClientSessionLogoutRequest{
//		Token: "ptu_wuk7tvtpswyjtlsx52b7yyi2l7zotv4a",
//	}
//
//	resp, err := authncli.Client.Session.Logout(ctx, input)
func (a *ClientSession) Logout(ctx context.Context, input ClientSessionLogoutRequest) (*pangea.PangeaResponse[ClientSessionLogoutResult], error) {
	return request.DoPost(ctx, a.Client, "v1/client/session/logout", &input, &ClientSessionLogoutResult{})
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
// @operationId authn_post_v1_client_session_refresh
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
	return request.DoPost(ctx, a.Client, "v1/client/session/refresh", &input, &ClientSessionRefreshResult{})
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
// @operationId authn_post_v1_session_list
//
// @example
//
//	input := authn.SessionListRequest{}
//	resp, err := authn.Session.List(ctx, input)
func (a *Session) List(ctx context.Context, input SessionListRequest) (*pangea.PangeaResponse[SessionListResult], error) {
	return request.DoPost(ctx, a.Client, "v1/session/list", &input, &SessionListResult{})
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
// @operationId authn_post_v1_session_invalidate
//
// @example
//
//	input := authn.SessionInvalidateRequest{
//		SessionID: "pmt_zppkzrjguxyblaia6itbiesejn7jejnr",
//	}
//
//	resp, err := authncli.Session.Invalidate(ctx, input)
func (a *Session) Invalidate(ctx context.Context, input SessionInvalidateRequest) (*pangea.PangeaResponse[SessionInvalidateResult], error) {
	return request.DoPost(ctx, a.Client, "v1/session/invalidate", &input, &SessionInvalidateResult{})
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
// @operationId authn_post_v1_session_logout
//
// @example
//
//	input := authn.SessionLogoutRequest{
//		UserID: "pui_xpkhwpnz2cmegsws737xbsqnmnuwtvm5",
//	}
//
//	resp, err := authncli.Session.Logout(ctx, input)
func (a *Session) Logout(ctx context.Context, input SessionLogoutRequest) (*pangea.PangeaResponse[SessionLogoutResult], error) {
	return request.DoPost(ctx, a.Client, "v1/session/logout", &input, &SessionLogoutResult{})
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

// TODO: docs
func (a *Agreements) Create(ctx context.Context, input AgreementCreateRequest) (*pangea.PangeaResponse[AgreementCreateResult], error) {
	return request.DoPost(ctx, a.Client, "v1/agreements/create", &input, &AgreementCreateResult{})
}

type AgreementDeleteRequest struct {
	pangea.BaseRequest

	Type AgreementType `json:"type"`
	ID   string        `json:"id"`
}

type AgreementDeleteResult struct{}

// TODO: docs
func (a *Agreements) Delete(ctx context.Context, input AgreementDeleteRequest) (*pangea.PangeaResponse[AgreementDeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/agreements/delete", &input, &AgreementDeleteResult{})
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

// TODO: docs
func (a *Agreements) List(ctx context.Context, input AgreementListRequest) (*pangea.PangeaResponse[AgreementListResult], error) {
	return request.DoPost(ctx, a.Client, "v1/agreements/list", &input, &AgreementListResult{})
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

// TODO: docs
func (a *Agreements) Update(ctx context.Context, input AgreementUpdateRequest) (*pangea.PangeaResponse[AgreementUpdateResult], error) {
	return request.DoPost(ctx, a.Client, "v1/agreements/update", &input, &AgreementUpdateResult{})
}
