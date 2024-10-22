package vault

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

type Client interface {
	StateChange(ctx context.Context, req *StateChangeRequest) (*pangea.PangeaResponse[StateChangeResult], error)
	Delete(ctx context.Context, req *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error)
	Get(ctx context.Context, req *GetRequest) (*pangea.PangeaResponse[GetResult], error)
	JWKGet(ctx context.Context, req *JWKGetRequest) (*pangea.PangeaResponse[JWKGetResult], error)
	List(ctx context.Context, req *ListRequest) (*pangea.PangeaResponse[ListResult], error)
	Update(ctx context.Context, req *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error)
	SecretStore(ctx context.Context, req *SecretStoreRequest) (*pangea.PangeaResponse[SecretStoreResult], error)
	SecretRotate(ctx context.Context, req *SecretRotateRequest) (*pangea.PangeaResponse[SecretRotateResult], error)
	SymmetricGenerate(ctx context.Context, req *SymmetricGenerateRequest) (*pangea.PangeaResponse[SymmetricGenerateResult], error)
	AsymmetricGenerate(ctx context.Context, req *AsymmetricGenerateRequest) (*pangea.PangeaResponse[AsymmetricGenerateResult], error)
	SymmetricStore(ctx context.Context, req *SymmetricStoreRequest) (*pangea.PangeaResponse[SymmetricStoreResult], error)
	AsymmetricStore(ctx context.Context, req *AsymmetricStoreRequest) (*pangea.PangeaResponse[AsymmetricStoreResult], error)
	KeyRotate(ctx context.Context, req *KeyRotateRequest) (*pangea.PangeaResponse[KeyRotateResult], error)
	Encrypt(ctx context.Context, req *EncryptRequest) (*pangea.PangeaResponse[EncryptResult], error)
	Decrypt(ctx context.Context, req *DecryptRequest) (*pangea.PangeaResponse[DecryptResult], error)
	Sign(ctx context.Context, req *SignRequest) (*pangea.PangeaResponse[SignResult], error)
	Verify(ctx context.Context, req *VerifyRequest) (*pangea.PangeaResponse[VerifyResult], error)
	JWTSign(ctx context.Context, req *JWTSignRequest) (*pangea.PangeaResponse[JWTSignResult], error)
	JWTVerify(ctx context.Context, req *JWTVerifyRequest) (*pangea.PangeaResponse[JWTVerifyResult], error)
	FolderCreate(ctx context.Context, req *FolderCreateRequest) (*pangea.PangeaResponse[FolderCreateResult], error)

	// Encrypt parts of a JSON object.
	EncryptStructured(ctx context.Context, input *EncryptStructuredRequest) (*pangea.PangeaResponse[EncryptStructuredResult], error)

	// Decrypt parts of a JSON object.
	DecryptStructured(ctx context.Context, input *EncryptStructuredRequest) (*pangea.PangeaResponse[EncryptStructuredResult], error)

	// Encrypt using a format-preserving algorithm (FPE).
	EncryptTransform(ctx context.Context, input *EncryptTransformRequest) (*pangea.PangeaResponse[EncryptTransformResult], error)

	// Decrypt using a format-preserving algorithm (FPE).
	DecryptTransform(ctx context.Context, input *DecryptTransformRequest) (*pangea.PangeaResponse[DecryptTransformResult], error)

	EncryptTransformStructured(ctx context.Context, input *EncryptTransformStructuredRequest) (*pangea.PangeaResponse[EncryptTransformStructuredResult], error)

	DecryptTransformStructured(ctx context.Context, input *EncryptTransformStructuredRequest) (*pangea.PangeaResponse[EncryptTransformStructuredResult], error)

	// Export a symmetric or asymmetric key.
	Export(ctx context.Context, input *ExportRequest) (*pangea.PangeaResponse[ExportResult], error)

	GetBulk(ctx context.Context, input *GetBulkRequest) (*pangea.PangeaResponse[GetBulkResult], error)

	// Base service methods
	pangea.BaseServicer
}

type vault struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &vault{
		BaseService: pangea.NewBaseService("vault", cfg),
	}
	return cli
}
