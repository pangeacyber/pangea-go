package share

import (
	"context"
	"errors"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

type ItemOrder string

const (
	IOasc ItemOrder = "asc"
	IOdes ItemOrder = "desc"
)

type ArchiveFormat string

const (
	AFzip ArchiveFormat = "zip"
	AFtar ArchiveFormat = "tar"
)

type LinkType string

const (
	LTupload   LinkType = "upload"
	LTdownload LinkType = "download"
	LTeditor   LinkType = "editor"
)

type AuthenticatorType string

const (
	ATemailOTP AuthenticatorType = "email_otp"
	ATpassword AuthenticatorType = "password"
	ATsmsOTP   AuthenticatorType = "sms_otp"
	ATsocial   AuthenticatorType = "social"
)

type ObjectOrderBy string

const (
	OOBid        ObjectOrderBy = "id"
	OOBcreatedAt ObjectOrderBy = "created_at"
	OOBname      ObjectOrderBy = "name"
	OOBparendID  ObjectOrderBy = "parent_id"
	OOBtype      ObjectOrderBy = "type"
	OOBupdatedAt ObjectOrderBy = "updated_at"
)

type ShareLinkOrderBy string

const (
	SLOBid             ShareLinkOrderBy = "id"
	SLOBbucketID       ShareLinkOrderBy = "bucket_id"
	SLOBtarget         ShareLinkOrderBy = "target"
	SLOBlinkType       ShareLinkOrderBy = "link_type"
	SLOBaccessCount    ShareLinkOrderBy = "access_count"
	SLOBmaxAccessCount ShareLinkOrderBy = "max_access_count"
	SLOBcreatedAt      ShareLinkOrderBy = "created_at"
	SLOBexpiresAt      ShareLinkOrderBy = "expires_at"
	SLOBlastAccessedAt ShareLinkOrderBy = "last_accessed_at"
	SLOBlink           ShareLinkOrderBy = "link"
)

type Metadata map[string]string
type Tags []string

type DeleteRequest struct {
	pangea.BaseRequest

	ID       string  `json:"id,omitempty"`        // The ID of the object to delete.
	Force    *bool   `json:"force,omitempty"`     // If true, delete a folder even if it's not empty. Deletes the contents of folder as well.
	BucketID *string `json:"bucket_id,omitempty"` // The bucket to use, if not the default.
}

type ItemData struct {
	BillableSize      int      `json:"billable_size"`                // The number of billable bytes (includes Metadata, Tags, etc.) for the object.
	CreatedAt         string   `json:"created_at"`                   // The date and time the object was created.
	ExternalBucketKey string   `json:"external_bucket_key"`          // The key in the external bucket that contains this file.
	Folder            string   `json:"folder"`                       // The full path to the folder the object is stored in.
	ID                string   `json:"id"`                           // The ID of a stored object.
	MD5               string   `json:"md5"`                          // The MD5 hash of the file contents. Cannot be written to.
	Metadata          Metadata `json:"metadata,omitempty"`           // A set of string-based key/value pairs used to provide additional data about an object.
	MetadataProtected Metadata `json:"metadata_protected,omitempty"` // Protected (read-only) metadata.
	Name              string   `json:"name"`                         // The name of the object.
	ParentID          string   `json:"parent_id"`                    // The parent ID (a folder). Blanks means the root folder.
	SHA256            string   `json:"sha256"`                       // The SHA256 hash of the file contents. Cannot be written to.
	SHA512            string   `json:"sha512"`                       // The SHA512 hash of the file contents. Cannot be written to.
	Size              int      `json:"size"`                         // The size of the object in bytes.
	Tags              Tags     `json:"tags,omitempty"`               // A list of user-defined tags
	TagsProtected     Tags     `json:"tags_protected,omitempty"`     // Protected (read-only) flags.
	Type              string   `json:"type"`                         // The type of the item (file or dir). Cannot be written to.
	UpdatedAt         string   `json:"updated_at"`                   // The date and time the object was last updated.
}

type DeleteResult struct {
	Count int `json:"count"` // Number of objects deleted.
}

type Bucket struct {
	Default         bool                    `json:"default"` // If true, is the default bucket.
	ID              string                  `json:"id"`      // The ID of a share bucket resource.
	Name            string                  `json:"name"`    // The bucket's friendly name.
	TransferMethods []pangea.TransferMethod `json:"transfer_methods"`
}

type BucketsResult struct {
	Buckets []Bucket `json:"buckets"` // A list of available buckets.
}

// @summary Buckets
//
// @description Get information on the accessible buckets.
//
// @operationId share_post_v1_buckets
//
// @example
//
//	res, err := shareClient.Buckets(ctx)
func (e *share) Buckets(ctx context.Context) (*pangea.PangeaResponse[BucketsResult], error) {
	return request.DoPost(ctx, e.Client, "v1/buckets", &pangea.BaseRequest{}, &BucketsResult{})
}

// @summary Delete
//
// @description Delete object by ID.
//
// @operationId store_post_v1_delete
//
// @example
//
//	input := &share.DeleteRequest{
//		ID: "pos_3djfmzg2db4c6donarecbyv5begtj2bm"
//	}
//
//	res, err := shareClient.Delete(ctx, input)
func (e *share) Delete(ctx context.Context, input *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error) {
	return request.DoPost(ctx, e.Client, "v1/delete", input, &DeleteResult{})
}

type FolderCreateRequest struct {
	pangea.BaseRequest

	Name     string   `json:"name,omitempty"`      // The name of an object.
	Metadata Metadata `json:"metadata,omitempty"`  // A set of string-based key/value pairs used to provide additional data about an object.
	ParentID string   `json:"parent_id,omitempty"` // The ID of a stored object.
	Folder   string   `json:"folder,omitempty"`    // The folder to place the folder in. Must match `parent_id` if also set.
	Tags     Tags     `json:"tags,omitempty"`      // A list of user-defined tags
	BucketID *string  `json:"bucket_id,omitempty"` // The bucket to use, if not the default.
}

type FolderCreateResult struct {
	Object ItemData `json:"object"` // Information on the created folder.
}

// @summary Create a folder
//
// @description Create a folder, either by name or path and parent_id.
//
// @operationId store_post_v1_folder_create
//
// @example
//
//	input := &share.FolderCreateRequest{
//		Metadata: share.Metadata{
//			"created_by": "jim",
//			"priority": "medium",
//		},
//		ParentID: "pos_3djfmzg2db4c6donarecbyv5begtj2bm",
//		Folder: "/",
//		Tags: share.Tags{"irs_2023", "personal"},
//	}
//
//	res, err := shareClient.FolderCreate(ctx, input)
func (e *share) FolderCreate(ctx context.Context, input *FolderCreateRequest) (*pangea.PangeaResponse[FolderCreateResult], error) {
	return request.DoPost(ctx, e.Client, "v1/folder/create", input, &FolderCreateResult{})
}

type GetRequest struct {
	pangea.BaseRequest

	ID             string                `json:"id,omitempty"`              // The ID of the object to retrieve.
	Password       *string               `json:"password,omitempty"`        // If the file was protected with a password, the password to decrypt with.
	TransferMethod pangea.TransferMethod `json:"transfer_method,omitempty"` // The requested transfer method for the file data.
	BucketID       *string               `json:"bucket_id,omitempty"`       // The bucket to use, if not the default.
}

type GetResult struct {
	Object  ItemData `json:"object"`             // File information.
	DestURL *string  `json:"dest_url,omitempty"` // A URL where the file can be downloaded from. (transfer_method: dest-url)
}

// @summary Get an object
//
// @description Get object.
//
// @operationId store_post_v1_get
//
// @example
//
//	input := &share.GetRequest{
//		ID: "pos_3djfmzg2db4c6donarecbyv5begtj2bm",
//	}
//
//	res, err := shareClient.Get(ctx, input)
func (e *share) Get(ctx context.Context, input *GetRequest) (*pangea.PangeaResponse[GetResult], error) {
	return request.DoPost(ctx, e.Client, "v1/get", input, &GetResult{})
}

type PutRequest struct {
	pangea.BaseRequest
	pangea.TransferRequest

	Size              *int        `json:"size,omitempty"`               // The size (in bytes) of the file. If the upload doesn't match, the call will fail.
	BucketID          *string     `json:"bucket_id,omitempty"`          // The bucket to use, if not the default.
	CRC32C            string      `json:"crc32c,omitempty"`             // The hexadecimal-encoded CRC32C hash of the file data, which will be verified by the server if provided.
	SHA256            string      `json:"sha256,omitempty"`             // The SHA256 hash of the file data, which will be verified by the server if provided.
	MD5               string      `json:"md5,omitempty"`                // The hexadecimal-encoded MD5 hash of the file data, which will be verified by the server if provided.
	Name              string      `json:"name,omitempty"`               // The name of the object to store.
	Format            *FileFormat `json:"format,omitempty"`             // The format of the file, which will be verified by the server if provided. Uploads not matching the supplied format will be rejected.
	Metadata          Metadata    `json:"metadata,omitempty"`           // A set of string-based key/value pairs used to provide additional data about an object.
	MimeType          string      `json:"mimetype,omitempty"`           // The MIME type of the file, which will be verified by the server if provided. Uploads not matching the supplied MIME type will be rejected.
	ParentID          string      `json:"parent_id,omitempty"`          // The parent ID of the object (a folder). Leave blank to keep in the root folder.
	Folder            string      `json:"folder,omitempty"`             // The path to the parent folder. Leave blank for the root folder. Path must resolve to `parent_id` if also set.
	Password          string      `json:"password,omitempty"`           // An optional password to protect the file with. Downloading the file will require this password.
	PasswordAlgorithm string      `json:"password_algorithm,omitempty"` // An optional password algorithm to protect the file with. See symmetric vault password_algorithm.
	SHA1              string      `json:"sha1,omitempty"`               // The hexadecimal-encoded SHA1 hash of the file data, which will be verified by the server if provided.
	SHA512            string      `json:"sha512,omitempty"`             // The hexadecimal-encoded SHA512 hash of the file data, which will be verified by the server if provided.
	SourceURL         string      `json:"source_url,omitempty"`         // The URL to fetch the file payload from (for transfer_method source-url).
	Tags              Tags        `json:"tags,omitempty"`               // A list of user-defined tags
}

type PutResult struct {
	Object ItemData `json:"object"`
}

// @summary Upload a file
//
// @description Upload a file.
//
// @operationId store_post_v1_put
//
// @example
//
//	input := &share.PutRequest{
//		TransferMethod: pangea.TMmultipart,
//		Metadata: share.Metadata{
//			"created_by": "jim",
//			"priority": "medium",
//		},
//		ParentID: "pos_3djfmzg2db4c6donarecbyv5begtj2bm",
//		Folder: "/",
//		Tags: share.Tags{"irs_2023", "personal"},
//	}
//
//	file, err := os.Open("./path/to/file.pdf")
//	if err != nil {
//		log.Fatal("Error opening file: %v", err)
//	}
//
//	res, err := shareClient.Put(ctx, input, file)
func (e *share) Put(ctx context.Context, input *PutRequest, file *os.File) (*pangea.PangeaResponse[PutResult], error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	if input.TransferMethod == pangea.TMpostURL {
		var err error
		params, err := pangea.GetUploadFileParams(file)
		if err != nil {
			return nil, err
		}
		input.CRC32C = params.CRC32C
		input.SHA256 = params.SHA256
		input.Size = pangea.Int(params.Size)
	} else if size, err := pangea.GetFileSize(file); err == nil && size == 0 {
		input.Size = pangea.Int(0)
	}

	name := "file"
	if input.TransferMethod == pangea.TMmultipart {
		name = "upload"
	}

	fd := pangea.FileData{
		File: file,
		Name: name,
	}

	return request.DoPostWithFile(ctx, e.Client, "v1/put", input, &PutResult{}, fd)
}

type UpdateRequest struct {
	pangea.BaseRequest

	ID                   string   `json:"id"`                               // An identifier for the file to update.
	Folder               string   `json:"folder,omitempty"`                 // Set the parent (folder). Leave blank for the root folder. Path must resolve to `parent_id` if also set.
	BucketID             *string  `json:"bucket_id,omitempty"`              // The bucket to use, if not the default.
	AddMetadata          Metadata `json:"add_metadata,omitempty"`           // A list of Metadata key/values to set in the object. If a provided key exists, the value will be replaced.
	AddPassword          string   `json:"add_password,omitempty"`           // Protect the file with the supplied password.
	AddPasswordAlgorithm string   `json:"add_password_algorithm,omitempty"` // The algorithm to use to password protect the file.
	AddTags              Tags     `json:"add_tags,omitempty"`               // A list of Tags to add. It is not an error to provide a tag which already exists.
	Name                 string   `json:"name,omitempty"`                   // Sets the object's Name.
	Metadata             Metadata `json:"metadata,omitempty"`               // Set the object's metadata.
	RemoveMetadata       Metadata `json:"remove_metadata,omitempty"`        // A list of metadata key/values to remove in the object. It is not an error for a provided key to not exist. If a provided key exists but doesn't match the provided value, it will not be removed.
	RemovePassword       string   `json:"remove_password,omitempty"`        // Remove the supplied password from the file.
	RemoveTags           Tags     `json:"remove_tags,omitempty"`            // A list of tags to remove. It is not an error to provide a tag which is not present.
	ParentID             string   `json:"parent_id,omitempty"`              // Set the parent (folder) of the object. Can be an empty string for the root folder.
	Tags                 Tags     `json:"tags,omitempty"`                   // Set the object's tags.
	UpdatedAt            string   `json:"updated_at,omitempty"`             // The date and time the object was last updated. If included, the update will fail if this doesn't match the date and time of the last update for the object.
}

type UpdateResult struct {
	Object ItemData `json:"object"`
}

// @summary Update a file
//
// @description Update a file.
//
// @operationId share_post_v1_update
//
// @example
//
//	input := &share.UpdateRequest{
//		ID: "pos_3djfmzg2db4c6donarecbyv5begtj2bm",
//		Folder: "/",
//		RemoveMetadata: share.Metadata{
//			"created_by": "jim",
//			"priority": "medium",
//		},
//		RemoveTags: share.Tags{"irs_2023", "personal"},
//	}
//
//	res, err := shareClient.Update(ctx, input)
func (e *share) Update(ctx context.Context, input *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error) {
	return request.DoPost(ctx, e.Client, "v1/update", input, &UpdateResult{})
}

// Just allowed to filter by folder now
type FilterList struct {
	pangea.FilterBase
	folder    *pangea.FilterEqual[string]
	createdAt *pangea.FilterRange[string]
	id        *pangea.FilterMatch[string]
	name      *pangea.FilterMatch[string]
	parentId  *pangea.FilterMatch[string]
	size      *pangea.FilterRange[string]
	tags      *pangea.FilterEqual[[]string]
	type_     *pangea.FilterMatch[string]
	updatedAt *pangea.FilterRange[string]
}

func NewFilterList() *FilterList {
	filter := make(pangea.Filter)
	return &FilterList{
		FilterBase: *pangea.NewFilterBase(filter),
		folder:     pangea.NewFilterEqual[string]("folder", &filter),
		createdAt:  pangea.NewFilterRange[string]("created_at", &filter),
		id:         pangea.NewFilterMatch[string]("id", &filter),
		name:       pangea.NewFilterMatch[string]("name", &filter),
		parentId:   pangea.NewFilterMatch[string]("parent_id", &filter),
		size:       pangea.NewFilterRange[string]("size", &filter),
		tags:       pangea.NewFilterEqual[[]string]("tags", &filter),
		type_:      pangea.NewFilterMatch[string]("type", &filter),
		updatedAt:  pangea.NewFilterRange[string]("updated_at", &filter),
	}
}

func (f *FilterList) Folder() *pangea.FilterEqual[string] {
	return f.folder
}

func (f *FilterList) Tags() *pangea.FilterEqual[[]string] {
	return f.tags
}

func (f *FilterList) CreatedAt() *pangea.FilterRange[string] {
	return f.createdAt
}

func (f *FilterList) ID() *pangea.FilterMatch[string] {
	return f.id
}

func (f *FilterList) Name() *pangea.FilterMatch[string] {
	return f.name
}

func (f *FilterList) ParentID() *pangea.FilterMatch[string] {
	return f.parentId
}

func (f *FilterList) Size() *pangea.FilterRange[string] {
	return f.size
}

func (f *FilterList) Type() *pangea.FilterMatch[string] {
	return f.type_
}

func (f *FilterList) UpdatedAt() *pangea.FilterRange[string] {
	return f.updatedAt
}

type ListRequest struct {
	pangea.BaseRequest

	BucketID                 *string       `json:"bucket_id,omitempty"`                   // The bucket to use, if not the default.
	IncludeExternalBucketKey *bool         `json:"include_external_bucket_key,omitempty"` // If true, include the `external_bucket_key` in results.
	Filter                   pangea.Filter `json:"filter,omitempty"`
	Last                     string        `json:"last,omitempty"`     // Reflected value from a previous response to obtain the next page of results.
	Order                    ItemOrder     `json:"order,omitempty"`    // Order results asc(ending) or desc(ending).
	OrderBy                  ObjectOrderBy `json:"order_by,omitempty"` // Which field to order results by.
	Size                     int           `json:"size,omitempty"`     // Maximum results to include in the response.
}

type ListResult struct {
	Count   int        `json:"count"`          // The total number of objects matched by the list request.
	Last    string     `json:"last,omitempty"` // Used to fetch the next page of the current listing when provided in a repeated request's last parameter.
	Objects []ItemData `json:"objects"`
}

// @summary List
//
// @description List or filter/search records.
//
// @operationId share_post_v1_list
//
// @example
//
//	input := &share.ListRequest{}
//
//	res, err := shareClient.List(ctx, input)
func (e *share) List(ctx context.Context, input *ListRequest) (*pangea.PangeaResponse[ListResult], error) {
	return request.DoPost(ctx, e.Client, "v1/list", input, &ListResult{})
}

type GetArchiveRequest struct {
	pangea.BaseRequest

	Ids            []string              `json:"ids"`                       // The IDs of the objects to include in the archive. Folders include all children.
	Format         ArchiveFormat         `json:"format,omitempty"`          // The format to use to build the archive.
	TransferMethod pangea.TransferMethod `json:"transfer_method,omitempty"` // The requested transfer method for the file data.
	BucketID       *string               `json:"bucket_id,omitempty"`       // The bucket to use, if not the default.
}

type GetArchiveResult struct {
	DestURL *string    `json:"dest_url,omitempty"` // A location where the archive can be downloaded from. (transfer_method: dest-url)
	Count   int        `json:"count"`              // Number of objects included in the archive.
	Objects []ItemData `json:"objects"`            // A list of all objects included in the archive.
}

// @summary Get archive
//
// @description Get an archive file of multiple objects.
//
// @operationId share_post_v1_get_archive
//
// @example
//
//	input := &share.GetArchiveRequest{
//		Ids: []string{"pos_3djfmzg2db4c6donarecbyv5begtj2bm"},
//	}
//
//	res, err := shareClient.GetArchive(ctx, input)
func (e *share) GetArchive(ctx context.Context, input *GetArchiveRequest) (*pangea.PangeaResponse[GetArchiveResult], error) {
	return request.DoPost(ctx, e.Client, "v1/get_archive", input, &GetArchiveResult{})
}

type Authenticator struct {
	AuthType    AuthenticatorType `json:"auth_type"`    // An authentication mechanism
	AuthContext string            `json:"auth_context"` // An email address, a phone number or a password to access share link
}

type ShareLinkCreateItem struct {
	Targets        []string        `json:"targets"`                    // List of storage IDs
	LinkType       LinkType        `json:"link_type,omitempty"`        // Type of link
	ExpiresAt      string          `json:"expires_at,omitempty"`       // The date and time the share link expires.
	MaxAccessCount *int            `json:"max_access_count,omitempty"` // The maximum number of times a user can be authenticated to access the share link.
	Authenticators []Authenticator `json:"authenticators,omitempty"`   // A list of authenticators
	Title          string          `json:"title,omitempty"`            // An optional title to use in accessing shares.
	Message        string          `json:"message,omitempty"`          // An optional message to use in accessing shares.
	NotifyEmail    string          `json:"notify_email,omitempty"`     // An email address
	Tags           Tags            `json:"tags,omitempty"`             // A list of user-defined tags
}

type ShareLinkCreateRequest struct {
	pangea.BaseRequest

	Links []ShareLinkCreateItem `json:"links"`
}

type ShareLinkItem struct {
	ID             string          `json:"id"`                         // The ID of a share link.
	BucketID       string          `json:"bucket_id"`                  // The ID of a share bucket resource.
	Targets        []string        `json:"targets"`                    // List of storage IDs
	LinkType       string          `json:"link_type"`                  // Type of link
	AccessCount    int             `json:"access_count"`               // The number of times a user has authenticated to access the share link.
	MaxAccessCount int             `json:"max_access_count"`           // The maximum number of times a user can be authenticated to access the share link.
	CreatedAt      string          `json:"created_at"`                 // The date and time the share link was created.
	ExpiresAt      string          `json:"expires_at"`                 // The date and time the share link expires.
	LastAccessedAt *string         `json:"last_accessed_at,omitempty"` // The date and time the share link was last accessed.
	Authenticators []Authenticator `json:"authenticators,omitempty"`   // A list of authenticators
	Title          string          `json:"title,omitempty"`            // An optional title to use in accessing shares.
	Message        string          `json:"message,omitempty"`          // An optional message to use in accessing shares.
	Link           string          `json:"link"`                       // A URL to access the file/folders shared with a link.
	NotifyEmail    string          `json:"notify_email,omitempty"`     // An email address
	Tags           Tags            `json:"tags,omitempty"`             // A list of user-defined tags
}

type ShareLinkCreateResult struct {
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

// @summary Create share links
//
// @description Create a share link.
//
// @operationId share_post_v1_share_link_create
//
// @example
//
//	authenticator := share.Authenticator{
//		AuthType: share.ATpassword,
//		AuthContext: "my_fav_Pa55word",
//	}
//
//	link := share.ShareLinkCreateItem{
//		Targets: []string{"pos_3djfmzg2db4c6donarecbyv5begtj2bm"},
//		LinkType: "download",
//		Authenticators: []Authenticator{authenticator},
//	}
//
//	input := &share.ShareLinkCreateRequest{
//		Links: []share.ShareLinkCreateItem{link},
//	}
//
//	res, err := shareClient.ShareLinkCreate(ctx, input)
func (e *share) ShareLinkCreate(ctx context.Context, input *ShareLinkCreateRequest) (*pangea.PangeaResponse[ShareLinkCreateResult], error) {
	return request.DoPost(ctx, e.Client, "v1/share/link/create", input, &ShareLinkCreateResult{})
}

type ShareLinkGetRequest struct {
	pangea.BaseRequest

	ID string `json:"id"` // The ID of a share link.
}

type ShareLinkGetResult struct {
	ShareLinkObject ShareLinkItem `json:"share_link_object"`
}

// @summary Get share link
//
// @description Get a share link.
//
// @operationId share_post_v1_share_link_get
//
// @example
//
//	input := &share.ShareLinkGetRequest{
//		ID: "psl_3djfmzg2db4c6donarecbyv5begtj2bm",
//	}
//
//	res, err := shareClient.ShareLinkGet(ctx, input)
func (e *share) ShareLinkGet(ctx context.Context, input *ShareLinkGetRequest) (*pangea.PangeaResponse[ShareLinkGetResult], error) {
	return request.DoPost(ctx, e.Client, "v1/share/link/get", input, &ShareLinkGetResult{})
}

type FilterShareLinkList struct {
	pangea.FilterBase
	id             *pangea.FilterMatch[string]
	target         *pangea.FilterMatch[string]
	linkType       *pangea.FilterMatch[string]
	accessCount    *pangea.FilterRange[string]
	maxAccessCount *pangea.FilterRange[string]
	createdAt      *pangea.FilterRange[string]
	expiresAt      *pangea.FilterRange[string]
	lastAccessedAt *pangea.FilterRange[string]
	link           *pangea.FilterMatch[string]
}

func NewFilterShareLinkList() *FilterShareLinkList {
	filter := make(pangea.Filter)
	return &FilterShareLinkList{
		FilterBase:     *pangea.NewFilterBase(filter),
		id:             pangea.NewFilterMatch[string]("id", &filter),
		target:         pangea.NewFilterMatch[string]("target", &filter),
		linkType:       pangea.NewFilterMatch[string]("link_type", &filter),
		accessCount:    pangea.NewFilterRange[string]("access_count", &filter),
		maxAccessCount: pangea.NewFilterRange[string]("max_access_count", &filter),
		createdAt:      pangea.NewFilterRange[string]("created_at", &filter),
		expiresAt:      pangea.NewFilterRange[string]("expires_at", &filter),
		lastAccessedAt: pangea.NewFilterRange[string]("last_accessed_at", &filter),
		link:           pangea.NewFilterMatch[string]("link", &filter),
	}
}

func (f *FilterShareLinkList) ID() *pangea.FilterMatch[string] {
	return f.id
}

func (f *FilterShareLinkList) Target() *pangea.FilterMatch[string] {
	return f.target
}

func (f *FilterShareLinkList) LinkType() *pangea.FilterMatch[string] {
	return f.linkType
}

func (f *FilterShareLinkList) Link() *pangea.FilterMatch[string] {
	return f.link
}

func (f *FilterShareLinkList) AccessCount() *pangea.FilterRange[string] {
	return f.accessCount
}

func (f *FilterShareLinkList) MaxAccessCount() *pangea.FilterRange[string] {
	return f.maxAccessCount
}

func (f *FilterShareLinkList) CreatedAt() *pangea.FilterRange[string] {
	return f.createdAt
}

func (f *FilterShareLinkList) ExpiresAt() *pangea.FilterRange[string] {
	return f.expiresAt
}

func (f *FilterShareLinkList) LastAccessedAt() *pangea.FilterRange[string] {
	return f.lastAccessedAt
}

type ShareLinkListRequest struct {
	pangea.BaseRequest

	BucketID *string           `json:"bucket_id,omitempty"` // The bucket to use, if not the default.
	Filter   pangea.Filter     `json:"filter,omitempty"`
	Last     string            `json:"last,omitempty"`     // Reflected value from a previous response to obtain the next page of results.
	Order    *ItemOrder        `json:"order,omitempty"`    // Order results asc(ending) or desc(ending).
	OrderBy  *ShareLinkOrderBy `json:"order_by,omitempty"` // Which field to order results by.
	Size     int               `json:"size,omitempty"`     // Maximum results to include in the response.
}

type ShareLinkListResult struct {
	Count            int             `json:"count"`  // The total number of share links matched by the list request.
	Last             *string         `json:"string"` // Used to fetch the next page of the current listing when provided in a repeated request's last parameter.
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

// @summary List share links
//
// @description Look up share links by filter options.
//
// @operationId share_post_v1_share_link_list
//
// @example
//
//	input := &share.ShareLinkListRequest{}
//
//	res, err := shareClient.ShareLinkList(ctx, input)
func (e *share) ShareLinkList(ctx context.Context, input *ShareLinkListRequest) (*pangea.PangeaResponse[ShareLinkListResult], error) {
	return request.DoPost(ctx, e.Client, "v1/share/link/list", input, &ShareLinkListResult{})
}

type ShareLinkDeleteRequest struct {
	pangea.BaseRequest

	BucketID *string  `json:"bucket_id,omitempty"` // The bucket to use, if not the default.
	Ids      []string `json:"ids"`
}

type ShareLinkDeleteResult struct {
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

// @summary Delete share links
//
// @description Delete share links.
//
// @operationId share_post_v1_share_link_delete
//
// @example
//
//	input := &share.ShareLinkDeleteRequest{
//		Ids: []string{"psl_3djfmzg2db4c6donarecbyv5begtj2bm"},
//	}
//
//	res, err := shareClient.ShareLinkDelete(ctx, input)
func (e *share) ShareLinkDelete(ctx context.Context, input *ShareLinkDeleteRequest) (*pangea.PangeaResponse[ShareLinkDeleteResult], error) {
	return request.DoPost(ctx, e.Client, "v1/share/link/delete", input, &ShareLinkDeleteResult{})
}

// @summary Request upload URL
//
// @description Request an upload URL.
//
// @operationId share_post_v1_put 2
//
// @example
//
//	input := &share.PutRequest{
//		TransferMethod: pangea.TMpostURL,
//		CRC32C: "515f7c32",
//		SHA256: "c0b56b1a154697f79d27d57a3a2aad4c93849aa2239cd23048fc6f45726271cc",
//		Size: 222089,
//		Metadata: share.Metadata{
//			"created_by": "jim",
//			"priority": "medium",
//		},
//		ParentID: "pos_3djfmzg2db4c6donarecbyv5begtj2bm",
//		Folder: "/",
//		Tags: share.Tags{"irs_2023", "personal"},
//	}
//
//	res, err := shareClient.RequestUploadURL(ctx, input)
func (e *share) RequestUploadURL(ctx context.Context, input *PutRequest) (*pangea.PangeaResponse[PutResult], error) {
	if input.TransferMethod == pangea.TMpostURL && (input.CRC32C == "" || input.SHA256 == "" || input.Size == nil) {
		return nil, errors.New("need to set CRC32C, SHA256 and Size in order to use TransferMethod TMpostURL")
	}

	return request.GetUploadURL(ctx, e.Client, "v1/put", input, &PutResult{})
}

type ShareLinkSendItem struct {
	Id    string `json:"id"`    // The ID of a share link.
	Email string `json:"email"` // An email address
}

type ShareLinkSendRequest struct {
	pangea.BaseRequest

	Links       []ShareLinkSendItem `json:"links"`
	SenderEmail string              `json:"sender_email"`          // An email address
	SenderName  string              `json:"sender_name,omitempty"` // The sender name information. Can be sender's full name for example.
}

type ShareLinkSendResult struct {
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

// @summary Send share link(s)
//
// @description Send share link(s).
//
// @operationId share_post_v1_share_link_send
//
// @example
//
//	res, err := client.ShareLinkSend(ctx, &share.ShareLinkSendRequest{
//		Links: []share.ShareLinkSendItem{
//			share.ShareLinkSendItem{
//				Id:    link.ID,
//				Email: "user@email.com",
//			},
//		},
//		SenderEmail: "sender@email.com",
//		SenderName:  "Sender Name",
//	})
func (e *share) ShareLinkSend(ctx context.Context, input *ShareLinkSendRequest) (*pangea.PangeaResponse[ShareLinkSendResult], error) {
	return request.DoPost(ctx, e.Client, "v1/share/link/send", input, &ShareLinkSendResult{})
}

type FileFormat string

const (
	FF3g2     FileFormat = "3G2"
	FF3gp     FileFormat = "3GP"
	FF3mf     FileFormat = "3MF"
	FF7z      FileFormat = "7Z"
	FFa       FileFormat = "A"
	FFaac     FileFormat = "AAC"
	FFaccdb   FileFormat = "ACCDB"
	FFaiff    FileFormat = "AIFF"
	FFamf     FileFormat = "AMF"
	FFamr     FileFormat = "AMR"
	FFape     FileFormat = "APE"
	FFasf     FileFormat = "ASF"
	FFatom    FileFormat = "ATOM"
	FFau      FileFormat = "AU"
	FFavi     FileFormat = "AVI"
	FFavif    FileFormat = "AVIF"
	FFbin     FileFormat = "BIN"
	FFbmp     FileFormat = "BMP"
	FFbpg     FileFormat = "BPG"
	FFbz2     FileFormat = "BZ2"
	FFcab     FileFormat = "CAB"
	FFclass   FileFormat = "CLASS"
	FFcpio    FileFormat = "CPIO"
	FFcrx     FileFormat = "CRX"
	FFcsv     FileFormat = "CSV"
	FFdae     FileFormat = "DAE"
	FFdbf     FileFormat = "DBF"
	FFdcm     FileFormat = "DCM"
	FFdeb     FileFormat = "DEB"
	FFdjvu    FileFormat = "DJVU"
	FFdll     FileFormat = "DLL"
	FFdoc     FileFormat = "DOC"
	FFdocx    FileFormat = "DOCX"
	FFdwg     FileFormat = "DWG"
	FFeot     FileFormat = "EOT"
	FFepub    FileFormat = "EPUB"
	FFexe     FileFormat = "EXE"
	FFfdf     FileFormat = "FDF"
	FFfits    FileFormat = "FITS"
	FFflac    FileFormat = "FLAC"
	FFflv     FileFormat = "FLV"
	FFgbr     FileFormat = "GBR"
	FFgeojson FileFormat = "GEOJSON"
	FFgif     FileFormat = "GIF"
	FFglb     FileFormat = "GLB"
	FFgml     FileFormat = "GML"
	FFgpx     FileFormat = "GPX"
	FFgz      FileFormat = "GZ"
	FFhar     FileFormat = "HAR"
	FFhdr     FileFormat = "HDR"
	FFheic    FileFormat = "HEIC"
	FFheif    FileFormat = "HEIF"
	FFhtml    FileFormat = "HTML"
	FFicns    FileFormat = "ICNS"
	FFico     FileFormat = "ICO"
	FFics     FileFormat = "ICS"
	FFiso     FileFormat = "ISO"
	FFjar     FileFormat = "JAR"
	FFjp2     FileFormat = "JP2"
	FFjpf     FileFormat = "JPF"
	FFjpg     FileFormat = "JPG"
	FFjpm     FileFormat = "JPM"
	FFjs      FileFormat = "JS"
	FFjson    FileFormat = "JSON"
	FFjxl     FileFormat = "JXL"
	FFjxr     FileFormat = "JXR"
	FFkml     FileFormat = "KML"
	FFlit     FileFormat = "LIT"
	FFlnk     FileFormat = "LNK"
	FFlua     FileFormat = "LUA"
	FFlz      FileFormat = "LZ"
	FFm3u     FileFormat = "M3U"
	FFm4a     FileFormat = "M4A"
	FFmacho   FileFormat = "MACHO"
	FFmdb     FileFormat = "MDB"
	FFmidi    FileFormat = "MIDI"
	FFmkv     FileFormat = "MKV"
	FFmobi    FileFormat = "MOBI"
	FFmov     FileFormat = "MOV"
	FFmp3     FileFormat = "MP3"
	FFmp4     FileFormat = "MP4"
	FFmpc     FileFormat = "MPC"
	FFmpeg    FileFormat = "MPEG"
	FFmqv     FileFormat = "MQV"
	FFmrc     FileFormat = "MRC"
	FFmsg     FileFormat = "MSG"
	FFmsi     FileFormat = "MSI"
	FFndjson  FileFormat = "NDJSON"
	FFnes     FileFormat = "NES"
	FFodc     FileFormat = "ODC"
	FFodf     FileFormat = "ODF"
	FFodg     FileFormat = "ODG"
	FFodp     FileFormat = "ODP"
	FFods     FileFormat = "ODS"
	FFodt     FileFormat = "ODT"
	FFoga     FileFormat = "OGA"
	FFogv     FileFormat = "OGV"
	FFotf     FileFormat = "OTF"
	FFotg     FileFormat = "OTG"
	FFotp     FileFormat = "OTP"
	FFots     FileFormat = "OTS"
	FFott     FileFormat = "OTT"
	FFowl     FileFormat = "OWL"
	FFp7s     FileFormat = "P7S"
	FFpat     FileFormat = "PAT"
	FFpdf     FileFormat = "PDF"
	FFphp     FileFormat = "PHP"
	FFpl      FileFormat = "PL"
	FFpng     FileFormat = "PNG"
	FFppt     FileFormat = "PPT"
	FFpptx    FileFormat = "PPTX"
	FFps      FileFormat = "PS"
	FFpsd     FileFormat = "PSD"
	FFpub     FileFormat = "PUB"
	FFpy      FileFormat = "PY"
	FFqcp     FileFormat = "QCP"
	FFrar     FileFormat = "RAR"
	FFrmvb    FileFormat = "RMVB"
	FFrpm     FileFormat = "RPM"
	FFrss     FileFormat = "RSS"
	FFrtf     FileFormat = "RTF"
	FFshp     FileFormat = "SHP"
	FFshx     FileFormat = "SHX"
	FFso      FileFormat = "SO"
	FFsqlite  FileFormat = "SQLITE"
	FFsrt     FileFormat = "SRT"
	FFsvg     FileFormat = "SVG"
	FFswf     FileFormat = "SWF"
	FFsxc     FileFormat = "SXC"
	FFtar     FileFormat = "TAR"
	FFtcl     FileFormat = "TCL"
	FFtcx     FileFormat = "TCX"
	FFtiff    FileFormat = "TIFF"
	FFtorrent FileFormat = "TORRENT"
	FFtsv     FileFormat = "TSV"
	FFttc     FileFormat = "TTC"
	FFttf     FileFormat = "TTF"
	FFtxt     FileFormat = "TXT"
	FFvcf     FileFormat = "VCF"
	FFvoc     FileFormat = "VOC"
	FFvtt     FileFormat = "VTT"
	FFwarc    FileFormat = "WARC"
	FFwasm    FileFormat = "WASM"
	FFwav     FileFormat = "WAV"
	FFwebm    FileFormat = "WEBM"
	FFwebp    FileFormat = "WEBP"
	FFwoff    FileFormat = "WOFF"
	FFwoff2   FileFormat = "WOFF2"
	FFx3d     FileFormat = "X3D"
	FFxar     FileFormat = "XAR"
	FFxcf     FileFormat = "XCF"
	FFxfdf    FileFormat = "XFDF"
	FFxlf     FileFormat = "XLF"
	FFxls     FileFormat = "XLS"
	FFxlsx    FileFormat = "XLSX"
	FFxml     FileFormat = "XML"
	FFxpm     FileFormat = "XPM"
	FFxz      FileFormat = "XZ"
	FFzip     FileFormat = "ZIP"
	FFzst     FileFormat = "ZST"
)
