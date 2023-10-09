package vault

type SecretStoreRequest struct {
	CommonStoreRequest
	Secret string `json:"secret"`
}

type SecretStoreResult struct {
	CommonStoreResult
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

type PangeaTokenRotateRequest struct {
	CommonRotateRequest
	RotationGracePeriod string `json:"rotation_grace_period"`
}

type PangeaTokenStoreRequest struct {
	CommonStoreRequest
	Token string `json:"secret"`
}
