package vault

type SecretStoreRequest struct {
	CommonStoreRequest
	Secret string `json:"secret"`
	Type   string `json:"type"`
}

type SecretStoreResult struct {
	CommonStoreResult
	Secret string `json:"secret"`
}

type SecretGenerateRequest struct {
	CommonGenerateRequest
	Secret string `json:"secret"`
}

type SecretGenerateResult struct {
	CommonGenerateResult
	Secret string `json:"secret"`
}

type SecretRotateRequest struct {
	CommonRotateRequest
	Secret string `json:"secret"`
}

type SecretRotateResult struct {
	CommonRotateResult
	Secret string `json:"secret"`
}
