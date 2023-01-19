package vault

// EncodedPublicKey is a PEM public key, with no further encoding (i.e. no base64)
// It may be used for example in openssh with no further processing
type EncodedPublicKey string

// EncodedPrivateKey is a PEM private key, with no further encoding (i.e. no base64).
// It may be used for example in openssh with no further processing
type EncodedPrivateKey string

// EncodedSymmetricKey is a base64 encoded key
type EncodedSymmetricKey string

type AsymmetricPurpose string

const (
	APsigning    AsymmetricPurpose = "signing"
	APencryption                   = "encryption"
)

type AsymmetricAlgorithm string

const (
	AAed25519 AsymmetricAlgorithm = "ed25519"
	AArsa                         = "rsa"
)

type SymmetricAlgorithm string

const (
	SYAaes SymmetricAlgorithm = "aes"
)

type SecretAlgorithm string

const (
	SAbase32 SecretAlgorithm = "base32"
)

type Metadata map[string]string
type Tags []string

type ItemType string

const (
	ITasymmetricKey ItemType = "asymmetric_key"
	ITsymmetricKey           = "symmetric_key"
	ITsecret                 = "secret"
	ITmasterKey              = "master_key"
)

type CommonStoreRequest struct {
	Type                  ItemType `json:"type"`
	Name                  string   `json:"name,omitempty"`
	Folder                string   `json:"folder,omitempty"`
	Metadata              Metadata `json:"metadata,omitempty"`
	Tags                  Tags     `json:"tags,omitempty"`
	AutoRotate            *bool    `json:"auto_rotate,omitempty"`
	RotationPolicy        string   `json:"rotation_policy,omitempty"`
	RetainPreviousVersion *bool    `json:"retain_previous_version,omitempty"`
	Expiration            string   `json:"expiration,omitempty"` //FIXME: datetime?
	Managed               *bool    `json:"managed,omitempty"`
}

type CommonStoreResult struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`
}

// `json:"name,omitempty"`

type CommonGenerateRequest struct {
	Type                  ItemType `json:"type"`
	Name                  string   `json:"name,omitempty"`
	Folder                string   `json:"folder,omitempty"`
	Metadata              Metadata `json:"metadata,omitempty"`
	Tags                  Tags     `json:"tags,omitempty"`
	AutoRotate            *bool    `json:"auto_rotate,omitempty"`
	RotationPolicy        string   `json:"rotation_policy,omitempty"`
	RetainPreviousVersion *bool    `json:"retain_previous_version,omitempty"`
	Expiration            string   `json:"expiration,omitempty"` //FIXME: datetime?
	Managed               *bool    `json:"managed,omitempty"`
	Store                 *bool    `json:"store,omitempty"`
}

type CommonGenerateResult struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Version *int   `json:"version,omitempty"`
}

type GetRequest struct {
	ID      string `json:"id"`
	Version *int   `json:"version,omitempty"`
	Verbose *bool  `json:"verbose,omitempty"`
}

type CommonGetResult struct {
	ID                    string   `json:"id"`
	Type                  ItemType `json:"type"`
	Version               int      `json:"version"`
	Name                  string   `json:"name,omitempty"`
	Folder                string   `json:"folder,omitempty"`
	Metadata              Metadata `json:"metadata,omitempty"`
	Tags                  Tags     `json:"tags,omitempty"`
	AutoRotate            *bool    `json:"auto_rotate,omitempty"` // FIXME: Should be this required?
	RotationPolicy        string   `json:"rotation_policy,omitempty"`
	RetainPreviousVersion *bool    `json:"retain_previous_version,omitempty"` // FIXME: Should be this required?
	Expiration            string   `json:"expiration,omitempty"`              //FIXME: datetime?
	LastRotated           string   `json:"last_rotated,omitempty"`            //FIXME: datetime?
	NextRotation          string   `json:"next_rotation,omitempty"`           //FIXME: datetime?
	CreatedAt             string   `json:"created_at,omitempty"`              //FIXME: datetime?
	RevokedAt             string   `json:"revoked_at,omitempty"`              //FIXME: datetime?
}

type GetResult struct {
	CommonGetResult
	PublicKey  *EncodedPublicKey    `json:"public_key,omitempty"`
	PrivateKey *EncodedPrivateKey   `json:"private_key,omitempty"`
	Algorithm  string               `json:"algorithm,omitempty"`
	Purpose    *AsymmetricPurpose   `json:"purpose,omitempty"`
	Key        *EncodedSymmetricKey `json:"key,omitempty"`
	Managed    *bool                `json:"managed,omitempty"`
	Secret     *string              `json:"secret,omitempty"`
}

type ListFolderData struct {
	Type   string `json:"type"`
	Name   string `json:"name,omitempty"`
	Folder string `json:"folder,omitempty"`
}

type ListItemData struct {
	ListFolderData
	ID             string   `json:"id"`
	CreatedAt      string   `json:"created_at,omitempty"` //FIXME: datetime?
	RevokedAt      string   `json:"revoked_at,omitempty"` //FIXME: datetime?
	Metadata       Metadata `json:"metadata,omitempty"`
	Tags           Tags     `json:"tags,omitempty"`
	Managed        *bool    `json:"managed,omitempty"`
	NextRotation   string   `json:"next_rotation,omitempty"` //FIXME: datetime?
	Expiration     string   `json:"expiration,omitempty"`    //FIXME: datetime?
	RotationPolicy string   `json:"rotation_policy,omitempty"`
	Version        int      `json:"version"`
	Identity       string   `json:"identity"`
}

type ListResult struct {
	Items []ListItemData `json:"items"`
	Count int            `json:"count"`
	Last  string         `json:"last,omitempty"`
}

type ListRequest struct {
	Filter       map[string]string   `json:"filter,omitempty"`
	Restrictions map[string][]string `json:"restrictions,omitempty"`
	Last         string              `json:"last,omitempty"`
	Size         int                 `json:"size,omitempty"`
	Order        string              `json:"order,omitempty"`
	OrderBy      string              `json:"order_by,omitempty"`
}

type CommonRotateRequest struct {
	ID string `json:"id"`
}

type CommonRotateResult struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Type    string `json:"type"`
}

type KeyRotateRequest struct {
	CommonRotateRequest
	PublicKey  *EncodedPublicKey    `json:"public_key,omitempty"`
	PrivateKey *EncodedPrivateKey   `json:"private_key,omitempty"`
	Key        *EncodedSymmetricKey `json:"key,omitempty"`
}

type KeyRotateResult struct {
	CommonRotateResult
	PublicKey  *EncodedPublicKey    `json:"public_key,omitempty"`
	PrivateKey *EncodedPrivateKey   `json:"private_key,omitempty"`
	Key        *EncodedSymmetricKey `json:"key,omitempty"`
	Algorithm  string               `json:"algorithm"`
}

type RevokeRequest struct {
	ID string `json:"id"`
}

type RevokeResult struct {
	ID string `json:"id"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type DeleteResult struct {
	ID string `json:"id"`
}

type UpdateRequest struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name,omitempty"`
	Folder                string   `json:"folder,omitempty"`
	Metadata              Metadata `json:"metadata,omitempty"`
	Tags                  Tags     `json:"tags,omitempty"`
	AutoRotate            *bool    `json:"auto_rotate,omitempty"` // FIXME: Should be this required?
	RotationPolicy        string   `json:"rotation_policy,omitempty"`
	RetainPreviousVersion *bool    `json:"retain_previous_version,omitempty"` // FIXME: Should be this required?
	Expiration            string   `json:"expiration,omitempty"`              //FIXME: datetime?
}

type UpdateResult struct {
	ID string `json:"id"`
}
