package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"

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

	ID      string  `json:"id"`                // The item ID
	Version *string `json:"version,omitempty"` // The key version(s). all for all versions, num for a specific version, -num for the num latest versions
}

type JWKGetResult struct {
	Keys []JWT `json:"keys"` // The JSON Web Key Set (JWKS) object. Fields with key information are base64URL encoded.
}

type JWTSignRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID      string `json:"id"`      // The item ID
	Payload string `json:"payload"` // The JWT payload (in JSON)
}

type JWTSignResult struct {
	JWS string `json:"jws"` // The signed JSON Web Token (JWS)
}

type JWTVerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	JWS string `json:"jws"` // The signed JSON Web Token (JWS)
}

type JWTVerifyResult struct {
	ValidSignature bool `json:"valid_signature"` // Indicates if messages have been verified.
}
