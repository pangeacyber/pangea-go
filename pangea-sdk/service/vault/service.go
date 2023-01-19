package vault

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Revoke(ctx context.Context, input *RevokeRequest) (*pangea.PangeaResponse[RevokeResult], error)
	Delete(ctx context.Context, input *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error)
	Get(ctx context.Context, input *GetRequest) (*pangea.PangeaResponse[GetResult], error)
	List(ctx context.Context, input *ListRequest) (*pangea.PangeaResponse[ListResult], error)
	Update(ctx context.Context, input *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error)
	SecretStore(ctx context.Context, input *SecretStoreRequest) (*pangea.PangeaResponse[SecretStoreResult], error)
	SecretRotate(ctx context.Context, input *SecretRotateRequest) (*pangea.PangeaResponse[SecretRotateResult], error)
	SymmetricGenerate(ctx context.Context, input *SymmetricGenerateRequest) (*pangea.PangeaResponse[SymmetricGenerateResult], error)
	AsymmetricGenerate(ctx context.Context, input *AsymmetricGenerateRequest) (*pangea.PangeaResponse[AsymmetricGenerateResult], error)
	SymmetricStore(ctx context.Context, input *SymmetricStoreRequest) (*pangea.PangeaResponse[SymmetricStoreResult], error)
	AsymmetricStore(ctx context.Context, input *AsymmetricStoreRequest) (*pangea.PangeaResponse[AsymmetricStoreResult], error)
	KeyRotate(ctx context.Context, input *KeyRotateRequest) (*pangea.PangeaResponse[KeyRotateResult], error)
	Encrypt(ctx context.Context, input *EncryptRequest) (*pangea.PangeaResponse[EncryptResult], error)
	Decrypt(ctx context.Context, input *DecryptRequest) (*pangea.PangeaResponse[DecryptResult], error)
	Sign(ctx context.Context, input *SignRequest) (*pangea.PangeaResponse[SignResult], error)
	Verify(ctx context.Context, input *VerifyRequest) (*pangea.PangeaResponse[VerifyResult], error)
}

type vault struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &vault{
		Client: pangea.NewClient("vault", cfg),
	}
	return cli
}
