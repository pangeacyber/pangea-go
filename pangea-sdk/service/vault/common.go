package vault

// EncodedPublicKey is a PEM public key, with no further encoding (i.e. no base64)
// It may be used for example in openssh with no further processing
type EncodedPublicKey string

// EncodedPrivateKey is a PEM private key, with no further encoding (i.e. no base64).
// It may be used for example in openssh with no further processing
type EncodedPrivateKey string

// EncodedSymmetricKey is a base64 encoded key
type EncodedSymmetricKey string

type KeyPurpose string

const (
	KPsigning    KeyPurpose = "signing"
	KPencryption KeyPurpose = "encryption"
	KPjwt        KeyPurpose = "jwt"
)

type AsymmetricAlgorithm string

const (
	AAed25519 AsymmetricAlgorithm = "ed25519"
	AArsa     AsymmetricAlgorithm = "rsa"
	AAes256   AsymmetricAlgorithm = "es256"
	AAes384   AsymmetricAlgorithm = "es384"
	AAes512   AsymmetricAlgorithm = "es512"
)

type SymmetricAlgorithm string

const (
	SYAaes   SymmetricAlgorithm = "aes"
	SYAhs256 SymmetricAlgorithm = "hs256"
	SYAhs384 SymmetricAlgorithm = "hs384"
	SYAhs512 SymmetricAlgorithm = "hs512"
)

type ItemVersionState string

const (
	IVSactive      ItemVersionState = "active"
	IVSdeactivated ItemVersionState = "deactivated"
	IVSsuspended   ItemVersionState = "suspended"
	IVScompromised ItemVersionState = "compromised"
	IVSdestroyed   ItemVersionState = "destroyed"
)

type ItemState string

const (
	ISenabled  ItemState = "enabled"
	ISdisabled ItemState = "disabled"
)

type Metadata map[string]string
type Tags []string

type ItemType string

const (
	ITasymmetricKey ItemType = "asymmetric_key"
	ITsymmetricKey  ItemType = "symmetric_key"
	ITsecret        ItemType = "secret"
	ITpangeaToken   ItemType = "pangea_token"
)

type ItemOrder string

const (
	IOasc ItemOrder = "asc"
	IOdes ItemOrder = "desc"
)

type ItemOrderBy string

const (
	IOBtype         ItemOrderBy = "type"
	IOBcreateAt     ItemOrderBy = "create_at"
	IOBdestroyedAt  ItemOrderBy = "destroyed_at"
	IOBidentity     ItemOrderBy = "identity"
	IOBmanaged      ItemOrderBy = "managed"
	IOBpurpose      ItemOrderBy = "purpose"
	IOBexpiration   ItemOrderBy = "expiration"
	IOBlastRotated  ItemOrderBy = "last_rotated"
	IOBnextRotation ItemOrderBy = "next_rotation"
	IOBname         ItemOrderBy = "name"
	IOBfolder       ItemOrderBy = "folder"
	IOBversion      ItemOrderBy = "version"
)

type CommonStoreRequest struct {
	Type              ItemType `json:"type"`
	Name              string   `json:"name,omitempty"`
	Folder            string   `json:"folder,omitempty"`
	Metadata          Metadata `json:"metadata,omitempty"`
	Tags              Tags     `json:"tags,omitempty"`
	RotationFrequency string   `json:"rotation_frequency,omitempty"`
	RotationState     string   `json:"rotation_state,omitempty"`
	Expiration        string   `json:"expiration,omitempty"`
}

type CommonStoreResult struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`
}

type CommonGenerateRequest struct {
	Type              ItemType `json:"type"`
	Name              string   `json:"name,omitempty"`
	Folder            string   `json:"folder,omitempty"`
	Metadata          Metadata `json:"metadata,omitempty"`
	Tags              Tags     `json:"tags,omitempty"`
	RotationFrequency string   `json:"rotation_frequency,omitempty"`
	RotationState     string   `json:"rotation_state,omitempty"`
	Expiration        string   `json:"expiration,omitempty"`
}

type CommonGenerateResult struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`
}

type GetRequest struct {
	ID           string            `json:"id"`
	Version      string            `json:"version,omitempty"`
	Verbose      *bool             `json:"verbose,omitempty"`
	VersionState *ItemVersionState `json:"version_state,omitempty"`
}

type ItemVersionData struct {
	Version   int               `json:"version"`
	State     string            `json:"state"`
	CreatedAt string            `json:"created_at"`
	DestroyAt *string           `json:"destroy_at,omitempty"`
	PublicKey *EncodedPublicKey `json:"public_key,omitempty"`
	Secret    *string           `json:"secret,omitempty"`
}

type ItemData struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	ItemState         string          `json:"item_state"`
	CurrentVersion    ItemVersionData `json:"current_version"`
	Name              string          `json:"name,omitempty"`
	Folder            string          `json:"folder,omitempty"`
	Metadata          Metadata        `json:"metadata,omitempty"`
	Tags              Tags            `json:"tags,omitempty"`
	RotationFrequency string          `json:"rotation_frequency,omitempty"`
	RotationState     string          `json:"rotation_state,omitempty"`
	LastRotated       string          `json:"last_rotated,omitempty"`
	NextRotation      string          `json:"next_rotation,omitempty"`
	Expiration        string          `json:"expiration,omitempty"`
	CreatedAt         string          `json:"created_at"`
	Algorithm         string          `json:"algorithm,omitempty"`
	Purpose           string          `json:"purpose,omitempty"`
}

type GetResult struct {
	ItemData
	Versions            []ItemVersionData `json:"versions"`
	RotationGracePeriod string            `json:"rotation_grace_period,omitempty"`
}

type ListItemData struct {
	ItemData
	CompromisedVersions []ItemVersionData `json:"compromised_versions"`
}

type ListResult struct {
	Items []ListItemData `json:"items"`
	Count int            `json:"count"`
	Last  string         `json:"last,omitempty"`
}

type ListRequest struct {
	Filter  map[string]string `json:"filter,omitempty"`
	Last    string            `json:"last,omitempty"`
	Size    int               `json:"size,omitempty"`
	Order   ItemOrder         `json:"order,omitempty"`
	OrderBy ItemOrderBy       `json:"order_by,omitempty"`
}

type CommonRotateRequest struct {
	ID            string           `json:"id"`
	RotationState ItemVersionState `json:"rotation_state,omitempty"`
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
	PublicKey *EncodedPublicKey `json:"public_key,omitempty"`
	Algorithm string            `json:"algorithm"`
	Purpose   string            `json:"purpose"`
}

type StateChangeRequest struct {
	ID            string           `json:"id"`
	State         ItemVersionState `json:"state"`
	Version       *int             `json:"version,omitempty"`
	DestroyPeriod string           `json:"destroy_period,omitempty"`
}

type StateChangeResult struct {
	ID        string  `json:"id"`
	Version   int     `json:"version"`
	State     string  `json:"state"`
	DestroyAt *string `json:"destroy_at,omitempty"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type DeleteResult struct {
	ID string `json:"id"`
}

type UpdateRequest struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name,omitempty"`
	Folder              string    `json:"folder,omitempty"`
	Metadata            Metadata  `json:"metadata,omitempty"`
	Tags                Tags      `json:"tags,omitempty"`
	RotationFrequency   string    `json:"rotation_frequency,omitempty"`
	RotationState       string    `json:"rotation_state,omitempty"`
	RotationGracePeriod string    `json:"rotation_grace_period,omitempty"`
	Expiration          string    `json:"expiration,omitempty"`
	ItemState           ItemState `json:"item_state,omitempty"`
}

type UpdateResult struct {
	ID string `json:"id"`
}
