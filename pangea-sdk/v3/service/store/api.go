package store

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type ItemOrder string

const (
	IOasc ItemOrder = "asc"
	IOdes           = "desc"
)

type ArchiveFormat string

const (
	AFzip ArchiveFormat = "zip"
	AFtar               = "tar"
)

type LinkType string

const (
	LTupload   LinkType = "upload"
	LTdownload          = "download"
	LTall               = "all"
)

type AuthenticatorType string

const (
	ATemailOTP AuthenticatorType = "email_otp"
	ATpassword                   = "password"
	ATsmsOTP                     = "sms_otp"
	ATsocial                     = "social"
)

type ObjectOrderBy string

const (
	OOBid        ObjectOrderBy = "id"
	OOBcreatedAt               = "created_at"
	OOBname                    = "name"
	OOBparendID                = "parent_id"
	OOBtype                    = "type"
	OOBupdatedAt               = "updated_at"
)

type ShareLinkOrderBy string

const (
	SLOBid             ShareLinkOrderBy = "id"
	SLOBstoragePoolID                   = "storage_pool_id"
	SLOBtarget                          = "target"
	SLOBlinkType                        = "link_type"
	SLOBaccessCount                     = "access_count"
	SLOBmaxAccessCount                  = "max_access_count"
	SLOBcreatedAt                       = "created_at"
	SLOBexpiresAt                       = "expires_at"
	SLOBlastAccessedAt                  = "last_accessed_at"
	SLOBlink                            = "link"
)

type Metadata map[string]string
type Tags []string

type DeleteRequest struct {
	pangea.BaseRequest

	ID    string `json:"id,omitempty"`
	Force *bool  `json:"force,omitempty"`
	Path  string `json:"path,omitempty"`
}

type ItemData struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	Size         int      `json:"size"`
	BillableSize int      `json:"billable_size"`
	Location     string   `json:"location"`
	Tags         Tags     `json:"tags,omitempty"`
	Metadata     Metadata `json:"metadata,omitempty"`
	MD5          string   `json:"md5"`
	SHA256       string   `json:"sha256"`
	SHA512       string   `json:"sha512"`
	ParentID     string   `json:"parent_id"`
}

type DeleteResult struct {
	Count int `json:"count"`
}

func (e *store) Delete(ctx context.Context, input *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/delete", input, &DeleteResult{})
}

type FolderCreateRequest struct {
	pangea.BaseRequest

	Name     string   `json:"name,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
	ParentID string   `json:"parent_id,omitempty"`
	Path     string   `json:"path,omitempty"`
	Tags     Tags     `json:"tags,omitempty"`
}

type FolderCreateResult struct {
	Object ItemData `json:"object"`
}

func (e *store) FolderCreate(ctx context.Context, input *FolderCreateRequest) (*pangea.PangeaResponse[FolderCreateResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/folder/create", input, &FolderCreateResult{})
}

type GetRequest struct {
	pangea.BaseRequest

	ID             string                 `json:"id,omitempty"`
	Path           string                 `json:"path,omitempty"`
	TransferMethod *pangea.TransferMethod `json:"transfer_method,omitempty"`
}

type GetResult struct {
	Object  ItemData `json:"object"`
	DestURL *string  `json:"dest_url,omitempty"`
}

func (e *store) Get(ctx context.Context, input *GetRequest) (*pangea.PangeaResponse[GetResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/get", input, &GetResult{})
}

type PutRequest struct {
	pangea.BaseRequest

	Name           string                `json:"name,omitempty"`
	Format         *FileFormat           `json:"format,omitempty"`
	Metadata       Metadata              `json:"metadata,omitempty"`
	MimeType       string                `json:"mimetype,omitempty"`
	ParentID       string                `json:"parent_id,omitempty"`
	Path           string                `json:"path,omitempty"`
	CRC32C         string                `json:"crc32c,omitempty"`
	MD5            string                `json:"md5,omitempty"`
	SHA1           string                `json:"sha1,omitempty"`
	SHA256         string                `json:"sha256,omitempty"`
	SHA512         string                `json:"sha512,omitempty"`
	Size           *int                  `json:"size,omitempty"`
	Tags           Tags                  `json:"tags,omitempty"`
	TransferMethod pangea.TransferMethod `json:"transfer_method,omitempty"`
}

type PutResult struct {
	Object ItemData `json:"object"`
}

func (e *store) Put(ctx context.Context, input *PutRequest, file *os.File) (*pangea.PangeaResponse[PutResult], error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	if input.TransferMethod == pangea.TMpostURL {
		var err error
		params, err := pangea.GetUploadFileParams(file)
		if err != nil {
			return nil, err
		}
		input.CRC32C = params.CRC
		input.SHA256 = params.SHA256
		input.Size = pangea.Int(params.Size)
	}

	name := "file"
	if input.TransferMethod == pangea.TMmultipart {
		name = "upload"
	}

	fd := pangea.FileData{
		File: file,
		Name: name,
	}

	return request.DoPostWithFile(ctx, e.Client, "v1beta/put", input, &PutResult{}, fd)
}

type UpdateRequest struct {
	pangea.BaseRequest

	ID             string   `json:"id"`
	Path           string   `json:"path,omitempty"`
	AddMetadata    Metadata `json:"add_metadata,omitempty"`
	RemoveMetadata Metadata `json:"remove_metadata,omitempty"`
	Metadata       Metadata `json:"metadata,omitempty"`
	AddTags        Tags     `json:"add_tags,omitempty"`
	RemoveTags     Tags     `json:"remove_tags,omitempty"`
	Tags           Tags     `json:"tags,omitempty"`
	ParentID       string   `json:"parent_id,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

type UpdateResult struct {
	Object ItemData `json:"object"`
}

func (e *store) Update(ctx context.Context, input *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/update", input, &UpdateResult{})
}

// Just allowed to filter by folder now
type FilterList struct {
	pangea.FilterBase
	folder *pangea.FilterMatch[string]
}

func NewFilterList() *FilterList {
	filter := make(pangea.Filter)
	return &FilterList{
		FilterBase: *pangea.NewFilterBase(filter),
		folder:     pangea.NewFilterMatch[string]("folder", &filter),
	}
}

func (f *FilterList) Folder() *pangea.FilterMatch[string] {
	return f.folder
}

type ListRequest struct {
	pangea.BaseRequest

	Filter  pangea.Filter `json:"filter,omitempty"`
	Last    string        `json:"last,omitempty"`
	Order   ItemOrder     `json:"order,omitempty"`
	OrderBy ObjectOrderBy `json:"order_by,omitempty"`
	Size    int           `json:"size,omitempty"`
}

type ListResult struct {
	Count   int        `json:"count"`
	Last    string     `json:"last,omitempty"`
	Objects []ItemData `json:"objects"`
}

func (e *store) List(ctx context.Context, input *ListRequest) (*pangea.PangeaResponse[ListResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/list", input, &ListResult{})
}

type GetArchiveRequest struct {
	pangea.BaseRequest

	Ids            []string              `json:"ids"`
	Format         ArchiveFormat         `json:"format,omitempty"`
	TransferMethod pangea.TransferMethod `json:"transfer_method,omitempty"`
}

type GetArchiveResult struct {
	DestURL *string `json:"dest_url,omitempty"`
	Count   int     `json:"count"`
}

func (e *store) GetArchive(ctx context.Context, input *GetArchiveRequest) (*pangea.PangeaResponse[GetArchiveResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/get_archive", input, &GetArchiveResult{})
}

type Authenticator struct {
	AuthType    AuthenticatorType `json:"auth_type"`
	AuthContext string            `json:"auth_context"`
}

type ShareLinkCreateItem struct {
	Targets        []string        `json:"targets"`
	LinkType       LinkType        `json:"link_type,omitempty"`
	ExpiresAt      string          `json:"expires_at,omitempty"`
	MaxAccessCount *int            `json:"max_access_count,omitempty"`
	Authenticators []Authenticator `json:"authenticators,omitempty"`
}

type ShareLinkCreateRequest struct {
	pangea.BaseRequest

	Links []ShareLinkCreateItem `json:"links"`
}

type ShareLinkItem struct {
	ID             string          `json:"id"`
	StoragePoolID  string          `json:"storage_pool_id"`
	Targets        []string        `json:"targets"`
	LinkType       string          `json:"link_type"`
	AccessCount    int             `json:"access_count"`
	MaxAccessCount int             `json:"max_access_count"`
	CreatedAt      string          `json:"created_at"`
	ExpiresAt      string          `json:"expires_at"`
	LastAccessedAt *string         `json:"last_accessed_at,omitempty"`
	Authenticators []Authenticator `json:"authenticators,omitempty"`
	Link           string          `json:"link"`
}

type ShareLinkCreateResult struct {
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

func (e *store) ShareLinkCreate(ctx context.Context, input *ShareLinkCreateRequest) (*pangea.PangeaResponse[ShareLinkCreateResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/share/link/create", input, &ShareLinkCreateResult{})
}

type ShareLinkGetRequest struct {
	pangea.BaseRequest

	ID string `json:"id"`
}

type ShareLinkGetResult struct {
	ShareLinkObject ShareLinkItem `json:"share_link_object"`
}

func (e *store) ShareLinkGet(ctx context.Context, input *ShareLinkGetRequest) (*pangea.PangeaResponse[ShareLinkGetResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/share/link/get", input, &ShareLinkGetResult{})
}

type FilterShareLinkList struct {
	pangea.FilterBase
	id             *pangea.FilterMatch[string]
	storagePoolID  *pangea.FilterMatch[string]
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
		storagePoolID:  pangea.NewFilterMatch[string]("storage_pool_id", &filter),
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

func (f *FilterShareLinkList) StoragePoolID() *pangea.FilterMatch[string] {
	return f.storagePoolID
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

	Filter  pangea.Filter     `json:"filter,omitempty"`
	Last    string            `json:"last,omitempty"`
	Order   *ItemOrder        `json:"order,omitempty"`
	OrderBy *ShareLinkOrderBy `json:"order_by,omitempty"`
	Size    int               `json:"size,omitempty"`
}

type ShareLinkListResult struct {
	Count            int             `json:"count"`
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

func (e *store) ShareLinkList(ctx context.Context, input *ShareLinkListRequest) (*pangea.PangeaResponse[ShareLinkListResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/share/link/list", input, &ShareLinkListResult{})
}

type ShareLinkDeleteRequest struct {
	pangea.BaseRequest

	Ids []string `json:"ids"`
}

type ShareLinkDeleteResult struct {
	ShareLinkObjects []ShareLinkItem `json:"share_link_objects"`
}

func (e *store) ShareLinkDelete(ctx context.Context, input *ShareLinkDeleteRequest) (*pangea.PangeaResponse[ShareLinkDeleteResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/share/link/delete", input, &ShareLinkDeleteResult{})
}

func (e *store) RequestUploadURL(ctx context.Context, input *PutRequest) (*pangea.PangeaResponse[PutResult], error) {
	if input.TransferMethod == pangea.TMpostURL && (input.CRC32C == "" || input.SHA256 == "" || input.Size == nil) {
		return nil, errors.New("Need to set CRC32C, SHA256 and Size in order to use TransferMethod TMpostURL")
	}

	return request.GetUploadURL(ctx, e.Client, "v1beta/put", input, &PutResult{})
}

func (fu *FileUploader) UploadFile(ctx context.Context, url string, tm pangea.TransferMethod, fd pangea.FileData) error {
	if tm == pangea.TMmultipart {
		return errors.New(fmt.Sprintf("%s is not supported in UploadFile. Use Put() instead", tm))
	}

	fds := pangea.FileData{
		File:    fd.File,
		Name:    "file",
		Details: fd.Details,
	}
	return fu.client.UploadFile(ctx, url, tm, fds)
}

type FileFormat string

const (
	FF3g2     FileFormat = "3G2"
	FF3gp                = "3GP"
	FF3mf                = "3MF"
	FF7z                 = "7Z"
	FFa                  = "A"
	FFaac                = "AAC"
	FFaccdb              = "ACCDB"
	FFaiff               = "AIFF"
	FFamf                = "AMF"
	FFamr                = "AMR"
	FFape                = "APE"
	FFasf                = "ASF"
	FFatom               = "ATOM"
	FFau                 = "AU"
	FFavi                = "AVI"
	FFavif               = "AVIF"
	FFbin                = "BIN"
	FFbmp                = "BMP"
	FFbpg                = "BPG"
	FFbz2                = "BZ2"
	FFcab                = "CAB"
	FFclass              = "CLASS"
	FFcpio               = "CPIO"
	FFcrx                = "CRX"
	FFcsv                = "CSV"
	FFdae                = "DAE"
	FFdbf                = "DBF"
	FFdcm                = "DCM"
	FFdeb                = "DEB"
	FFdjvu               = "DJVU"
	FFdll                = "DLL"
	FFdoc                = "DOC"
	FFdocx               = "DOCX"
	FFdwg                = "DWG"
	FFeot                = "EOT"
	FFepub               = "EPUB"
	FFexe                = "EXE"
	FFfdf                = "FDF"
	FFfits               = "FITS"
	FFflac               = "FLAC"
	FFflv                = "FLV"
	FFgbr                = "GBR"
	FFgeojson            = "GEOJSON"
	FFgif                = "GIF"
	FFglb                = "GLB"
	FFgml                = "GML"
	FFgpx                = "GPX"
	FFgz                 = "GZ"
	FFhar                = "HAR"
	FFhdr                = "HDR"
	FFheic               = "HEIC"
	FFheif               = "HEIF"
	FFhtml               = "HTML"
	FFicns               = "ICNS"
	FFico                = "ICO"
	FFics                = "ICS"
	FFiso                = "ISO"
	FFjar                = "JAR"
	FFjp2                = "JP2"
	FFjpf                = "JPF"
	FFjpg                = "JPG"
	FFjpm                = "JPM"
	FFjs                 = "JS"
	FFjson               = "JSON"
	FFjxl                = "JXL"
	FFjxr                = "JXR"
	FFkml                = "KML"
	FFlit                = "LIT"
	FFlnk                = "LNK"
	FFlua                = "LUA"
	FFlz                 = "LZ"
	FFm3u                = "M3U"
	FFm4a                = "M4A"
	FFmacho              = "MACHO"
	FFmdb                = "MDB"
	FFmidi               = "MIDI"
	FFmkv                = "MKV"
	FFmobi               = "MOBI"
	FFmov                = "MOV"
	FFmp3                = "MP3"
	FFmp4                = "MP4"
	FFmpc                = "MPC"
	FFmpeg               = "MPEG"
	FFmqv                = "MQV"
	FFmrc                = "MRC"
	FFmsg                = "MSG"
	FFmsi                = "MSI"
	FFndjson             = "NDJSON"
	FFnes                = "NES"
	FFodc                = "ODC"
	FFodf                = "ODF"
	FFodg                = "ODG"
	FFodp                = "ODP"
	FFods                = "ODS"
	FFodt                = "ODT"
	FFoga                = "OGA"
	FFogv                = "OGV"
	FFotf                = "OTF"
	FFotg                = "OTG"
	FFotp                = "OTP"
	FFots                = "OTS"
	FFott                = "OTT"
	FFowl                = "OWL"
	FFp7s                = "P7S"
	FFpat                = "PAT"
	FFpdf                = "PDF"
	FFphp                = "PHP"
	FFpl                 = "PL"
	FFpng                = "PNG"
	FFppt                = "PPT"
	FFpptx               = "PPTX"
	FFps                 = "PS"
	FFpsd                = "PSD"
	FFpub                = "PUB"
	FFpy                 = "PY"
	FFqcp                = "QCP"
	FFrar                = "RAR"
	FFrmvb               = "RMVB"
	FFrpm                = "RPM"
	FFrss                = "RSS"
	FFrtf                = "RTF"
	FFshp                = "SHP"
	FFshx                = "SHX"
	FFso                 = "SO"
	FFsqlite             = "SQLITE"
	FFsrt                = "SRT"
	FFsvg                = "SVG"
	FFswf                = "SWF"
	FFsxc                = "SXC"
	FFtar                = "TAR"
	FFtcl                = "TCL"
	FFtcx                = "TCX"
	FFtiff               = "TIFF"
	FFtorrent            = "TORRENT"
	FFtsv                = "TSV"
	FFttc                = "TTC"
	FFttf                = "TTF"
	FFtxt                = "TXT"
	FFvcf                = "VCF"
	FFvoc                = "VOC"
	FFvtt                = "VTT"
	FFwarc               = "WARC"
	FFwasm               = "WASM"
	FFwav                = "WAV"
	FFwebm               = "WEBM"
	FFwebp               = "WEBP"
	FFwoff               = "WOFF"
	FFwoff2              = "WOFF2"
	FFx3d                = "X3D"
	FFxar                = "XAR"
	FFxcf                = "XCF"
	FFxfdf               = "XFDF"
	FFxlf                = "XLF"
	FFxls                = "XLS"
	FFxlsx               = "XLSX"
	FFxml                = "XML"
	FFxpm                = "XPM"
	FFxz                 = "XZ"
	FFzip                = "ZIP"
	FFzst                = "ZST"
)
