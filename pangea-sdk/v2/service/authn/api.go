package authn

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/request"
	v "github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
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

func (a *Client) Userinfo(ctx context.Context, input ClientUserinfoRequest) (*pangea.PangeaResponse[ClientUserinfoResult], error) {
	return request.DoPost(ctx, a.client, "v1/client/userinfo", &input, &ClientUserinfoResult{})
}

type ClientJWKSResult struct {
	Keys []v.JWT `json:"keys"`
}

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
type Filter map[string]any

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

type UserListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Filter  Filter          `json:"filter,omitempty"`
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

func (a *UserLogin) Password(ctx context.Context, input UserLoginPasswordRequest) (*pangea.PangeaResponse[UserLoginResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/login/password", &input, &UserLoginResult{})
}

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

type UserInviteListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Filter  Filter                `json:"filter,omitempty"`
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

func (a *UserPassword) Reset(ctx context.Context, input UserPasswordResetRequest) (*pangea.PangeaResponse[UserPasswordResetResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/password/reset", &input, &UserPasswordResetResult{})
}

// #   - path: authn::/v1/flow/complete
// # https://dev.pangea.cloud/docs/api/authn#complete-a-login-or-signup-flow
type FlowCompleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID string `json:"flow_id"`
}

type FlowCompleteResult struct {
	RefreshToken LoginToken `json:"refresh_token"`
	ActiveToken  LoginToken `json:"active_token"`
}

func (a *Flow) Complete(ctx context.Context, input FlowCompleteRequest) (*pangea.PangeaResponse[FlowCompleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/complete", &input, &FlowCompleteResult{})
}

// #   - path: authn::/v1/flow/enroll/mfa/complete
// # https://dev.pangea.cloud/docs/api/authn#complete-mfa-enrollment-by-verifying-a-trial-mfa-code

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

func (a *FlowReset) Password(ctx context.Context, input FlowResetPasswordRequest) (*pangea.PangeaResponse[FlowResetPasswordResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/reset/password", &input, &FlowResetPasswordResult{})
}

// #   - path: authn::/v1/flow/enroll/mfa/start
// # https://dev.pangea.cloud/docs/api/authn#start-the-process-of-enrolling-an-mfa
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

func (a *FlowEnrollMFA) Start(ctx context.Context, input FlowEnrollMFAStartRequest) (*pangea.PangeaResponse[FlowEnrollMFAStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/enroll/mfa/start", &input, &FlowEnrollMFAStartResult{})
}

// #   - path: authn::/v1/flow/signup/password
// # https://dev.pangea.cloud/docs/api/authn#signup-a-new-account-using-a-password
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

func (a *FlowSignup) Password(ctx context.Context, input FlowSignupPasswordRequest) (*pangea.PangeaResponse[FlowSignupPasswordResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/signup/password", &input, &FlowSignupPasswordResult{})
}

// #   - path: authn::/v1/flow/signup/social
// # https://dev.pangea.cloud/docs/api/authn#signup-a-new-account-using-a-social-provider
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

func (a *FlowSignup) Social(ctx context.Context, input FlowSignupSocialRequest) (*pangea.PangeaResponse[FlowSignupSocialResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/signup/social", &input, &FlowSignupSocialResult{})
}

// #   - path: authn::/v1/flow/start
// # https://dev.pangea.cloud/docs/api/authn#start-a-new-signup-or-signin-flow
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

func (a *Flow) Start(ctx context.Context, input FlowStartRequest) (*pangea.PangeaResponse[FlowStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/start", &input, &FlowStartResult{})
}

// #   - path: authn::/v1/flow/verify/captcha
// # https://dev.pangea.cloud/docs/api/authn#verify-a-captcha-during-a-signup-or-signin-flow
type FlowVerifyCaptchaRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID string `json:"flow_id"`
	Code   string `json:"code"`
}

type FlowVerifyCaptchaResult struct {
	CommonFlowResult
}

func (a *FlowVerify) Captcha(ctx context.Context, input FlowVerifyCaptchaRequest) (*pangea.PangeaResponse[FlowVerifyCaptchaResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/captcha", &input, &FlowVerifyCaptchaResult{})
}

// #   - path: authn::/v1/flow/verify/email
// # https://dev.pangea.cloud/docs/api/authn#verify-an-email-address-during-a-signup-or-signin-flow
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

func (a *FlowVerify) Email(ctx context.Context, input FlowVerifyEmailRequest) (*pangea.PangeaResponse[FlowVerifyEmailResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/email", &input, &FlowVerifyEmailResult{})
}

// #   - path: authn::/v1/flow/verify/mfa/complete
// # https://dev.pangea.cloud/docs/api/authn#complete-mfa-verification
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

func (a *FlowVerifyMFA) Complete(ctx context.Context, input FlowVerifyMFACompleteRequest) (*pangea.PangeaResponse[FlowVerifyMFACompleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/mfa/complete", &input, &FlowVerifyMFACompleteResult{})
}

// #   - path: authn::/v1/flow/verify/mfa/start
// # https://dev.pangea.cloud/docs/api/authn#start-the-process-of-mfa-verification
type FlowVerifyMFAStartRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	FlowID      string      `json:"flow_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
}

type FlowVerifyMFAStartResult struct {
	CommonFlowResult
}

func (a *FlowVerifyMFA) Start(ctx context.Context, input FlowVerifyMFAStartRequest) (*pangea.PangeaResponse[FlowVerifyMFAStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/mfa/start", &input, &FlowVerifyMFAStartResult{})
}

// #   - path: authn::/v1/flow/verify/password
// # https://dev.pangea.cloud/docs/api/authn#sign-in-with-a-password
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

func (a *FlowVerify) Password(ctx context.Context, input FlowVerifyPasswordRequest) (*pangea.PangeaResponse[FlowVerifyPasswordResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/password", &input, &FlowVerifyPasswordResult{})
}

// #   - path: authn::/v1/flow/verify/social
// # https://dev.pangea.cloud/docs/api/authn#signin-with-a-social-provider
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

func (a *FlowVerify) Social(ctx context.Context, input FlowVerifySocialRequest) (*pangea.PangeaResponse[FlowVerifySocialResult], error) {
	return request.DoPost(ctx, a.Client, "v1/flow/verify/social", &input, &FlowVerifySocialResult{})
}

// #   - path: authn::/v1/user/mfa/delete
// # https://dev.pangea.cloud/docs/api/authn#delete-mfa-enrollment-for-a-user
type UserMFADeleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
}

type UserMFADeleteResult struct {
}

func (a *UserMFA) Delete(ctx context.Context, input UserMFADeleteRequest) (*pangea.PangeaResponse[UserMFADeleteResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/delete", &input, &UserMFADeleteResult{})
}

// #   - path: authn::/v1/user/mfa/enroll
// # https://dev.pangea.cloud/docs/api/authn#enroll-mfa-for-a-user
type UserMFAEnrollRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Code        string      `json:"code"`
}

type UserMFAEnrollResult struct {
}

func (a *UserMFA) Enroll(ctx context.Context, input UserMFAEnrollRequest) (*pangea.PangeaResponse[UserMFAEnrollResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/enroll", &input, &UserMFAEnrollResult{})
}

// #   - path: authn::/v1/user/mfa/start
// # https://dev.pangea.cloud/docs/api/authn#start-mfa-verification-for-a-user
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

func (a *UserMFA) Start(ctx context.Context, input UserMFAStartRequest) (*pangea.PangeaResponse[UserMFAStartResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/start", &input, &UserMFAStartResult{})
}

// #   - path: authn::/v1/user/mfa/verify
// # https://dev.pangea.cloud/docs/api/authn#verify-an-mfa-code
type UserMFAVerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Code        string      `json:"code"`
}

type UserMFAVerifyResult struct {
}

func (a *UserMFA) Verify(ctx context.Context, input UserMFAVerifyRequest) (*pangea.PangeaResponse[UserMFAVerifyResult], error) {
	return request.DoPost(ctx, a.Client, "v1/user/mfa/verify", &input, &UserMFAVerifyResult{})
}

// #   - path: authn::/v1/user/verify
// # https://dev.pangea.cloud/docs/api/authn#verify-a-user
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

type ClientSessionListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Token   string             `json:"token"`
	Filter  Filter             `json:"filter,omitempty"`
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

func (a *ClientSession) Refresh(ctx context.Context, input ClientSessionRefreshRequest) (*pangea.PangeaResponse[ClientSessionRefreshResult], error) {
	return request.DoPost(ctx, a.Client, "v1/client/session/refresh", &input, &ClientSessionRefreshResult{})
}

type SessionListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Filter  Filter             `json:"filter,omitempty"`
	Last    string             `json:"last,omitempty"`
	Order   ItemOrder          `json:"order,omitempty"`
	OrderBy SessionListOrderBy `json:"order_by,omitempty"`
}

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

func (a *Session) Logout(ctx context.Context, input SessionLogoutRequest) (*pangea.PangeaResponse[SessionLogoutResult], error) {
	return request.DoPost(ctx, a.Client, "v1/session/logout", &input, &SessionLogoutResult{})
}
