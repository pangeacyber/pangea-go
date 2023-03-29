package vault

import (
	"context"
	"errors"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// @summary State change
//
// @description Change the state of a specific version of a secret or key
//
// @example
//
//	input := &vault.StateChangeRequest{
//		ID:    pangea.String("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		State: vault.IVSdeactivated,
//	}
//
//	stateChangeResponse, err := vaultcli.StateChange(ctx, input)
//
func (v *Vault) StateChange(ctx context.Context, input *StateChangeRequest) (*pangea.PangeaResponse[StateChangeResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/state/change", input)
	if err != nil {
		return nil, err
	}
	out := StateChangeResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[StateChangeResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @summary Delete
//
// @description Delete a secret or key
//
// @example
//
//	input := &vault.DeleteRequest{
//		ID: pangea.String("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//	}
//
//	deleteResponse, err := vaultcli.Delete(ctx, input)
//
func (v *Vault) Delete(ctx context.Context, input *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/delete", input)
	if err != nil {
		return nil, err
	}
	out := DeleteResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[DeleteResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @summary Retrieve
//
// @description Retrieve a secret or key, and any associated information
//
// @example
//
//	input := &vault.GetRequest{
//		ID:           pangea.String("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		Version:      pangea.Int(1),
//		Verbose:      pangea.Bool(true),
//		VersionState: vault.IVSactive,
//	}
//
//	getResponse, err := vaultcli.Get(ctx, input)
//
func (v *Vault) Get(ctx context.Context, input *GetRequest) (*pangea.PangeaResponse[GetResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/get", input)
	if err != nil {
		return nil, err
	}
	out := GetResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[GetResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) JWKGet(ctx context.Context, input *JWKGetRequest) (*pangea.PangeaResponse[JWKGetResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/get/jwk", input)
	if err != nil {
		return nil, err
	}
	out := JWKGetResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[JWKGetResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) List(ctx context.Context, input *ListRequest) (*pangea.PangeaResponse[ListResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/list", input)
	if err != nil {
		return nil, err
	}
	out := ListResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[ListResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) Update(ctx context.Context, input *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/update", input)
	if err != nil {
		return nil, err
	}
	out := UpdateResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UpdateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) SecretStore(ctx context.Context, input *SecretStoreRequest) (*pangea.PangeaResponse[SecretStoreResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITsecret

	req, err := v.Client.NewRequest("POST", "v1/secret/store", input)
	if err != nil {
		return nil, err
	}
	out := SecretStoreResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SecretStoreResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) PangeaTokenStore(ctx context.Context, input *PangeaTokenStoreRequest) (*pangea.PangeaResponse[SecretStoreResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITpangeaToken

	req, err := v.Client.NewRequest("POST", "v1/secret/store", input)
	if err != nil {
		return nil, err
	}
	out := SecretStoreResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SecretStoreResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) SecretRotate(ctx context.Context, input *SecretRotateRequest) (*pangea.PangeaResponse[SecretRotateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/secret/rotate", input)
	if err != nil {
		return nil, err
	}
	out := SecretRotateResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SecretRotateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) PangeaTokenRotate(ctx context.Context, input *PangeaTokenRotateRequest) (*pangea.PangeaResponse[SecretRotateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/secret/rotate", input)
	if err != nil {
		return nil, err
	}
	out := SecretRotateResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SecretRotateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) SymmetricGenerate(ctx context.Context, input *SymmetricGenerateRequest) (*pangea.PangeaResponse[SymmetricGenerateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITsymmetricKey

	req, err := v.Client.NewRequest("POST", "v1/key/generate", input)
	if err != nil {
		return nil, err
	}
	out := SymmetricGenerateResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SymmetricGenerateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) AsymmetricGenerate(ctx context.Context, input *AsymmetricGenerateRequest) (*pangea.PangeaResponse[AsymmetricGenerateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITasymmetricKey

	req, err := v.Client.NewRequest("POST", "v1/key/generate", input)
	if err != nil {
		return nil, err
	}
	out := AsymmetricGenerateResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[AsymmetricGenerateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) AsymmetricStore(ctx context.Context, input *AsymmetricStoreRequest) (*pangea.PangeaResponse[AsymmetricStoreResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITasymmetricKey

	req, err := v.Client.NewRequest("POST", "v1/key/store", input)
	if err != nil {
		return nil, err
	}
	out := AsymmetricStoreResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[AsymmetricStoreResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) SymmetricStore(ctx context.Context, input *SymmetricStoreRequest) (*pangea.PangeaResponse[SymmetricStoreResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITsymmetricKey

	req, err := v.Client.NewRequest("POST", "v1/key/store", input)
	if err != nil {
		return nil, err
	}
	out := SymmetricStoreResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SymmetricStoreResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) KeyRotate(ctx context.Context, input *KeyRotateRequest) (*pangea.PangeaResponse[KeyRotateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/rotate", input)
	if err != nil {
		return nil, err
	}
	out := KeyRotateResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[KeyRotateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) Encrypt(ctx context.Context, input *EncryptRequest) (*pangea.PangeaResponse[EncryptResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/encrypt", input)
	if err != nil {
		return nil, err
	}
	out := EncryptResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[EncryptResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) Decrypt(ctx context.Context, input *DecryptRequest) (*pangea.PangeaResponse[DecryptResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/decrypt", input)
	if err != nil {
		return nil, err
	}
	out := DecryptResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[DecryptResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) Sign(ctx context.Context, input *SignRequest) (*pangea.PangeaResponse[SignResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/sign", input)
	if err != nil {
		return nil, err
	}
	out := SignResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SignResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) Verify(ctx context.Context, input *VerifyRequest) (*pangea.PangeaResponse[VerifyResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/verify", input)
	if err != nil {
		return nil, err
	}
	out := VerifyResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[VerifyResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) JWTSign(ctx context.Context, input *JWTSignRequest) (*pangea.PangeaResponse[JWTSignResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/sign/jwt", input)
	if err != nil {
		return nil, err
	}
	out := JWTSignResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[JWTSignResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

func (v *Vault) JWTVerify(ctx context.Context, input *JWTVerifyRequest) (*pangea.PangeaResponse[JWTVerifyResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := v.Client.NewRequest("POST", "v1/key/verify/jwt", input)
	if err != nil {
		return nil, err
	}
	out := JWTVerifyResult{}
	resp, err := v.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[JWTVerifyResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}
