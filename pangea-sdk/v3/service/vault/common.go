package vault

import "github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"

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
	AAed25519                 AsymmetricAlgorithm = "ED25519"
	AAes256                   AsymmetricAlgorithm = "ES256"
	AAes384                   AsymmetricAlgorithm = "ES384"
	AAes512                   AsymmetricAlgorithm = "ES512"
	AArsa2048_pkcs1v15_sha256 AsymmetricAlgorithm = "RSA-PKCS1V15-2048-SHA256"
	AArsa2048_oaep_sha256     AsymmetricAlgorithm = "RSA-OAEP-2048-SHA256"
	AAes256K                  AsymmetricAlgorithm = "ES256K"
	AArsa2048_oaep_sha1       AsymmetricAlgorithm = "RSA-OAEP-2048-SHA1"
	AArsa2048_oaep_sha512     AsymmetricAlgorithm = "RSA-OAEP-2048-SHA512"
	AArsa3072_oaep_sha1       AsymmetricAlgorithm = "RSA-OAEP-3072-SHA1"
	AArsa3072_oaep_sha256     AsymmetricAlgorithm = "RSA-OAEP-3072-SHA256"
	AArsa3072_oaep_sha512     AsymmetricAlgorithm = "RSA-OAEP-3072-SHA512"
	AArsa4096_oaep_sha1       AsymmetricAlgorithm = "RSA-OAEP-4096-SHA1"
	AArsa4096_oaep_sha256     AsymmetricAlgorithm = "RSA-OAEP-4096-SHA256"
	AArsa4096_oaep_sha512     AsymmetricAlgorithm = "RSA-OAEP-4096-SHA512"
	AArsa2048_pss_sha256      AsymmetricAlgorithm = "RSA-PSS-2048-SHA256"
	AArsa3072_pss_sha256      AsymmetricAlgorithm = "RSA-PSS-3072-SHA256"
	AA4096_pss_sha256         AsymmetricAlgorithm = "RSA-PSS-4096-SHA256"
	AArsa4096_pss_sha512      AsymmetricAlgorithm = "RSA-PSS-4096-SHA512"
	AArsa                     AsymmetricAlgorithm = "RSA-PKCS1V15-2048-SHA256" // deprecated, use AArsa2048_pkcs1v15_sha256 instead
)

type SymmetricAlgorithm string

const (
	SYAhs256      SymmetricAlgorithm = "HS256"
	SYAhs384      SymmetricAlgorithm = "HS384"
	SYAhs512      SymmetricAlgorithm = "HS512"
	SYAaes128_cfb SymmetricAlgorithm = "AES-CFB-128"
	SYAaes256_cfb SymmetricAlgorithm = "AES-CFB-256"
	SYAaes256_gcm SymmetricAlgorithm = "AES-GCM-256"
	SYAaes        SymmetricAlgorithm = "AES-CFB-128" // deprecated, use SYAaes128_cfb instead
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
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Type              ItemType         `json:"type"`
	Name              string           `json:"name,omitempty"`
	Folder            string           `json:"folder,omitempty"`
	Metadata          Metadata         `json:"metadata,omitempty"`
	Tags              Tags             `json:"tags,omitempty"`
	RotationFrequency string           `json:"rotation_frequency,omitempty"`
	RotationState     ItemVersionState `json:"rotation_state,omitempty"`
	Expiration        string           `json:"expiration,omitempty"`
}

type CommonStoreResult struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`
}

type CommonGenerateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Type              ItemType         `json:"type"`
	Name              string           `json:"name,omitempty"`
	Folder            string           `json:"folder,omitempty"`
	Metadata          Metadata         `json:"metadata,omitempty"`
	Tags              Tags             `json:"tags,omitempty"`
	RotationFrequency string           `json:"rotation_frequency,omitempty"`
	RotationState     ItemVersionState `json:"rotation_state,omitempty"`
	Expiration        string           `json:"expiration,omitempty"`
}

type CommonGenerateResult struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Version int    `json:"version"`
}

type GetRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

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
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Filter  map[string]string `json:"filter,omitempty"`
	Last    string            `json:"last,omitempty"`
	Size    int               `json:"size,omitempty"`
	Order   ItemOrder         `json:"order,omitempty"`
	OrderBy ItemOrderBy       `json:"order_by,omitempty"`
}

type CommonRotateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

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
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

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
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID string `json:"id"`
}

type DeleteResult struct {
	ID string `json:"id"`
}

type UpdateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID                  string           `json:"id"`
	Name                string           `json:"name,omitempty"`
	Folder              string           `json:"folder,omitempty"`
	Metadata            Metadata         `json:"metadata,omitempty"`
	Tags                Tags             `json:"tags,omitempty"`
	RotationFrequency   string           `json:"rotation_frequency,omitempty"`
	RotationState       ItemVersionState `json:"rotation_state,omitempty"`
	RotationGracePeriod string           `json:"rotation_grace_period,omitempty"`
	Expiration          string           `json:"expiration,omitempty"`
	ItemState           ItemState        `json:"item_state,omitempty"`
}

type UpdateResult struct {
	ID string `json:"id"`
}

type FolderCreateRequest struct {
	pangea.BaseRequest

	Name     string   `json:"name"`
	Folder   string   `json:"folder"`
	Metadata Metadata `json:"metadata,omitempty"`
	Tags     Tags     `json:"tags,omitempty"`
}

type FolderCreateResult struct {
	ID string `json:"id"`
}
