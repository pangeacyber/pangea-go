package authn

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Scopes []string

type UserinfoResult struct {
	Token     string            `json:"token"`
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Life      string            `json:"life"`
	Expire    string            `json:"expire"`
	Identity  string            `json:"identity"`
	Email     string            `json:"email"`
	Scopes    *Scopes           `json:"scopes,omitempty"`
	Profile   map[string]string `json:"profile"`
	CreatedAt string            `json:"created_at"`
}

type UserinfoRequest struct {
	Code string `json:"code"`
}

func (a *AuthN) Userinfo(ctx context.Context, input UserinfoRequest) (*pangea.PangeaResponse[UserinfoResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/userinfo", input)
	if err != nil {
		return nil, err
	}

	var out UserinfoResult
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserinfoResult]{
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

type UserCreateRequest struct {
	Email         string   `json:"email"`
	Authenticator string   `json:"authenticator"`
	IDProvider    string   `json:"id_provider"`
	Verified      *bool    `json:"verified,omitempty"`
	RequireMFA    *bool    `json:"require_mfa,omitempty"`
	Profile       *Profile `json:"profile,omitempty"`
	Scopes        *Scopes  `json:"scopes,omitempty"`
}

type UserCreateResult struct {
	Identity    string            `json:"identity"`
	Email       string            `json:"email"`
	Profile     map[string]string `json:"profile"`
	IDProvider  IDProvider        `json:"id_provider"`
	RequireMFA  bool              `json:"require_mfa"`
	Verified    bool              `json:"verified"`
	LastLoginAt *bool             `json:"last_login_at,omitempty"`
	Disabled    *bool             `json:"disabled"`
	MFAProvider *[]string         `json:"mfa_provider,omitempty"`
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
	Identity    string            `json:"identity"`
	Email       string            `json:"email"`
	Profile     map[string]string `json:"profile"`
	Scopes      *Scopes           `json:"scopes,omitempty"`
	IDProvider  IDProvider        `json:"id_provider"`
	MFAProvider *[]string         `json:"mfa_provider"`
	RequireMFA  bool              `json:"require_mfa"`
	Verified    bool              `json:"verified"`
	Disabled    bool              `json:"disabled"`
	LastLoginAt string            `json:"last_login_at"`
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
	Scopes     Scopes `json:"scopes"`
	GlobScopes Scopes `json:"glob_scopes"`
}

type UserInfo struct {
	Profile  Profile `json:"profile"`
	Identity string  `json:"identity"`
	Email    string  `json:"email"`
	Scopes   Scopes  `json:"scopes"`
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
	Email  string  `json:"email"`
	Secret string  `json:"secret"`
	Scopes *Scopes `json:"scopes,omitempty"`
}

type UserLoginResult struct {
	Token     string            `json:"token"`
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Life      int               `json:"life"`
	Expire    string            `json:"expire"`
	Identity  string            `json:"identity"`
	Email     string            `json:"email"`
	Scopes    *Scopes           `json:"scopes,omitempty"`
	Profile   map[string]string `json:"profile"`
	CreatedAt string            `json:"created_at"`
}

func (a *User) Login(ctx context.Context, input UserLoginRequest) (*pangea.PangeaResponse[UserLoginResult], error) {
	req, err := a.Client.NewRequest("POST", "v1/user/login", input)
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
	Identity    string            `json:"identity"`
	Email       string            `json:"email"`
	Profile     map[string]string `json:"profile"`
	IDProvider  IDProvider        `json:"id_provider"`
	MFAProvider []string          `json:"mfa_provider"`
	RequireMFA  bool              `json:"require_mfa"`
	Verified    bool              `json:"verified"`
	LastLoginAt string            `json:"last_login_at"`
	Disabled    *bool             `json:"disabled,omitempty"`
}

func (a *Profile) Get(ctx context.Context, input UserProfileGetRequest) (*pangea.PangeaResponse[UserProfileGetResult], error) {
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
	Profile     map[string]string `json:"profile"`
	Identity    *string           `json:"identity,omitempty"`
	Email       *string           `json:"email,omitempty"`
	RequireMFA  *bool             `json:"require_mfa,omitempty"`
	MFAValue    *string           `json:"mfa_value,omitempty"`
	MFAProvider *[]string         `json:"mfa_provider,omitempty"`
}

type UserProfileUpdateResult struct {
	Identity    string            `json:"identity"`
	Email       string            `json:"email"`
	Profile     map[string]string `json:"profile"`
	IDProvider  IDProvider        `json:"id_provider"`
	MFAProvider []string          `json:"mfa_provider"`
	RequireMFA  bool              `json:"require_mfa"`
	Verified    bool              `json:"verified"`
	LastLoginAt string            `json:"last_login_at"`
	Disabled    *bool             `json:"disabled,omitempty"`
}

func (a *Profile) Update(ctx context.Context, input UserProfileUpdateRequest) (*pangea.PangeaResponse[UserProfileUpdateResult], error) {
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

type UserInvite struct {
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
	Invites []UserInvite `json:"invites"`
}

func (a *Invites) List(ctx context.Context) (*pangea.PangeaResponse[UserInviteListResult], error) {
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

func (a *Invites) Delete(ctx context.Context, input UserInviteDeleteRequest) (*pangea.PangeaResponse[UserInviteDeleteResult], error) {
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
