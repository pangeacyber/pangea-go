package vault

type SymmetricStoreRequest struct {
	CommonStoreRequest
	Key       EncodedSymmetricKey `json:"key"`
	Algorithm SymmetricAlgorithm  `json:"algorithm"`
	Managed   *bool               `json:"managed,omitempty"`
}

type SymmetricStoreResult struct {
	CommonStoreResult
	Algorithm *SymmetricAlgorithm  `json:"algorithm,omitempty"`
	Key       *EncodedSymmetricKey `json:"key,omitempty"`
	Managed   *bool                `json:"managed,omitempty"`
	Purpose   *KeyPurpose          `json:"purpose,omitempty"`
}

type SymmetricGenerateRequest struct {
	CommonGenerateRequest
	Algorithm *SymmetricAlgorithm `json:"algorithm,omitempty"`
	Managed   *bool               `json:"managed,omitempty"`
	Purpose   *KeyPurpose         `json:"purpose,omitempty"`
}

type SymmetricGenerateResult struct {
	CommonGenerateResult
	Algorithm SymmetricAlgorithm `json:"algorithm"`
	Key       *EncodedPublicKey  `json:"key,omitempty"`
}

type EncryptRequest struct {
	ID        string `json:"id"`
	PlainText string `json:"plain_text"`
}

type EncryptResult struct {
	ID         string             `json:"id"`
	Version    int                `json:"version"`
	Algorithm  SymmetricAlgorithm `json:"algorithm"`
	CipherText string             `json:"cipher_text"`
}

type DecryptRequest struct {
	ID         string `json:"id"`
	Version    *int   `json:"version,omitempty"`
	CipherText string `json:"cipher_text"`
}

type DecryptResult struct {
	ID        string             `json:"id"`
	Version   *int               `json:"version,omitempty"`
	Algorithm SymmetricAlgorithm `json:"algorithm"`
	PlainText string             `json:"plain_text"`
}
