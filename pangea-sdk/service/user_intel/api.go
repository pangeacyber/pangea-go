package user_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type HashType string

const (
	HTsha265 HashType = "sha256"
	HTsha1   HashType = "sha1"
	HTmd5    HashType = "md5"
)

// @summary Look up breached users
//
// @description Determine if an email address, username, phone number, or IP address was exposed in a security breach.
//
// @operationId user_intel_post_v1_user_breached
//
// @example
//
//	input := &user_intel.UserBreachedRequest{
//		PhoneNumber: "8005550123",
//		Raw:         true,
//		Verbose:     true,
//		Provider:    "spycloud",
//	}
//
// out, err := userintel.UserBreached(ctx, input)
func (e *userIntel) UserBreached(ctx context.Context, input *UserBreachedRequest) (*pangea.PangeaResponse[UserBreachedResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/user/breached", input)
	if err != nil {
		return nil, err
	}
	out := UserBreachedResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserBreachedResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type UserBreachedRequest struct {
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	IP          string `json:"ip,omitempty"`
	Start       string `json:"start,omitempty"`
	End         string `json:"end,omitempty"`
	Verbose     bool   `json:"verbose,omitempty"`
	Raw         bool   `json:"raw,omitempty"`
	Provider    string `json:"provider,omitempty"`
}

type UserBreachedData struct {
	FoundInBreach bool `json:"found_in_breach"`
	BreachCount   int  `json:"breach_count,omitempty"`
}

type UserBreachedResult struct {
	Data       UserBreachedData `json:"data"`
	Parameters interface{}      `json:"parameters,omitempty"`
	RawData    interface{}      `json:"raw_data,omitempty"`
}

// @summary Look up breached passwords
//
// @description Determine if a password has been exposed in a security breach using a 5 character prefix of the password hash.
//
// @operationId user_intel_post_v1_password_breached
//
// @example
//
//	input := &user_intel.UserPasswordBreachedRequest{
//		HashType:   user_intel.HTsha265,
//		HashPrefix: "5baa6",
//		Raw:        true,
//		Verbose:    true,
//		Provider:   "spycloud",
//	}
//
// out, err := userintel.PasswordBreached(ctx, input)
func (e *userIntel) PasswordBreached(ctx context.Context, input *UserPasswordBreachedRequest) (*pangea.PangeaResponse[UserPasswordBreachedResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/password/breached", input)
	if err != nil {
		return nil, err
	}
	out := UserPasswordBreachedResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UserPasswordBreachedResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type UserPasswordBreachedRequest struct {
	HashType   HashType `json:"hash_type,omitempty"`
	HashPrefix string   `json:"hash_prefix,omitempty"`
	Verbose    bool     `json:"verbose,omitempty"`
	Raw        bool     `json:"raw,omitempty"`
	Provider   string   `json:"provider,omitempty"`
}

type UserPasswordBreachedData struct {
	FoundInBreach bool `json:"found_in_breach"`
	BreachCount   int  `json:"breach_count,omitempty"`
}

type UserPasswordBreachedResult struct {
	Data       UserPasswordBreachedData `json:"data"`
	Parameters interface{}              `json:"parameters,omitempty"`
	RawData    interface{}              `json:"raw_data,omitempty"`
}
