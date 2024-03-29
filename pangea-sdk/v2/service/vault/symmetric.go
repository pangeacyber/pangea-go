package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"

type SymmetricStoreRequest struct {
	CommonStoreRequest
	Key       EncodedSymmetricKey `json:"key"`
	Algorithm SymmetricAlgorithm  `json:"algorithm"`
	Purpose   KeyPurpose          `json:"purpose,omitempty"`
}

type SymmetricStoreResult struct {
	CommonStoreResult
	Algorithm string `json:"algorithm"`
	Purpose   string `json:"purpose"`
}

type SymmetricGenerateRequest struct {
	CommonGenerateRequest
	Algorithm SymmetricAlgorithm `json:"algorithm"`
	Purpose   KeyPurpose         `json:"purpose"`
}

type SymmetricGenerateResult struct {
	CommonGenerateResult
	Algorithm string `json:"algorithm"`
	Purpose   string `json:"purpose"`
}

type EncryptRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID             string  `json:"id"`
	PlainText      string  `json:"plain_text"`
	Version        *int    `json:"version,omitempty"`
	AdditionalData *string `json:"additional_data,omitempty"`
}

type EncryptResult struct {
	ID         string `json:"id"`
	Version    int    `json:"version"`
	Algorithm  string `json:"algorithm"`
	CipherText string `json:"cipher_text"`
}

type DecryptRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID             string  `json:"id"`
	CipherText     string  `json:"cipher_text"`
	Version        *int    `json:"version,omitempty"`
	AdditionalData *string `json:"additional_data,omitempty"`
}

type DecryptResult struct {
	ID        string `json:"id"`
	Version   int    `json:"version"`
	Algorithm string `json:"algorithm"`
	PlainText string `json:"plain_text"`
}
