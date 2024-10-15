package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"

type AsymmetricGenerateRequest struct {
	CommonGenerateRequest
	Algorithm  AsymmetricAlgorithm `json:"algorithm,omitempty"`  // The algorithm of the key
	Purpose    KeyPurpose          `json:"purpose,omitempty"`    // The purpose of the key
	Exportable *bool               `json:"exportable,omitempty"` // Whether the key is exportable or not.
}

type AsymmetricGenerateResult struct {
	ItemData
}

type AsymmetricStoreRequest struct {
	CommonStoreRequest
	Algorithm  AsymmetricAlgorithm `json:"algorithm"`            // The algorithm of the key
	PublicKey  EncodedPublicKey    `json:"public_key"`           // The public key (in PEM format)
	PrivateKey EncodedPrivateKey   `json:"private_key"`          // The private key (in PEM format)
	Purpose    KeyPurpose          `json:"purpose,omitempty"`    // The purpose of the key
	Exportable *bool               `json:"exportable,omitempty"` // Whether the key is exportable or not.
}

type AsymmetricStoreResult struct {
	ItemData
}

type SignRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID      string `json:"id"`                // The ID of the item
	Version *int   `json:"version,omitempty"` // The item version
	Message string `json:"message"`           // The message to be signed
}

type SignResult struct {
	ID        string            `json:"id"`                   // The ID of the item
	Version   int               `json:"version"`              // The item version
	Signature string            `json:"signature"`            // The signature of the message
	Algorithm string            `json:"algorithm"`            // The algorithm of the key
	PublicKey *EncodedPublicKey `json:"public_key,omitempty"` // The public key (in PEM format)
}

type VerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID        string `json:"id"`                // The ID of the item
	Version   *int   `json:"version,omitempty"` // The item version
	Message   string `json:"message"`           // A message to be verified
	Signature string `json:"signature"`         // The message signature
}

type VerifyResult struct {
	ID             string `json:"id"`              // The ID of the item
	Version        int    `json:"version"`         // The item version
	Algorithm      string `json:"algorithm"`       // The algorithm of the key
	ValidSignature bool   `json:"valid_signature"` // Indicates if messages have been verified.
}
