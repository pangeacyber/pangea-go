package vault

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
	ID        string `json:"id"`
	PlainText string `json:"plain_text"`
	Version   *int   `json:"version,omitempty"`
}

type EncryptResult struct {
	ID         string `json:"id"`
	Version    int    `json:"version"`
	Algorithm  string `json:"algorithm"`
	CipherText string `json:"cipher_text"`
}

type DecryptRequest struct {
	ID         string `json:"id"`
	CipherText string `json:"cipher_text"`
	Version    *int   `json:"version,omitempty"`
}

type DecryptResult struct {
	ID        string `json:"id"`
	Version   int    `json:"version"`
	Algorithm string `json:"algorithm"`
	PlainText string `json:"plain_text"`
}
