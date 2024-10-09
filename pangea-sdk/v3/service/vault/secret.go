package vault

type SecretStoreRequest struct {
	CommonStoreRequest
	Secret              string `json:"secret,omitempty"`                // The secret value
	Token               string `json:"token,omitempty"`                 // The Pangea Token value
	ClientSecret        string `json:"client_secret,omitempty"`         // The oauth client secret
	ClientID            string `json:"client_id,omitempty"`             // The oauth client ID
	ClientSecretID      string `json:"client_secret_id,omitempty"`      // The oauth client secret ID
	RotationGracePeriod string `json:"rotation_grace_period,omitempty"` // Grace period for the previous version of the secret
}

type SecretStoreResult struct {
	ItemData
	Secret              string `json:"secret,omitempty"`                // The secret value
	Token               string `json:"token,omitempty"`                 // The Pangea Token value
	ClientSecret        string `json:"client_secret,omitempty"`         // The oauth client secret
	ClientID            string `json:"client_id,omitempty"`             // The oauth client ID
	ClientSecretID      string `json:"client_secret_id,omitempty"`      // The oauth client secret ID
	RotationGracePeriod string `json:"rotation_grace_period,omitempty"` // Grace period for the previous version of the secret
}

type SecretRotateRequest struct {
	CommonRotateRequest
	RotationGracePeriod string `json:"rotation_grace_period,omitempty"` // Grace period for the previous version of the secret
	Secret              string `json:"secret,omitempty"`
}

type SecretRotateResult struct {
	ItemData
}
