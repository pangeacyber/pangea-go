package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"

type SymmetricStoreRequest struct {
	CommonStoreRequest
	Key        EncodedSymmetricKey `json:"key"`
	Algorithm  SymmetricAlgorithm  `json:"algorithm"`
	Purpose    KeyPurpose          `json:"purpose,omitempty"`
	Exportable *bool               `json:"exportable,omitempty"`
}

type SymmetricStoreResult struct {
	CommonStoreResult
	Algorithm string `json:"algorithm"`
	Purpose   string `json:"purpose"`
}

type SymmetricGenerateRequest struct {
	CommonGenerateRequest
	Algorithm  SymmetricAlgorithm `json:"algorithm"`
	Purpose    KeyPurpose         `json:"purpose"`
	Exportable *bool              `json:"exportable,omitempty"`
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

// Parameters for an encrypt/decrypt structured request.
type EncryptStructuredRequest struct {
	pangea.BaseRequest

	// The ID of the key to use. It must be an item of type `symmetric_key` or
	// `asymmetric_key` and purpose `encryption`.
	ID string `json:"id"`

	// Structured data for applying bulk operations.
	StructuredData map[string]interface{} `json:"structured_data"`

	// A filter expression. It must point to string elements of the
	// `StructuredData` field.
	Filter string `json:"filter"`

	// The item version. Defaults to the current version.
	Version *int `json:"version,omitempty"`

	// User provided authentication data.
	AdditionalData *string `json:"additional_data,omitempty"`
}

// Result of an encrypt/decrypt structured request.
type EncryptStructuredResult struct {
	// The ID of the item.
	ID string `json:"id"`

	// The item version.
	Version int `json:"version"`

	// The algorithm of the key.
	Algorithm string `json:"algorithm"`

	// Structured data with filtered fields encrypted/decrypted.
	StructuredData map[string]interface{} `json:"structured_data"`
}

type TransformAlphabet string

const (
	TAalphalower        TransformAlphabet = "alphalower"        // Lowercase alphabet (a-z).
	TAalphanumeric      TransformAlphabet = "alphanumeric"      // Alphanumeric (a-z, A-Z, 0-9).
	TAalphanumericlower TransformAlphabet = "alphanumericlower" // Lowercase alphabet with numbers (a-z, 0-9).
	TAalphanumericupper TransformAlphabet = "alphanumericupper" // Uppercase alphabet with numbers (A-Z, 0-9).
	TAalphaupper        TransformAlphabet = "alphaupper"        // Uppercase alphabet (A-Z).
	TAnumeric           TransformAlphabet = "numeric"           // Numeric (0-9).
)

// Parameters for an encrypt transform request.
type EncryptTransformRequest struct {
	pangea.BaseRequest

	// The ID of the key to use.
	ID string `json:"id"`

	// Message to be encrypted.
	PlainText string `json:"plain_text"`

	// Set of characters to use for format-preserving encryption (FPE).
	Alphabet TransformAlphabet `json:"alphabet"`

	// User provided tweak string. If not provided, a random string will be
	// generated and returned. The user must securely store the tweak source
	// which will be needed to decrypt the data.
	Tweak *string `json:"tweak,omitempty"`

	// The item version. Defaults to the current version.
	Version *int `json:"version,omitempty"`
}

// Result of an encrypt transform request.
type EncryptTransformResult struct {
	// The item ID.
	ID string `json:"id"`

	// The encrypted message.
	CipherText string `json:"cipher_text"`

	// The item version.
	Version int `json:"version"`

	// The algorithm of the key.
	Algorithm string `json:"algorithm"`

	// User provided tweak string. If not provided, a random string will be
	// generated and returned. The user must securely store the tweak source
	// which will be needed to decrypt the data.
	Tweak string `json:"tweak"`

	// Set of characters to use for format-preserving encryption (FPE).
	Alphabet TransformAlphabet `json:"alphabet"`
}

// Parameters for a decrypt transform request.
type DecryptTransformRequest struct {
	pangea.BaseRequest

	// The ID of the key to use.
	ID string `json:"id"`

	// A message encrypted by Vault.
	CipherText string `json:"cipher_text"`

	// User provided tweak string. If not provided, a random string will be
	// generated and returned. The user must securely store the tweak source
	// which will be needed to decrypt the data.
	Tweak string `json:"tweak"`

	// Set of characters to use for format-preserving encryption (FPE).
	Alphabet TransformAlphabet `json:"alphabet"`

	// The item version. Defaults to the current version.
	Version *int `json:"version,omitempty"`
}

// Result of a decrypt transform request.
type DecryptTransformResult struct {
	// The item ID.
	ID string `json:"id"`

	// Decrypted message.
	PlainText string `json:"plain_text"`

	// The item version.
	Version int `json:"version"`

	// The algorithm of the key.
	Algorithm string `json:"algorithm"`
}
