package vault

import (
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

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
	KPfpe        KeyPurpose = "fpe" // Format-preserving encryption.
)

type AsymmetricAlgorithm string

const (
	AAed25519                               AsymmetricAlgorithm = "ED25519"
	AAes256                                 AsymmetricAlgorithm = "ES256"
	AAes384                                 AsymmetricAlgorithm = "ES384"
	AAes512                                 AsymmetricAlgorithm = "ES512"
	AArsa2048_pkcs1v15_sha256               AsymmetricAlgorithm = "RSA-PKCS1V15-2048-SHA256"
	AArsa2048_oaep_sha256                   AsymmetricAlgorithm = "RSA-OAEP-2048-SHA256"
	AAes256K                                AsymmetricAlgorithm = "ES256K"
	AArsa2048_oaep_sha1                     AsymmetricAlgorithm = "RSA-OAEP-2048-SHA1"
	AArsa2048_oaep_sha512                   AsymmetricAlgorithm = "RSA-OAEP-2048-SHA512"
	AArsa3072_oaep_sha1                     AsymmetricAlgorithm = "RSA-OAEP-3072-SHA1"
	AArsa3072_oaep_sha256                   AsymmetricAlgorithm = "RSA-OAEP-3072-SHA256"
	AArsa3072_oaep_sha512                   AsymmetricAlgorithm = "RSA-OAEP-3072-SHA512"
	AArsa4096_oaep_sha1                     AsymmetricAlgorithm = "RSA-OAEP-4096-SHA1"
	AArsa4096_oaep_sha256                   AsymmetricAlgorithm = "RSA-OAEP-4096-SHA256"
	AArsa4096_oaep_sha512                   AsymmetricAlgorithm = "RSA-OAEP-4096-SHA512"
	AArsa2048_pss_sha256                    AsymmetricAlgorithm = "RSA-PSS-2048-SHA256"
	AArsa3072_pss_sha256                    AsymmetricAlgorithm = "RSA-PSS-3072-SHA256"
	AA4096_pss_sha256                       AsymmetricAlgorithm = "RSA-PSS-4096-SHA256" // deprecated by typo. use AArsa4096_pss_sha256 instead
	AArsa4096_pss_sha256                    AsymmetricAlgorithm = "RSA-PSS-4096-SHA256"
	AArsa4096_pss_sha512                    AsymmetricAlgorithm = "RSA-PSS-4096-SHA512"
	AArsa                                   AsymmetricAlgorithm = "RSA-PKCS1V15-2048-SHA256" // deprecated, use AArsa2048_pkcs1v15_sha256 instead
	AAed25519_dilithium2_beta               AsymmetricAlgorithm = "ED25519-DILITHIUM2-BETA"
	AAed488_dilithium3_beta                 AsymmetricAlgorithm = "ED448-DILITHIUM3-BETA"
	AAsphincsplus_128f_shake256_simple_beta AsymmetricAlgorithm = "SPHINCSPLUS-128F-SHAKE256-SIMPLE-BETA"
	AAsphincsplus_128f_shake256_robust_beta AsymmetricAlgorithm = "SPHINCSPLUS-128F-SHAKE256-ROBUST-BETA"
	AAsphincsplus_192f_shake256_simple_beta AsymmetricAlgorithm = "SPHINCSPLUS-192F-SHAKE256-SIMPLE-BETA"
	AAsphincsplus_192f_shake256_robust_beta AsymmetricAlgorithm = "SPHINCSPLUS-192F-SHAKE256-ROBUST-BETA"
	AAsphincsplus_256f_shake256_simple_beta AsymmetricAlgorithm = "SPHINCSPLUS-256F-SHAKE256-SIMPLE-BETA"
	AAsphincsplus_256f_shake256_robust_beta AsymmetricAlgorithm = "SPHINCSPLUS-256F-SHAKE256-ROBUST-BETA"
	AAsphincsplus_128f_sha256_simple_beta   AsymmetricAlgorithm = "SPHINCSPLUS-128F-SHA256-SIMPLE-BETA"
	AAsphincsplus_128f_sha256_robust_beta   AsymmetricAlgorithm = "SPHINCSPLUS-128F-SHA256-ROBUST-BETA"
	AAsphincsplus_192f_sha256_simple_beta   AsymmetricAlgorithm = "SPHINCSPLUS-192F-SHA256-SIMPLE-BETA"
	AAsphincsplus_192f_sha256_robust_beta   AsymmetricAlgorithm = "SPHINCSPLUS-192F-SHA256-ROBUST-BETA"
	AAsphincsplus_256f_sha256_simple_beta   AsymmetricAlgorithm = "SPHINCSPLUS-256F-SHA256-SIMPLE-BETA"
	AAsphincsplus_256f_sha256_robust_beta   AsymmetricAlgorithm = "SPHINCSPLUS-256F-SHA256-ROBUST-BETA"
	AAfalcon1024_beta                       AsymmetricAlgorithm = "FALCON-1024-BETA"
)

type SymmetricAlgorithm string

const (
	SYAhs256         SymmetricAlgorithm = "HS256"
	SYAhs384         SymmetricAlgorithm = "HS384"
	SYAhs512         SymmetricAlgorithm = "HS512"
	SYAaes128_cfb    SymmetricAlgorithm = "AES-CFB-128"
	SYAaes256_cfb    SymmetricAlgorithm = "AES-CFB-256"
	SYAaes256_gcm    SymmetricAlgorithm = "AES-GCM-256"
	SYAaes128_cbc    SymmetricAlgorithm = "AES-CBC-128"
	SYAaes256_cbc    SymmetricAlgorithm = "AES-CBC-256"
	SYAaes           SymmetricAlgorithm = "AES-CFB-128"        // deprecated, use SYAaes128_cfb instead
	SYAaes_ff3_1_128 SymmetricAlgorithm = "AES-FF3-1-128-BETA" // 128-bit encryption using the FF3-1 algorithm. Beta feature.
	SYAaes_ff3_1_256 SymmetricAlgorithm = "AES-FF3-1-256-BETA" // 256-bit encryption using the FF3-1 algorithm. Beta feature.
)

type ItemVersionState string

const (
	IVSactive      ItemVersionState = "active"
	IVSdeactivated ItemVersionState = "deactivated"
	IVSsuspended   ItemVersionState = "suspended"
	IVScompromised ItemVersionState = "compromised"
	IVSdestroyed   ItemVersionState = "destroyed"
	IVSinherited   ItemVersionState = "inherited"
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
	ITasymmetricKey              ItemType = "asymmetric_key"
	ITsymmetricKey               ItemType = "symmetric_key"
	ITsecret                     ItemType = "secret"
	ITpangeaToken                ItemType = "pangea_token"
	ITfolder                     ItemType = "folder"
	ITpangeaClientSecret         ItemType = "pangea_client_secret"
	ITpangeaPlatformClientSecret ItemType = "pangea_platform_client_secret"
)

type ItemOrder string

const (
	IOasc ItemOrder = "asc"
	IOdes ItemOrder = "desc"
)

type ItemOrderBy string

const (
	IOBid           ItemOrderBy = "id"
	IOBtype         ItemOrderBy = "type"
	IOBcreateAt     ItemOrderBy = "create_at"
	IOBdestroyedAt  ItemOrderBy = "destroyed_at"
	IOBalgorithm    ItemOrderBy = "algorithm"
	IOBpurpose      ItemOrderBy = "purpose"
	IOBdisabledAt   ItemOrderBy = "disabled_at"
	IOBlastRotated  ItemOrderBy = "last_rotated"
	IOBnextRotation ItemOrderBy = "next_rotation"
	IOBname         ItemOrderBy = "name"
	IOBfolder       ItemOrderBy = "folder"
	IOBitemState    ItemOrderBy = "item_state"
)

// Algorithm of an exported public key.
type ExportEncryptionAlgorithm string

const (
	EEArsa4096_oaep_sha512    ExportEncryptionAlgorithm = "RSA-OAEP-4096-SHA512"
	EEArsa4096_no_padding_kem ExportEncryptionAlgorithm = "RSA-NO-PADDING-4096-KEM"
)

type ExportEncryptionType string

const (
	EETasymmetric ExportEncryptionType = "asymmetric"
	EETkem        ExportEncryptionType = "kem"
)

type CommonStoreRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Type              ItemType         `json:"type"`                         // The type of the item
	Name              string           `json:"name,omitempty"`               // The name of this item
	Folder            string           `json:"folder,omitempty"`             // The folder where this item is stored
	Metadata          Metadata         `json:"metadata,omitempty"`           // User-provided metadata
	Tags              Tags             `json:"tags,omitempty"`               // A list of user-defined tags
	RotationFrequency string           `json:"rotation_frequency,omitempty"` // Period of time between item rotations.
	RotationState     ItemVersionState `json:"rotation_state,omitempty"`     // State to which the previous version should transition upon rotation
	DisabledAt        string           `json:"disabled_at,omitempty"`        // Timestamp indicating when the item will be disabled
}

type CommonStoreResult struct {
	ID      string `json:"id"`      // The ID of the item
	Type    string `json:"type"`    // The type of the item
	Version int    `json:"version"` // The item version
}

type CommonGenerateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Type              ItemType `json:"type"`                         // The type of the item
	Name              string   `json:"name,omitempty"`               // The name of this item
	Folder            string   `json:"folder,omitempty"`             // The folder where this item is stored
	Metadata          Metadata `json:"metadata,omitempty"`           // User-provided metadata
	Tags              Tags     `json:"tags,omitempty"`               // A list of user-defined tags
	RotationFrequency string   `json:"rotation_frequency,omitempty"` // Period of time between item rotations.
	RotationState     string   `json:"rotation_state,omitempty"`     // State to which the previous version should transition upon rotation
	DisabledAt        string   `json:"disabled_at,omitempty"`        // Timestamp indicating when the item will be disabled
}

type GetRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID      string `json:"id"`
	Version string `json:"version,omitempty"`
}

type ItemVersionData struct {
	Version        int               `json:"version"`                // The item version
	State          string            `json:"state"`                  // The state of the item version
	CreatedAt      string            `json:"created_at"`             // Timestamp indicating when the item was created
	DestroyedAt    *string           `json:"destroyed_at,omitempty"` // Timestamp indicating when the item version will be destroyed
	RotatedAt      *string           `json:"rotated_at,omitempty"`   // Timestamp indicating when the item version will be rotated
	PublicKey      *EncodedPublicKey `json:"public_key,omitempty"`
	Secret         *string           `json:"secret,omitempty"`
	Token          *string           `json:"token,omitempty"`
	ClientSecret   *string           `json:"client_secret,omitempty"`
	ClientSecretID *string           `json:"client_secret_id,omitempty"`
}

type ItemData struct {
	ID                  string             `json:"id"`                              // The ID of the item
	Type                string             `json:"type"`                            // The type of the item
	NumVersions         int                `json:"num_versions"`                    // Latest version number
	Enabled             bool               `json:"enabled"`                         // True if the item is enabled
	Name                string             `json:"name,omitempty"`                  // The name of this item
	Folder              string             `json:"folder,omitempty"`                // The folder where this item is stored
	Metadata            Metadata           `json:"metadata,omitempty"`              // User-provided metadata
	Tags                Tags               `json:"tags,omitempty"`                  // A list of user-defined tags
	RotationFrequency   string             `json:"rotation_frequency,omitempty"`    // Period of time between item rotations.
	RotationState       string             `json:"rotation_state,omitempty"`        // State to which the previous version should transition upon rotation
	LastRotated         string             `json:"last_rotated,omitempty"`          // Timestamp of the last rotation (if any)
	NextRotation        string             `json:"next_rotation,omitempty"`         // Timestamp of the next rotation, if auto rotation is enabled.
	DisabledAt          string             `json:"disabled_at,omitempty"`           // Timestamp indicating when the item will be disabled
	CreatedAt           string             `json:"created_at"`                      // Timestamp indicating when the item was created
	Algorithm           string             `json:"algorithm,omitempty"`             // The algorithm of the key
	Purpose             string             `json:"purpose,omitempty"`               // The purpose of the key
	RotationGracePeriod string             `json:"rotation_grace_period,omitempty"` // Grace period for the previous version of the secret
	Exportable          *bool              `json:"exportable,omitempty"`            // Whether the key is exportable or not.
	ClientID            string             `json:"client_id,omitempty"`
	InheritedSettings   *InheritedSettings `json:"inherited_settings,omitempty"` // For settings that inherit a value from a parent folder, the full path of the folder where the value is set
	ItemVersions        []ItemVersionData  `json:"item_versions"`
}

type InheritedSettings struct {
	RotationFrequency   string `json:"rotation_frequency,omitempty"`
	RotationState       string `json:"rotation_state,omitempty"`
	RotationGracePeriod string `json:"rotation_grace_period,omitempty"`
}

type GetResult struct {
	ItemData
}

type FilterList struct {
	pangea.FilterBase

	_type        *pangea.FilterMatch[string]
	id           *pangea.FilterMatch[string]
	algorithm    *pangea.FilterMatch[string]
	purpose      *pangea.FilterMatch[string]
	name         *pangea.FilterMatch[string]
	folder       *pangea.FilterMatch[string]
	itemState    *pangea.FilterMatch[string]
	createdAt    *pangea.FilterRange[string]
	destroyedAt  *pangea.FilterRange[string]
	expiration   *pangea.FilterRange[string]
	lastRotated  *pangea.FilterRange[string]
	nextRotation *pangea.FilterRange[string]
}

func NewFilterList() *FilterList {
	filter := make(pangea.Filter)
	return &FilterList{
		FilterBase:   *pangea.NewFilterBase(filter),
		_type:        pangea.NewFilterMatch[string]("type", &filter),
		id:           pangea.NewFilterMatch[string]("id", &filter),
		algorithm:    pangea.NewFilterMatch[string]("algorithm", &filter),
		purpose:      pangea.NewFilterMatch[string]("purpose", &filter),
		name:         pangea.NewFilterMatch[string]("name", &filter),
		folder:       pangea.NewFilterMatch[string]("folder", &filter),
		itemState:    pangea.NewFilterMatch[string]("item_state", &filter),
		createdAt:    pangea.NewFilterRange[string]("created_at", &filter),
		destroyedAt:  pangea.NewFilterRange[string]("destroyed_at", &filter),
		expiration:   pangea.NewFilterRange[string]("expiration", &filter),
		lastRotated:  pangea.NewFilterRange[string]("last_rotated", &filter),
		nextRotation: pangea.NewFilterRange[string]("last_rotation", &filter),
	}
}

func (fu *FilterList) Type() *pangea.FilterMatch[string] {
	return fu._type
}

func (fu *FilterList) ID() *pangea.FilterMatch[string] {
	return fu.id
}

func (fu *FilterList) Algorithm() *pangea.FilterMatch[string] {
	return fu.algorithm
}

func (fu *FilterList) Purpose() *pangea.FilterMatch[string] {
	return fu.purpose
}

func (fu *FilterList) Name() *pangea.FilterMatch[string] {
	return fu.name
}

func (fu *FilterList) Folder() *pangea.FilterMatch[string] {
	return fu.folder
}

func (fu *FilterList) ItemStated() *pangea.FilterMatch[string] {
	return fu.itemState
}

func (fu *FilterList) CreatedAt() *pangea.FilterRange[string] {
	return fu.createdAt
}

func (fu *FilterList) DestroyedAt() *pangea.FilterRange[string] {
	return fu.destroyedAt
}

func (fu *FilterList) Expiration() *pangea.FilterRange[string] {
	return fu.expiration
}

func (fu *FilterList) LastRotated() *pangea.FilterRange[string] {
	return fu.lastRotated
}

func (fu *FilterList) NextRotation() *pangea.FilterRange[string] {
	return fu.nextRotation
}

type GetBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Filter  pangea.Filter `json:"filter,omitempty"`   // A set of filters to help you customize your search
	Last    string        `json:"last,omitempty"`     // Internal ID returned in the previous look up response. Used for pagination.
	Size    int           `json:"size,omitempty"`     // Maximum number of items in the response
	Order   ItemOrder     `json:"order,omitempty"`    // Ordering direction
	OrderBy ItemOrderBy   `json:"order_by,omitempty"` // Property used to order the results
}

type GetBulkResult struct {
	Items []ItemData `json:"items"`
	Last  string     `json:"last,omitempty"`
}

type ListItemData struct {
	ItemData
	CompromisedVersions []ItemVersionData `json:"compromised_versions"`
}

type ListResult struct {
	Items []ListItemData `json:"items"`
	Last  string         `json:"last,omitempty"` // Internal ID returned in the previous look up response. Used for pagination.
}

type ListInclude string

const (
	LIsecrets   ListInclude = "secrets"
	LIencrypted ListInclude = "encrypted"
)

type ListRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Filter  pangea.Filter `json:"filter,omitempty"`   // A set of filters to help you customize your search.
	Last    string        `json:"last,omitempty"`     // Internal ID returned in the previous look up response. Used for pagination.
	Size    int           `json:"size,omitempty"`     // Maximum number of items in the response
	Order   ItemOrder     `json:"order,omitempty"`    // Ordering direction
	OrderBy ItemOrderBy   `json:"order_by,omitempty"` // Property used to order the results
}

type CommonRotateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID            string           `json:"id"`                       // The ID of the key
	RotationState ItemVersionState `json:"rotation_state,omitempty"` // State to which the previous version should transition upon rotation
}

type KeyRotateRequest struct {
	CommonRotateRequest
	PublicKey  *EncodedPublicKey    `json:"public_key,omitempty"`  // The public key (in PEM format)
	PrivateKey *EncodedPrivateKey   `json:"private_key,omitempty"` // The private key (in PEM format)
	Key        *EncodedSymmetricKey `json:"key,omitempty"`         // The key material
}

type KeyRotateResult struct {
	ItemData
}

type StateChangeRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID            string           `json:"id"`                       // The item ID
	State         ItemVersionState `json:"state"`                    // The new state of the item version
	Version       *int             `json:"version,omitempty"`        // The item version
	DestroyPeriod string           `json:"destroy_period,omitempty"` // Period of time for the destruction of a compromised key. Only applicable if state=compromised (format: a positive number followed by a time period (secs, mins, hrs, days, weeks, months, years) or an abbreviation
}

type StateChangeResult struct {
	ItemData
}

type DeleteRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID        string `json:"id"`                  // The item ID
	Recursive *bool  `json:"recursive,omitempty"` // true for recursive deleting all the items inside a folder. Valid only for folders
}

type DeleteResult struct {
	ID string `json:"id"`
}

type UpdateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	ID                  string           `json:"id"`                              // The item ID
	Name                string           `json:"name,omitempty"`                  // The name of this item
	Folder              string           `json:"folder,omitempty"`                // The parent folder where this item is stored
	Metadata            Metadata         `json:"metadata,omitempty"`              // User-provided metadata
	Tags                Tags             `json:"tags,omitempty"`                  // A list of user-defined tags
	DisabledAt          string           `json:"disabled_at,omitempty"`           // Timestamp indicating when the item will be disabled
	Enabled             *bool            `json:"enabled,omitempty"`               // True if the item is enabled
	RotationFrequency   string           `json:"rotation_frequency,omitempty"`    // Period of time between item rotations, never to disable rotation or inherited to inherit the value from the parent folder or from the default settings (format: a positive number followed by a time period (secs, mins, hrs, days, weeks, months, years) or an abbreviation
	RotationState       ItemVersionState `json:"rotation_state,omitempty"`        // State to which the previous version should transition upon rotation or inherited to inherit the value from the parent folder or from the default settings
	RotationGracePeriod string           `json:"rotation_grace_period,omitempty"` // Grace period for the previous version of the Pangea Token or inherited to inherit the value from the parent folder or from the default settings (format: a positive number followed by a time period (secs, mins, hrs, days, weeks, months, years) or an abbreviation
}

type UpdateResult struct {
	ItemData
}

type FolderCreateRequest struct {
	pangea.BaseRequest

	Name                string           `json:"name"`                            // The name of this folder
	Folder              string           `json:"folder,omitempty"`                // The parent folder where this folder is stored
	Metadata            Metadata         `json:"metadata,omitempty"`              // User-provided metadata
	Tags                Tags             `json:"tags,omitempty"`                  // A list of user-defined tags
	RotationFrequency   string           `json:"rotation_frequency,omitempty"`    // Period of time between item rotations, never to disable rotation or inherited to inherit the value from the parent folder or from the default settings (format: a positive number followed by a time period (secs, mins, hrs, days, weeks, months, years) or an abbreviation
	RotationState       ItemVersionState `json:"rotation_state,omitempty"`        // State to which the previous version should transition upon rotation or inherited to inherit the value from the parent folder or from the default settings
	RotationGracePeriod string           `json:"rotation_grace_period,omitempty"` // Grace period for the previous version of the Pangea Token or inherited to inherit the value from the parent folder or from the default settings (format: a positive number followed by a time period (secs, mins, hrs, days, weeks, months, years) or an abbreviation
}

type FolderCreateResult struct {
	ItemData
}

type ExportRequest struct {
	pangea.BaseRequest

	ID                  string                     `json:"id"`                              // The ID of the item.
	Version             *int                       `json:"version,omitempty"`               // The item version.
	AsymmetricPublicKey *string                    `json:"asymmetric_public_key,omitempty"` // Public key in PEM format used to encrypt exported key(s).
	AsymmetricAlgorithm *ExportEncryptionAlgorithm `json:"asymmetric_algorithm,omitempty"`  // The algorithm of the public key.
	KEMPassword         *string                    `json:"kem_password,omitempty"`          // This is the password that will be used along with a salt to derive the symmetric key that is used to encrypt the exported key material. Required if encryption_type is kem.
}

type ExportResult struct {
	ID         string  `json:"id"`                    // The ID of the item.
	Type       string  `json:"type"`                  // The type of the key.
	Version    int     `json:"version"`               // The item version.
	Enabled    bool    `json:"enabled"`               // True if the item is enabled.
	Algorithm  string  `json:"algorithm"`             // The algorithm of the key.
	PublicKey  *string `json:"public_key,omitempty"`  // The public key (in PEM format).
	PrivateKey *string `json:"private_key,omitempty"` // The private key (in PEM format), it could be encrypted or not based on 'encryption_type' value.
	Key        *string `json:"key,omitempty"`         // The key material.

	// Encryption information
	EncryptionType      string `json:"encryption_type"`                // Encryption format of the exported key(s). It could be none if returned in plain text, asymmetric if it is encrypted just with the public key sent in asymmetric_public_key, or kem if it was encrypted using KEM protocol.
	AsymmetricAlgorithm string `json:"asymmetric_algorithm,omitempty"` // The algorithm of the public key used to encrypt exported material
	SymmetricAlgorithm  string `json:"symmetric_algorithm,omitempty"`  // The algorithm of the symmetric key used to encrypt exported material
	KDF                 string `json:"kdf,omitempty"`                  // Key derivation function used to derivate the symmetric key when `encryption_type` is `kem`
	HashAlgorithm       string `json:"hash_algorithm,omitempty"`       // Hash algorithm used to derivate the symmetric key when `encryption_type` is `kem`
	IterationCount      int    `json:"iteration_count,omitempty"`      // Iteration count used to derivate the symmetric key when `encryption_type` is `kem`
	EncryptedSalt       string `json:"encrypted_salt,omitempty"`       // Salt used to derivate the symmetric key when `encryption_type` is `kem`, encrypted with the public key provided in `asymmetric_key`
}

func GetSymmetricKeyLength(algorithm string) (int, error) {
	switch algorithm {
	case string(SYAaes256_cbc), string(SYAaes256_cfb), string(SYAaes256_gcm):
		return 32, nil
	case string(SYAaes128_cfb), string(SYAaes128_cbc):
		return 16, nil
	case string(SYAhs256):
		return 64, nil
	case string(SYAhs384):
		return 96, nil
	case string(SYAhs512):
		return 128, nil
	default:
		return 0, fmt.Errorf("invalid algorithm: '%s'", algorithm)
	}
}
