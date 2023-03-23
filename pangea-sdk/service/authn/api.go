package authn

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type ClientUserinfoResult struct {
	Token     string            `json:"token"`
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Life      string            `json:"life"`
	Expire    string            `json:"expire"`
	Identity  string            `json:"identity"`
	Email     string            `json:"email"`
	Scopes    *[]string         `json:"scopes,omitempty"`
	Profile   map[string]string `json:"profile"`
	CreatedAt string            `json:"created_at"`
}

type ClientUserinfoRequest struct {
	Code string `json:"code"`
}

func (a *Client) Userinfo(ctx context.Context, input ClientUserinfoRequest) (*pangea.PangeaResponse[ClientUserinfoResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/client/userinfo", input)
	if err != nil {
		return nil, err
	}

	var out ClientUserinfoResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[ClientUserinfoResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type PasswordUpdateRequest struct {
	Email     string `json:"email"`
	OldSecret string `json:"old_secret"`
	NewSecret string `json:"new_secret"`
}

type PasswordUpdateResult struct {
}

func (a *Password) Update(ctx context.Context, input PasswordUpdateRequest) (*pangea.PangeaResponse[PasswordUpdateResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/password/update", input)
	if err != nil {
		return nil, err
	}

	var out PasswordUpdateResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[PasswordUpdateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
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

type UserCreateRequest struct {
	Email         string       `json:"email"`
	Authenticator string       `json:"authenticator"`
	IDProvider    string       `json:"id_provider"`
	Verified      *bool        `json:"verified,omitempty"`
	RequireMFA    *bool        `json:"require_mfa,omitempty"`
	Profile       *UserProfile `json:"profile,omitempty"`
	Scopes        *[]string    `json:"scopes,omitempty"`
}

type UserCreateResult struct {
	Identity     string            `json:"identity"`
	Email        string            `json:"email"`
	Profile      map[string]string `json:"profile"`
	IDProvider   IDProvider        `json:"id_provider"`
	RequireMFA   bool              `json:"require_mfa"`
	Verified     bool              `json:"verified"`
	LastLoginAt  *bool             `json:"last_login_at,omitempty"`
	Disabled     *bool             `json:"disabled"`
	MFAProviders *[]MFAProvider    `json:"mfa_providers,omitempty"`
}

func (a *User) Create(ctx context.Context, input UserCreateRequest) (*pangea.PangeaResponse[UserCreateResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/create", input)
	if err != nil {
		return nil, err
	}

	var out UserCreateResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserCreateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserDeleteRequest struct {
	Email string `json:"email"`
}

type UserDeleteResult struct {
}

func (a *User) Delete(ctx context.Context, input UserDeleteRequest) (*pangea.PangeaResponse[UserDeleteResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/delete", input)
	if err != nil {
		return nil, err
	}

	var out UserDeleteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserDeleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserUpdateRequest struct {
	Identity      *string `json:"identity,omitempty"`
	Email         *string `json:"email,omitempty"`
	Authenticator *string `json:"authenticator,omitempty"`
	Disabled      *bool   `json:"disabled,omitempty"`
	RequireMFA    *bool   `json:"require_mfa,omitempty"`
}

type UserUpdateResult struct {
	Identity     string            `json:"identity"`
	Email        string            `json:"email"`
	Profile      map[string]string `json:"profile"`
	Scopes       *[]string         `json:"scopes,omitempty"`
	IDProvider   IDProvider        `json:"id_provider"`
	MFAProviders *[]MFAProvider    `json:"mfa_providers,omitempty"`
	RequireMFA   bool              `json:"require_mfa"`
	Verified     bool              `json:"verified"`
	Disabled     bool              `json:"disabled"`
	LastLoginAt  string            `json:"last_login_at"`
}

func (a *User) Update(ctx context.Context, input UserUpdateRequest) (*pangea.PangeaResponse[UserUpdateResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/update", input)
	if err != nil {
		return nil, err
	}

	var out UserUpdateResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserUpdateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserInviteRequest struct {
	Inviter    string  `json:"inviter"`
	Email      string  `json:"email"`
	Callback   string  `json:"callback"`
	State      string  `json:"state"`
	InviteOrg  *string `json:"invite_org,omitempty"`
	RequireMFA *bool   `json:"require_mfa,omitempty"`
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
	req, err := a.Client.NewRequest("POST", "v1/user/invite", input)
	if err != nil {
		return nil, err
	}

	var out UserInviteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserInviteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserListRequest struct {
	Scopes     []string `json:"scopes"`
	GlobScopes []string `json:"glob_scopes"`
}

type UserInfo struct {
	Profile  UserProfile `json:"profile"`
	Identity string      `json:"identity"`
	Email    string      `json:"email"`
	Scopes   []string    `json:"scopes"`
}

type UserListResult struct {
	Users []UserInfo `json:"users"`
}

func (a *User) List(ctx context.Context, input UserListRequest) (*pangea.PangeaResponse[UserListResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/list", input)
	if err != nil {
		return nil, err
	}

	var out UserListResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserListResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserLoginRequest struct {
	Email  string    `json:"email"`
	Secret string    `json:"secret"`
	Scopes *[]string `json:"scopes,omitempty"`
}

type LoginToken struct {
	Token     string      `json:"token"`
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Life      int         `json:"life"`
	Expire    string      `json:"expire"`
	Identity  string      `json:"identity"`
	Email     string      `json:"email"`
	Profile   UserProfile `json:"profile"`
	Scopes    []string    `json:"scopes"`
	CreatedAt string      `json:"created_at"`
}

type UserLoginResult struct {
	RefreshToken LoginToken `json:"refres_token"`
	ActiveToken  LoginToken `json:"active_token"`
}

type UserLoginPasswordRequest struct {
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	ExtraProfile *UserProfile `json:"extra_profile,omitempty"`
}

type UserLoginSocialRequest struct {
	Email        string       `json:"email"`
	Provider     IDProvider   `json:"provider"`
	SocialID     string       `json:"social_id"`
	ExtraProfile *UserProfile `json:"extra_profile,omitempty"`
}

func (a *UserLogin) Password(ctx context.Context, input UserLoginPasswordRequest) (*pangea.PangeaResponse[UserLoginResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/login/password", input)
	if err != nil {
		return nil, err
	}

	var out UserLoginResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserLoginResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

func (a *UserLogin) Social(ctx context.Context, input UserLoginSocialRequest) (*pangea.PangeaResponse[UserLoginResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/login/social", input)
	if err != nil {
		return nil, err
	}

	var out UserLoginResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserLoginResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserProfileGetRequest struct {
	Identity *string `json:"identity,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type UserProfileGetResult struct {
	Identity     string            `json:"identity"`
	Email        string            `json:"email"`
	Profile      map[string]string `json:"profile"`
	IDProvider   IDProvider        `json:"id_provider"`
	MFAProviders *[]MFAProvider    `json:"mfa_providers,omitempty"`
	RequireMFA   bool              `json:"require_mfa"`
	Verified     bool              `json:"verified"`
	LastLoginAt  string            `json:"last_login_at"`
	Disabled     *bool             `json:"disabled,omitempty"`
}

func (a *UserProfile) Get(ctx context.Context, input UserProfileGetRequest) (*pangea.PangeaResponse[UserProfileGetResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/profile/get", input)
	if err != nil {
		return nil, err
	}

	var out UserProfileGetResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserProfileGetResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserProfileUpdateRequest struct {
	Profile    map[string]string `json:"profile"`
	Identity   *string           `json:"identity,omitempty"`
	Email      *string           `json:"email,omitempty"`
	RequireMFA *bool             `json:"require_mfa,omitempty"`
	Disabled   *bool             `json:"disabled,omitempty"`
	Verified   *bool             `json:"verified,omitempty"`
}

type UserProfileUpdateResult struct {
	Identity     string            `json:"identity"`
	Email        string            `json:"email"`
	Profile      map[string]string `json:"profile"`
	IDProvider   IDProvider        `json:"id_provider"`
	MFAProviders *[]MFAProvider    `json:"mfa_providers"`
	RequireMFA   bool              `json:"require_mfa"`
	Verified     bool              `json:"verified"`
	LastLoginAt  string            `json:"last_login_at"`
	Disabled     *bool             `json:"disabled,omitempty"`
}

func (a *UserProfile) Update(ctx context.Context, input UserProfileUpdateRequest) (*pangea.PangeaResponse[UserProfileUpdateResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/profile/update", input)
	if err != nil {
		return nil, err
	}

	var out UserProfileUpdateResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserProfileUpdateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
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

type UserInviteListResult struct {
	Invites []UserInviteData `json:"invites"`
}

func (a *UserInvite) List(ctx context.Context) (*pangea.PangeaResponse[UserInviteListResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/invite/list", make(map[string]string))
	if err != nil {
		return nil, err
	}

	var out UserInviteListResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserInviteListResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type UserInviteDeleteRequest struct {
	ID string `json:"id"`
}

type UserInviteDeleteResult struct {
}

func (a *UserInvite) Delete(ctx context.Context, input UserInviteDeleteRequest) (*pangea.PangeaResponse[UserInviteDeleteResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/invite/delete", input)
	if err != nil {
		return nil, err
	}

	var out UserInviteDeleteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserInviteDeleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/complete
// # https://dev.pangea.cloud/docs/api/authn#complete-a-login-or-signup-flow
type FlowCompleteRequest struct {
	FlowID string `json:"flow_id"`
}

type FlowCompleteResult struct {
	RefreshToken LoginToken `json:"refresh_token"`
	ActiveToken  LoginToken `json:"active_token"`
}

func (a *Flow) Complete(ctx context.Context, input FlowCompleteRequest) (*pangea.PangeaResponse[FlowCompleteResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/complete", input)
	if err != nil {
		return nil, err
	}

	var out FlowCompleteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowCompleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/enroll/mfa/complete
// # https://dev.pangea.cloud/docs/api/authn#complete-mfa-enrollment-by-verifying-a-trial-mfa-code

type FlowEnrollMFACompleteRequest struct {
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
	RedirectURI string `json:"redirect_uri"`
}

type PasswordSignupData struct {
	PasswordCharsMin int `json:"password_chars_min"`
	PasswordCharsMax int `json:"password_chars_max"`
	PasswordLowerMin int `json:"password_lower_min"`
	PasswordUpperMin int `json:"passwrod_upper_min"`
	PasswordPunctMin int `json:"password_punct_min"`
}

type VerifyCaptchaData struct {
	SikeKey string `json:"site_key"`
}

type VerifyMFAStartData struct {
	MFAProviders *[]MFAProvider `json:"mfa_providers,omitempty"`
}

type VerifyPasswordData struct {
	PasswordCharsMin int `json:"password_chars_min"`
	PasswordCharsMax int `json:"password_chars_max"`
	PasswordLowerMin int `json:"password_lower_min"`
	PasswordUpperMin int `json:"passwrod_upper_min"`
	PasswordPunctMin int `json:"password_punct_min"`
}

type SignupData struct {
	SocialSignup   SocialSignupData   `json:"social_signup"`
	PasswordSignup PasswordSignupData `json:"password_signup"`
}

type VerifySocialData struct {
	RedirectURI string `json:"redirect_uri"`
}

type CommonFlowResult struct {
	FlowID            string                 `json:"flow_id"`
	NextStep          string                 `json:"next_step"`
	Error             *string                `json:"error,omitempty"`
	Complete          *map[string]any        `json:"complete,omitempty"`
	EnrollMFAstart    *EnrollMFAStartData    `json:"enroll_mfa_start,omitempty"`
	EnrollMFAComplete *EnrollMFACompleteData `json:"enroll_mfa_complete,omitempty"`
	Signup            *SignupData            `json:"signup,omitempty"`
	VerifyCaptcha     *VerifyCaptchaData     `json:"verify_captcha,omitempty"`
	VerifyEmail       *map[string]any        `json:"verify_email,omitempty"`
	VerifyMFAStart    *VerifyMFAStartData    `json:"verify_mfa_start,omitempty"`
	VerifyMFAComplete *map[string]any        `json:"verify_mfa_complete,omitempty"`
	VerifyPassword    *VerifyPasswordData    `json:"verify_password,omitempty"`
	VerifySocial      *VerifySocialData      `json:"verify_social,omitempty"`
}

type FlowEnrollMFACompleteResult struct {
	CommonFlowResult
}

func (a *FlowEnrollMFA) Complete(ctx context.Context, input FlowEnrollMFACompleteRequest) (*pangea.PangeaResponse[FlowEnrollMFACompleteResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/enroll/mfa/complete", input)
	if err != nil {
		return nil, err
	}

	var out FlowEnrollMFACompleteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowEnrollMFACompleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type FlowResetPasswordRequest struct {
	FlowID   string `json:"flow_id"`
	Password string `json:"password"`
	CBState  string `json:"cb_state,omitempty"`
	CBCode   string `json:"cb_code,omitempty"`
}

type FlowResetPasswordResult struct {
	CommonFlowResult
}

func (a *FlowReset) Password(ctx context.Context, input FlowResetPasswordRequest) (*pangea.PangeaResponse[FlowResetPasswordResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/reset/password", input)
	if err != nil {
		return nil, err
	}

	var out FlowResetPasswordResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowResetPasswordResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/enroll/mfa/start
// # https://dev.pangea.cloud/docs/api/authn#start-the-process-of-enrolling-an-mfa
type FlowEnrollMFAStartRequest struct {
	FlowID      string      `json:"flow_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Phone       string      `json:"phone,omitempty"`
}

type FlowEnrollMFAStartResult struct {
	CommonFlowResult
}

func (a *FlowEnrollMFA) Start(ctx context.Context, input FlowEnrollMFAStartRequest) (*pangea.PangeaResponse[FlowEnrollMFAStartResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/enroll/mfa/start", input)
	if err != nil {
		return nil, err
	}

	var out FlowEnrollMFAStartResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowEnrollMFAStartResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/signup/password
// # https://dev.pangea.cloud/docs/api/authn#signup-a-new-account-using-a-password
type FlowSignupPasswordRequest struct {
	FlowID    string `json:"flow_id"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type FlowSignupPasswordResult struct {
	CommonFlowResult
}

func (a *FlowSignup) Password(ctx context.Context, input FlowSignupPasswordRequest) (*pangea.PangeaResponse[FlowSignupPasswordResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/signup/password", input)
	if err != nil {
		return nil, err
	}

	var out FlowSignupPasswordResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowSignupPasswordResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/signup/social
// # https://dev.pangea.cloud/docs/api/authn#signup-a-new-account-using-a-social-provider
type FlowSignupSocialRequest struct {
	FlowID  string `json:"flow_id"`
	CBState string `json:"cb_state"`
	CBCode  string `json:"cb_code"`
}

type FlowSignupSocialResult struct {
	CommonFlowResult
}

func (a *FlowSignup) Social(ctx context.Context, input FlowSignupSocialRequest) (*pangea.PangeaResponse[FlowSignupSocialResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/signup/social", input)
	if err != nil {
		return nil, err
	}

	var out FlowSignupSocialResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowSignupSocialResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/start
// # https://dev.pangea.cloud/docs/api/authn#start-a-new-signup-or-signin-flow
type FlowStartRequest struct {
	CBURI     string      `json:"cb_uri"`
	Email     *string     `json:"email,omitempty"`
	FlowTypes *[]string   `json:"flow_types,omitempty"`
	Provider  *IDProvider `json:"provider,omitempty"`
}

type FlowStartResult struct {
	CommonFlowResult
}

func (a *Flow) Start(ctx context.Context, input FlowStartRequest) (*pangea.PangeaResponse[FlowStartResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/start", input)
	if err != nil {
		return nil, err
	}

	var out FlowStartResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowStartResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/verify/captcha
// # https://dev.pangea.cloud/docs/api/authn#verify-a-captcha-during-a-signup-or-signin-flow
type FlowVerifyCaptchaRequest struct {
	FlowID string `json:"flow_id"`
	Code   string `json:"code"`
}

type FlowVerifyCaptchaResult struct {
	CommonFlowResult
}

func (a *FlowVerify) Captcha(ctx context.Context, input FlowVerifyCaptchaRequest) (*pangea.PangeaResponse[FlowVerifyCaptchaResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/verify/captcha", input)
	if err != nil {
		return nil, err
	}

	var out FlowVerifyCaptchaResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowVerifyCaptchaResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/verify/email
// # https://dev.pangea.cloud/docs/api/authn#verify-an-email-address-during-a-signup-or-signin-flow
type FlowVerifyEmailRequest struct {
	FlowID  string `json:"flow_id"`
	CBState string `json:"cb_state"`
	CBCode  string `json:"cb_code"`
}

type FlowVerifyEmailResult struct {
	CommonFlowResult
}

func (a *FlowVerify) Email(ctx context.Context, input FlowVerifyEmailRequest) (*pangea.PangeaResponse[FlowVerifyEmailResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/verify/email", input)
	if err != nil {
		return nil, err
	}

	var out FlowVerifyEmailResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowVerifyEmailResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/verify/mfa/complete
// # https://dev.pangea.cloud/docs/api/authn#complete-mfa-verification
type FlowVerifyMFACompleteRequest struct {
	FlowID string  `json:"flow_id"`
	Code   *string `json:"code,omitempty"`
	Cancel *bool   `json:"cancel,omitempty"`
}

type FlowVerifyMFACompleteResult struct {
	CommonFlowResult
}

func (a *FlowVerifyMFA) Complete(ctx context.Context, input FlowVerifyMFACompleteRequest) (*pangea.PangeaResponse[FlowVerifyMFACompleteResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/verify/mfa/captcha", input)
	if err != nil {
		return nil, err
	}

	var out FlowVerifyMFACompleteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowVerifyMFACompleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/verify/mfa/start
// # https://dev.pangea.cloud/docs/api/authn#start-the-process-of-mfa-verification
type FlowVerifyMFAStartRequest struct {
	FlowID      string      `json:"flow_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
}

type FlowVerifyMFAStartResult struct {
	CommonFlowResult
}

func (a *FlowVerifyMFA) Start(ctx context.Context, input FlowVerifyMFAStartRequest) (*pangea.PangeaResponse[FlowVerifyMFAStartResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/verify/mfa/start", input)
	if err != nil {
		return nil, err
	}

	var out FlowVerifyMFAStartResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowVerifyMFAStartResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/verify/password
// # https://dev.pangea.cloud/docs/api/authn#sign-in-with-a-password
type FlowVerifyPasswordRequest struct {
	FlowID   string  `json:"flow_id"`
	Password *string `json:"password,omitempty"`
	Cancel   *bool   `json:"cancel,omitempty"`
}

type FlowVerifyPasswordResult struct {
	CommonFlowResult
}

func (a *FlowVerify) Password(ctx context.Context, input FlowVerifyPasswordRequest) (*pangea.PangeaResponse[FlowVerifyPasswordResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/verify/password", input)
	if err != nil {
		return nil, err
	}

	var out FlowVerifyPasswordResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowVerifyPasswordResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/flow/verify/social
// # https://dev.pangea.cloud/docs/api/authn#signin-with-a-social-provider
type FlowVerifySocialRequest struct {
	FlowID  string `json:"flow_id"`
	CBState string `json:"cb_state"`
	CBCode  string `json:"cb_code"`
}

type FlowVerifySocialResult struct {
	CommonFlowResult
}

func (a *FlowVerify) Social(ctx context.Context, input FlowVerifySocialRequest) (*pangea.PangeaResponse[FlowVerifySocialResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/flow/verify/social", input)
	if err != nil {
		return nil, err
	}

	var out FlowVerifySocialResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FlowVerifySocialResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/user/mfa/delete
// # https://dev.pangea.cloud/docs/api/authn#delete-mfa-enrollment-for-a-user
type UserMFADeleteRequest struct {
	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
}

type UserMFADeleteResult struct {
}

func (a *UserMFA) Delete(ctx context.Context, input UserMFADeleteRequest) (*pangea.PangeaResponse[UserMFADeleteResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/mfa/delete", input)
	if err != nil {
		return nil, err
	}

	var out UserMFADeleteResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserMFADeleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/user/mfa/enroll
// # https://dev.pangea.cloud/docs/api/authn#enroll-mfa-for-a-user
type UserMFAEnrollRequest struct {
	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Code        string      `json:"code"`
}

type UserMFAEnrollResult struct {
}

func (a *UserMFA) Enroll(ctx context.Context, input UserMFAEnrollRequest) (*pangea.PangeaResponse[UserMFAEnrollResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/mfa/enroll", input)
	if err != nil {
		return nil, err
	}

	var out UserMFAEnrollResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserMFAEnrollResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/user/mfa/start
// # https://dev.pangea.cloud/docs/api/authn#start-mfa-verification-for-a-user
type UserMFAStartRequest struct {
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
	TOTPSecret UserMFAStartTOTPSecret `json:"totp_secret"`
}

func (a *UserMFA) Start(ctx context.Context, input UserMFAStartRequest) (*pangea.PangeaResponse[UserMFAStartResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/mfa/start", input)
	if err != nil {
		return nil, err
	}

	var out UserMFAStartResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserMFAStartResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/user/mfa/verify
// # https://dev.pangea.cloud/docs/api/authn#verify-an-mfa-code
type UserMFAVerifyRequest struct {
	UserID      string      `json:"user_id"`
	MFAProvider MFAProvider `json:"mfa_provider"`
	Code        string      `json:"code"`
}

type UserMFAVerifyResult struct {
}

func (a *UserMFA) Verify(ctx context.Context, input UserMFAVerifyRequest) (*pangea.PangeaResponse[UserMFAVerifyResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/mfa/verify", input)
	if err != nil {
		return nil, err
	}

	var out UserMFAVerifyResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserMFAVerifyResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// #   - path: authn::/v1/user/verify
// # https://dev.pangea.cloud/docs/api/authn#verify-a-user
type UserVerifyRequest struct {
	IDProvider    IDProvider `json:"id_provider"`
	Email         string     `json:"email"`
	Authenticator string     `json:"authenticator"`
}

type UserVerifyResult struct {
	Identity     string            `json:"identity"`
	Email        string            `json:"email"`
	Profile      map[string]string `json:"profile"`
	Scopes       *[]string         `json:"scopes"`
	IDProvider   string            `json:"id_provider"`
	MFAProviders []string          `json:"mfa_providers"`
	RequireMFA   bool              `json:"require_mfa"`
	Verified     bool              `json:"verified"`
	Disable      bool              `json:"disable"`
	LastLoginAt  *string           `json:"last_login_at,omitempty"`
}

func (a *User) Verify(ctx context.Context, input UserVerifyRequest) (*pangea.PangeaResponse[UserVerifyResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/verify", input)
	if err != nil {
		return nil, err
	}

	var out UserVerifyResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserVerifyResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type ClientSessionInvalidateRequest struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
}

type ClientSessionInvalidateResult struct {
}

func (a *ClientSession) Invalidate(ctx context.Context, input ClientSessionInvalidateRequest) (*pangea.PangeaResponse[ClientSessionInvalidateResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/client/session/invalidate", input)
	if err != nil {
		return nil, err
	}

	var out ClientSessionInvalidateResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[ClientSessionInvalidateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
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

type SessionListOrder string

const (
	SLOasc  SessionListOrder = "asc"
	SLOdesc                  = "desc"
)

type ClientSessionListRequest struct {
	Token   string             `json:"token"`
	Filter  map[string]string  `json:"filter,omitempty"`
	Last    string             `json:"last,omitempty"`
	Order   SessionListOrder   `json:"order,omitempty"`
	OrderBy SessionListOrderBy `json:"order_by,omitempty"`
}

type SessionToken struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Life      int               `json:"list"`
	Expire    string            `json:"expire"`
	Email     string            `json:"email"`
	Scopes    []string          `json:"scopes"`
	Profile   map[string]string `json:"profile"`
	CreatedAt string            `json:"created_at"`
}

type SessionItem struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Life        int               `json:"list"`
	Expire      string            `json:"expire"`
	Identity    string            `json:"identity"`
	Email       string            `json:"email"`
	Scopes      []string          `json:"scopes"`
	Profile     map[string]string `json:"profile"`
	CreatedAt   string            `json:"created_at"`
	ActiveToken SessionToken      `json:"active_token"`
	Last        string            `json:"last"`
}

type SessionListResult struct {
	Sessions []SessionItem `json:"sessions"`
}

func (a *ClientSession) List(ctx context.Context, input ClientSessionListRequest) (*pangea.PangeaResponse[SessionListResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/client/session/invalidate", input)
	if err != nil {
		return nil, err
	}

	var out SessionListResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SessionListResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type ClientSessionLogoutRequest struct {
	Token string `json:"token"`
}

type ClientSessionLogoutResult struct {
}

func (a *ClientSession) Logout(ctx context.Context, input ClientSessionLogoutRequest) (*pangea.PangeaResponse[ClientSessionLogoutResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/client/session/logout", input)
	if err != nil {
		return nil, err
	}

	var out ClientSessionLogoutResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[ClientSessionLogoutResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type ClientSessionRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	ActiveToken  string `json:"active_token"`
}

type ClientSessionRefreshResult struct {
	RefreshToken LoginToken `json:"refresh_token"`
	ActiveToken  LoginToken `json:"active_token"`
}

func (a *ClientSession) Refresh(ctx context.Context, input ClientSessionRefreshRequest) (*pangea.PangeaResponse[ClientSessionRefreshResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/client/session/refresh", input)
	if err != nil {
		return nil, err
	}

	var out ClientSessionRefreshResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[ClientSessionRefreshResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type SessionListRequest struct {
	Filter  map[string]string  `json:"filter,omitempty"`
	Last    string             `json:"last,omitempty"`
	Order   SessionListOrder   `json:"order,omitempty"`
	OrderBy SessionListOrderBy `json:"order_by,omitempty"`
}

func (a *Session) List(ctx context.Context, input SessionListRequest) (*pangea.PangeaResponse[SessionListResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/session/list", input)
	if err != nil {
		return nil, err
	}

	var out SessionListResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SessionListResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type SessionInvalidateRequest struct {
	SessionID string `json:"session_id"`
}

type SessionInvalidateResult struct {
}

func (a *Session) Invalidate(ctx context.Context, input SessionInvalidateRequest) (*pangea.PangeaResponse[SessionInvalidateResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/session/invalidate", input)
	if err != nil {
		return nil, err
	}

	var out SessionInvalidateResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SessionInvalidateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type SessionLogoutRequest struct {
	UserID string `json:"user_id"`
}

type SessionLogoutResult struct {
}

func (a *Session) Logout(ctx context.Context, input SessionLogoutRequest) (*pangea.PangeaResponse[SessionLogoutResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/session/invalidate", input)
	if err != nil {
		return nil, err
	}

	var out SessionLogoutResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SessionLogoutResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}
