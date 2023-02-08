package vault

type AsymmetricGenerateRequest struct {
	CommonGenerateRequest
	Algorithm *AsymmetricAlgorithm `json:"algorithm,omitempty"`
	Purpose   *KeyPurpose          `json:"purpose,omitempty"`
}

type AsymmetricGenerateResult struct {
	CommonGenerateResult
	PublicKey  EncodedPublicKey   `json:"public_key"`
	PrivateKey *EncodedPrivateKey `json:"private_key,omitempty"`
}

type AsymmetricStoreRequest struct {
	CommonStoreRequest
	Algorithm  AsymmetricAlgorithm `json:"algorithm"`
	PublicKey  EncodedPublicKey    `json:"public_key"`
	PrivateKey EncodedPrivateKey   `json:"private_key"`
	Managed    *bool               `json:"managed,omitempty"`
	Purpose    *KeyPurpose         `json:"purpose,omitempty"`
}

type AsymmetricStoreResult struct {
	CommonStoreResult
	Algorithm  AsymmetricAlgorithm `json:"algorithm"`
	PublicKey  EncodedPublicKey    `json:"public_key"`
	PrivateKey *EncodedPrivateKey  `json:"private_key,omitempty"`
}

type SignRequest struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type SignResult struct {
	ID        string              `json:"id"`
	Version   int                 `json:"version"`
	Signature string              `json:"signature"`
	Algorithm AsymmetricAlgorithm `json:"algorithm"`
	PublicKey *EncodedPublicKey   `json:"public_key,omitempty"`
}

type VerifyRequest struct {
	ID        string `json:"id"`
	Version   *int   `json:"version,omitempty"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type VerifyResult struct {
	ID             string              `json:"id"`
	Version        int                 `json:"version"`
	Algorithm      AsymmetricAlgorithm `json:"algorithm"`
	ValidSignature bool                `json:"valid_signature"`
}
