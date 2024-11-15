package user_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

type HashType string

const (
	HTsha265 HashType = "sha256"
	HTsha1   HashType = "sha1"
	HTsha512 HashType = "sha512"
	HTntlm   HashType = "ntlm"
)

// @summary Look up breached users
//
// @description Determine if an email address, username, phone number, or IP
// address was exposed in a security breach.
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
	return request.DoPost(ctx, e.Client, "v1/user/breached", input, &UserBreachedResult{})
}

// @summary Look up breached users V2
//
// @description Determine if an email address, username, phone number, or IP
// address was exposed in a security breach.
//
// @operationId user_intel_post_v2_user_breached
//
// @example
//
// phoneNumbers := [...]string{"8005550123"}
//
//	input := &user_intel.UserBreachedBulkRequest{
//		PhoneNumbers: phoneNumbers,
//		Raw:          true,
//		Verbose:      true,
//		Provider:     "spycloud",
//	}
//
// out, err := userintel.UserBreachedBulk(ctx, input)
func (e *userIntel) UserBreachedBulk(ctx context.Context, input *UserBreachedBulkRequest) (*pangea.PangeaResponse[UserBreachedBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/user/breached", input, &UserBreachedBulkResult{})
}

type UserBreachedRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	IP          string `json:"ip,omitempty"`
	Start       string `json:"start,omitempty"`
	End         string `json:"end,omitempty"`
	Verbose     *bool  `json:"verbose,omitempty"`
	Raw         *bool  `json:"raw,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Cursor      string `json:"cursor,omitempty"` // A token given in the raw response from SpyCloud. Post this back to paginate results
}

type UserBreachedBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Emails       []string `json:"emails,omitempty"`
	Usernames    []string `json:"usernames,omitempty"`
	PhoneNumbers []string `json:"phone_numbers,omitempty"`
	IPs          []string `json:"ips,omitempty"`
	Domains      []string `json:"domains,omitempty"`
	Start        string   `json:"start,omitempty"`
	End          string   `json:"end,omitempty"`
	Verbose      *bool    `json:"verbose,omitempty"`
	Raw          *bool    `json:"raw,omitempty"`
	Provider     string   `json:"provider,omitempty"`
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

type UserBreachedBulkResult struct {
	Data       map[string]UserBreachedData `json:"data"`
	Parameters interface{}                 `json:"parameters,omitempty"`
	RawData    interface{}                 `json:"raw_data,omitempty"`
}

// @summary Look up breached passwords
//
// @description Determine if a password has been exposed in a security breach
// using a 5 character prefix of the password hash.
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
//	out, err := userintel.PasswordBreached(ctx, input)
func (e *userIntel) PasswordBreached(ctx context.Context, input *UserPasswordBreachedRequest) (*pangea.PangeaResponse[UserPasswordBreachedResult], error) {
	return request.DoPost(ctx, e.Client, "v1/password/breached", input, &UserPasswordBreachedResult{})
}

// @summary Look up breached passwords V2
//
// @description Determine if a password has been exposed in a security breach
// using a 5 character prefix of the password hash.
//
// @operationId user_intel_post_v2_password_breached
//
// @example
//
//	hashPrefixes := [...]string{"5baa6"}
//
//	input := &user_intel.UserPasswordBreachedBulkRequest{
//		HashType:     user_intel.HTsha265,
//		HashPrefixes: hashPrefixes,
//		Raw:          true,
//		Verbose:      true,
//		Provider:     "spycloud",
//	}
//
//	out, err := userintel.PasswordBreachedBulk(ctx, input)
func (e *userIntel) PasswordBreachedBulk(ctx context.Context, input *UserPasswordBreachedBulkRequest) (*pangea.PangeaResponse[UserPasswordBreachedBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/password/breached", input, &UserPasswordBreachedBulkResult{})
}

type UserPasswordBreachedRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	HashType   HashType `json:"hash_type,omitempty"`
	HashPrefix string   `json:"hash_prefix,omitempty"`
	Verbose    *bool    `json:"verbose,omitempty"`
	Raw        *bool    `json:"raw,omitempty"`
	Provider   string   `json:"provider,omitempty"`
}

type UserPasswordBreachedBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	HashType     HashType `json:"hash_type,omitempty"`
	HashPrefixes []string `json:"hash_prefixes,omitempty"`
	Verbose      *bool    `json:"verbose,omitempty"`
	Raw          *bool    `json:"raw,omitempty"`
	Provider     string   `json:"provider,omitempty"`
}

type UserPasswordBreachedData struct {
	FoundInBreach bool `json:"found_in_breach"`
	BreachCount   int  `json:"breach_count,omitempty"`
}

type UserPasswordBreachedResult struct {
	Data       UserPasswordBreachedData `json:"data"`
	Parameters map[string]any           `json:"parameters,omitempty"`
	RawData    map[string]any           `json:"raw_data,omitempty"`
}

type UserPasswordBreachedBulkResult struct {
	Data       map[string]UserPasswordBreachedData `json:"data"`
	Parameters map[string]any                      `json:"parameters,omitempty"`
	RawData    map[string]any                      `json:"raw_data,omitempty"`
}
