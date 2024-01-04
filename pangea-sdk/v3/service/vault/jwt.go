package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"

type JWT struct {
	Alg string  `json:"alg"`
	Kid *string `json:"kid,omitempty"`
	Kty string  `json:"kty"`
	Use *string `json:"use,omitempty"`
	Crv *string `json:"crv,omitempty"`
	D   *string `json:"d,omitempty"`
	X   *string `json:"x,omitempty"`
	Y   *string `json:"y,omitempty"`
	N   *string `json:"n,omitempty"`
	E   *string `json:"e,omitempty"`
}

type JWKGetRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID      string  `json:"id"`
	Version *string `json:"version,omitempty"`
}

type JWKGetResult struct {
	Keys []JWT `json:"keys"`
}

type JWTSignRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID      string `json:"id"`
	Payload string `json:"payload"`
}

type JWTSignResult struct {
	JWS string `json:"jws"`
}

type JWTVerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	JWS string `json:"jws"`
}

type JWTVerifyResult struct {
	ValidSignature bool `json:"valid_signature"`
}
