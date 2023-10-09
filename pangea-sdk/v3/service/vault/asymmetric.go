package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"

type AsymmetricGenerateRequest struct {
	CommonGenerateRequest
	Algorithm AsymmetricAlgorithm `json:"algorithm,omitempty"`
	Purpose   KeyPurpose          `json:"purpose,omitempty"`
}

type AsymmetricGenerateResult struct {
	CommonGenerateResult
	PublicKey EncodedPublicKey `json:"public_key"`
	Algorithm string           `json:"algorithm"`
	Purpose   string           `json:"purpose"`
}

type AsymmetricStoreRequest struct {
	CommonStoreRequest
	Algorithm  AsymmetricAlgorithm `json:"algorithm"`
	PublicKey  EncodedPublicKey    `json:"public_key"`
	PrivateKey EncodedPrivateKey   `json:"private_key"`
	Purpose    KeyPurpose          `json:"purpose,omitempty"`
}

type AsymmetricStoreResult struct {
	CommonStoreResult
	PublicKey EncodedPublicKey `json:"public_key"`
	Algorithm string           `json:"algorithm"`
	Purpose   string           `json:"purpose"`
}

type SignRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID      string `json:"id"`
	Message string `json:"message"`
	Version *int   `json:"version,omitempty"`
}

type SignResult struct {
	ID        string            `json:"id"`
	Version   int               `json:"version"`
	Signature string            `json:"signature"`
	Algorithm string            `json:"algorithm"`
	PublicKey *EncodedPublicKey `json:"public_key,omitempty"`
}

type VerifyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID        string `json:"id"`
	Version   *int   `json:"version,omitempty"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type VerifyResult struct {
	ID             string `json:"id"`
	Version        int    `json:"version"`
	Algorithm      string `json:"algorithm"`
	ValidSignature bool   `json:"valid_signature"`
}
