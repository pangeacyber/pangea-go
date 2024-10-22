package vault

import (
	"context"
	"errors"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

// @summary State change
//
// @description Change the state of a specific version of a secret or key.
//
// @operationId vault_post_v2_state_change
//
// @example
//
//	input := &vault.StateChangeRequest{
//		ID:    pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		State: vault.IVSdeactivated,
//	}
//
//	scr, err := vaultcli.StateChange(ctx, input)
func (v *vault) StateChange(ctx context.Context, input *StateChangeRequest) (*pangea.PangeaResponse[StateChangeResult], error) {
	return request.DoPost(ctx, v.Client, "v2/state/change", input, &StateChangeResult{})
}

// @summary Delete
//
// @description Delete a secret or key.
//
// @operationId vault_post_v2_delete
//
// @example
//
//	input := &vault.DeleteRequest{
//		ID: pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//	}
//
//	dr, err := vaultcli.Delete(ctx, input)
func (v *vault) Delete(ctx context.Context, input *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error) {
	return request.DoPost(ctx, v.Client, "v2/delete", input, &DeleteResult{})
}

// @summary Retrieve
//
// @description Retrieve a secret or key, and any associated information.
//
// @operationId vault_post_v2_get
//
// @example
//
//	input := &vault.GetRequest{
//		ID:           pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		Version:      pangea.StringValue(1),
//		Verbose:      pangea.Bool(true),
//		VersionState: &vault.IVSactive,
//	}
//
//	gr, err := vaultcli.Get(ctx, input)
func (v *vault) Get(ctx context.Context, input *GetRequest) (*pangea.PangeaResponse[GetResult], error) {
	return request.DoPost(ctx, v.Client, "v2/get", input, &GetResult{})
}

// @summary Get Bulk
//
// @description Retrieve a list of secrets, keys and folders.
//
// @operationId vault_post_v2_get_bulk
//
// @example
//
// filter.Folder().Set(pangea.String("/"))
// gbr, err := client.GetBulk(
//
//	ctx,
//	&vault.GetBulkRequest{
//		Filter: filter.Filter(),
//		Size:   5,
//	},
//
// )
func (v *vault) GetBulk(ctx context.Context, input *GetBulkRequest) (*pangea.PangeaResponse[GetBulkResult], error) {
	return request.DoPost(ctx, v.Client, "v2/get_bulk", input, &GetBulkResult{})
}

// @summary JWT Retrieve
//
// @description Retrieve a key in JWK format.
//
// @operationId vault_post_v2_jwk_get
//
// @example
//
//	input := &vault.JWKGetRequest{
//		ID: pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//	}
//
//	jr, err := vaultcli.JWKGet(ctx, input)
func (v *vault) JWKGet(ctx context.Context, input *JWKGetRequest) (*pangea.PangeaResponse[JWKGetResult], error) {
	return request.DoPost(ctx, v.Client, "v2/jwk/get", input, &JWKGetResult{})
}

// @summary List
//
// @description Retrieve a list of secrets, keys and folders, and their associated information.
//
// @operationId vault_post_v2_list
//
// @example
//
//	input := &vault.ListRequest{
//		Filter: map[string]string{
//			"folder": "/",
//			"type": "asymmetric_key",
//			"name__contains": "test",
//			"metadata_key1": "value1",
//			"created_at__lt": "2023-12-12T00:00:00Z",
//		},
//		Last:    pangea.StringValue("WyIvdGVzdF8yMDdfc3ltbWV0cmljLyJd"),
//		Size:    pangea.IntValue(20),
//		Order:   vault.IOasc,
//		OrderBy: vault.IOBname,
//	}
//
//	lr, err := vaultcli.List(ctx, input)
func (v *vault) List(ctx context.Context, input *ListRequest) (*pangea.PangeaResponse[ListResult], error) {
	return request.DoPost(ctx, v.Client, "v2/list", input, &ListResult{})
}

// @summary Update
//
// @description Update information associated with a secret or key.
//
// @operationId vault_post_v2_update
//
// @example
//
//	input := &vault.Update Request{
//		ID: pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		Name: pangea.StringValue("my-very-secret-secret"),
//		Folder: pangea.StringValue("/personal"),
//		Metadata: vault.Metadata{
//			"created_by": pangea.StringValue("John Doe"),
//			"used_in":    pangea.StringValue("Google products"),
//		},
//		Tags: vault.Tags{
//			pangea.StringValue("irs_2023"),
//			pangea.StringValue("personal"),
//		},
//		DisabledAt:          pangea.StringValue("2025-01-01T10:00:00Z"),
//		ItemState:           vault.ISdisabled,
//	}
//
//	ur, err := vaultcli.Update(ctx, input)
func (v *vault) Update(ctx context.Context, input *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error) {
	return request.DoPost(ctx, v.Client, "v2/update", input, &UpdateResult{})
}

// @summary Secret store
//
// @description Import a secret
//
// @operationId vault_post_v2_secret_store 1
//
// @example
//
//	input := &vault.SecretStoreRequest{
//		Secret: pangea.StringValue("12sdfgs4543qv@#%$casd"),
//		Type: vault.ITsecret
//		CommonStoreRequest: vault.CommonStoreRequest{
//			Name: pangea.StringValue("my-very-secret-secret"),
//			Folder: pangea.StringValue("/personal"),
//			Metadata: vault.Metadata{
//				"created_by": pangea.StringValue("John Doe"),
//				"used_in":    pangea.StringValue("Google products"),
//			},
//			Tags: vault.Tags{
//				pangea.StringValue("irs_2023"),
//				pangea.StringValue("personal"),
//			},
//		},
//	}
//
//	ssr, err := vaultcli.SecretStore(ctx, input)
func (v *vault) SecretStore(ctx context.Context, input *SecretStoreRequest) (*pangea.PangeaResponse[SecretStoreResult], error) {
	return request.DoPost(ctx, v.Client, "v2/secret/store", input, &SecretStoreResult{})
}

// @summary Secret rotate
//
// @description Rotate a secret.
//
// @operationId vault_post_v2_secret_rotate 1
//
// @example
//
//	input := &vault.SecretRotateRequest{
//		Secret: pangea.StringValue("12sdfgs4543qv@#%$casd"),
//		CommonRotateRequest: vault.CommonRotateRequest{
//			ID:           pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//			RotationState vault.IVSdeactivated,
//		},
//	}
//
//	srr, err := vaultcli.SecretRotate(ctx, input)
func (v *vault) SecretRotate(ctx context.Context, input *SecretRotateRequest) (*pangea.PangeaResponse[SecretRotateResult], error) {
	return request.DoPost(ctx, v.Client, "v2/secret/rotate", input, &SecretRotateResult{})
}

// @summary Symmetric generate
//
// @description Generate a symmetric key.
//
// @operationId vault_post_v2_key_generate 2
//
// @example
//
//	input := &vault.SymmetricGenerateRequest{
//		Algorithm: vault.SYAaes128_cfb,
//		Purpose:   vault.KPencryption,
//		CommonGenerateRequest: vault.CommonGenerateRequest{
//			Name:   pangea.StringValue("my-very-secret-secret"),
//			Folder: pangea.StringValue("/personal"),
//			Metadata: vault.Metadata{
//				"created_by": pangea.StringValue("John Doe"),
//				"used_in":    pangea.StringValue("Google products"),
//			},
//			Tags: vault.Tags{
//				pangea.StringValue("irs_2023"),
//				pangea.StringValue("personal"),
//			},
//			RotationFrequency: pangea.StringValue("10d"),
//			RotationState:     pangea.StringValue("deactivated"),
//			DisabledAt:        pangea.StringValue("2025-01-01T10:00:00Z"),
//		},
//	}
//
//	sgr, err := vaultcli.SymmetricGenerate(ctx, input)
func (v *vault) SymmetricGenerate(ctx context.Context, input *SymmetricGenerateRequest) (*pangea.PangeaResponse[SymmetricGenerateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITsymmetricKey

	return request.DoPost(ctx, v.Client, "v2/key/generate", input, &SymmetricGenerateResult{})
}

// @summary Asymmetric generate
//
// @description Generate an asymmetric key.
//
// @operationId vault_post_v2_key_generate 1
//
// @example
//
//	input := &vault.AsymmetricGenerateRequest{
//		Algorithm: vault.AArsa2048_pkcs1v15_sha256,
//		Purpose:   vault.KPsigning,
//		CommonGenerateRequest: vault.CommonGenerateRequest{
//			Name:   pangea.StringValue("my-very-secret-secret"),
//			Folder: pangea.StringValue("/personal"),
//			Metadata: vault.Metadata{
//				"created_by": pangea.StringValue("John Doe"),
//				"used_in":    pangea.StringValue("Google products"),
//			},
//			Tags: vault.Tags{
//				pangea.StringValue("irs_2023"),
//				pangea.StringValue("personal"),
//			},
//			RotationFrequency: pangea.StringValue("10d"),
//			RotationState:     pangea.StringValue("deactivated"),
//			DisabledAt:        pangea.StringValue("2025-01-01T10:00:00Z"),
//		},
//	}
//
//	agr, err := vaultcli.AsymmetricGenerate(ctx, input)
func (v *vault) AsymmetricGenerate(ctx context.Context, input *AsymmetricGenerateRequest) (*pangea.PangeaResponse[AsymmetricGenerateResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITasymmetricKey

	return request.DoPost(ctx, v.Client, "v2/key/generate", input, &AsymmetricGenerateResult{})
}

// @summary Asymmetric store
//
// @description Import an asymmetric key.
//
// @operationId vault_post_v2_key_store 1
//
// @example
//
//	var PUBLIC_KEY = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEA8s5JopbEPGBylPBcMK+L5PqHMqPJW/5KYPgBHzZGncc=\n-----END PUBLIC KEY-----"
//	var PRIVATE_KEY = "private key example"
//
//	input := &vault.AsymmetricStoreRequest{
//		Algorithm:  vault.AArsa2048_pkcs1v15_sha256,
//		PublicKey:  vault.EncodedPublicKey(PUBLIC_KEY),
//		PrivateKey: vault.EncodedPrivateKey(PRIVATE_KEY),
//		Purpose:    vault.KPsigning,
//		CommonStoreRequest: vault.CommonStoreRequest{
//			Name: pangea.StringValue("my-very-secret-secret"),
//			Folder: pangea.StringValue("/personal"),
//			Metadata: vault.Metadata{
//				"created_by": pangea.StringValue("John Doe"),
//				"used_in":    pangea.StringValue("Google products"),
//			},
//			Tags: vault.Tags{
//				pangea.StringValue("irs_2023"),
//				pangea.StringValue("personal"),
//			},
//			RotationFrequency: pangea.StringValue("10d"),
//			RotationState:     pangea.StringValue("deactivated"),
//			DisabledAt:        pangea.StringValue("2025-01-01T10:00:00Z"),
//		},
//	}
//
//	asr, err := vaultcli.AsymmetricStore(ctx, input)
func (v *vault) AsymmetricStore(ctx context.Context, input *AsymmetricStoreRequest) (*pangea.PangeaResponse[AsymmetricStoreResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITasymmetricKey

	return request.DoPost(ctx, v.Client, "v2/key/store", input, &AsymmetricStoreResult{})
}

// @summary Symmetric store
//
// @description Import a symmetric key.
//
// @operationId vault_post_v2_key_store 2
//
// @example
//
//	input := &vault.SymmetricStoreRequest{
//		Key: vault.EncodedSymmetricKey("lJkk0gCLux+Q+rPNqLPEYw=="),
//		Algorithm: vault.SYAaes128_cfb,
//		Purpose: vault.KPencryption,
//		CommonStoreRequest: vault.CommonStoreRequest{
//			Name: pangea.StringValue("my-very-secret-secret"),
//			Folder: pangea.StringValue("/personal"),
//			Metadata: vault.Metadata{
//				"created_by": pangea.StringValue("John Doe"),
//				"used_in":    pangea.StringValue("Google products"),
//			},
//			Tags: vault.Tags{
//				pangea.StringValue("irs_2023"),
//				pangea.StringValue("personal"),
//			},
//			RotationFrequency: pangea.StringValue("10d"),
//			RotationState:     pangea.StringValue("deactivated"),
//			DisabledAt:        pangea.StringValue("2025-01-01T10:00:00Z"),
//		},
//	}
//
//	ssr, err := vaultcli.SymmetricStore(ctx, input)
func (v *vault) SymmetricStore(ctx context.Context, input *SymmetricStoreRequest) (*pangea.PangeaResponse[SymmetricStoreResult], error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}
	input.Type = ITsymmetricKey

	return request.DoPost(ctx, v.Client, "v2/key/store", input, &SymmetricStoreResult{})
}

// @summary Key rotate
//
// @description Manually rotate a symmetric or asymmetric key.
//
// @operationId vault_post_v2_key_rotate
//
// @example
//
//	var SYMMETRIC_KEY = "lJkk0gCLux+Q+rPNqLPEYw=="
//
//	input := &vault.KeyRotateRequest{
//		CommonRotateRequest: vault.CommonRotateRequest{
//			ID:            pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//			RotationState: vault.IVSdeactivated,
//		},
//		Key: &vault.EncodedSymmetricKey(SYMMETRIC_KEY),
//	}
//
//	krr, err := vaultcli.KeyRotate(ctx, input)
func (v *vault) KeyRotate(ctx context.Context, input *KeyRotateRequest) (*pangea.PangeaResponse[KeyRotateResult], error) {
	return request.DoPost(ctx, v.Client, "v2/key/rotate", input, &KeyRotateResult{})
}

// @summary Encrypt
//
// @description Encrypt a message using a key.
//
// @operationId vault_post_v2_encrypt
//
// @example
//
//	msg := "message to encrypt..."
//	data := base64.StdEncoding.EncodeToString([]byte(msg))
//
//	input := &vault.EncryptRequest{
//		ID: pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		PlainText: data,
//		Version: pangea.Int(1),
//	}
//
//	enc, err := vaultcli.Encrypt(ctx, input)
func (v *vault) Encrypt(ctx context.Context, input *EncryptRequest) (*pangea.PangeaResponse[EncryptResult], error) {
	return request.DoPost(ctx, v.Client, "v2/encrypt", input, &EncryptResult{})
}

// @summary Decrypt
//
// @description Decrypt a message using a key.
//
// @operationId vault_post_v2_decrypt
//
// @example
//
//	input := &vault.DecryptRequest{
//		ID: pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		CipherText: pangea.StringValue("lJkk0gCLux+Q+rPNqLPEYw=="),
//		Version: pangea.Int(1),
//	}
//
//	dr, err := vaultcli.Decrypt(ctx, input)
func (v *vault) Decrypt(ctx context.Context, input *DecryptRequest) (*pangea.PangeaResponse[DecryptResult], error) {
	return request.DoPost(ctx, v.Client, "v2/decrypt", input, &DecryptResult{})
}

// @summary Sign
//
// @description Sign a message using a key.
//
// @operationId vault_post_v2_sign
//
// @example
//
//	msg := "message to sign..."
//	data := base64.StdEncoding.EncodeToString([]byte(msg))
//
//	input := &vault.SignRequest{
//		ID: pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		Message: data,
//		Version: pangea.Int(1),
//	}
//
//	sr, err := vaultcli.Sign(ctx, input)
func (v *vault) Sign(ctx context.Context, input *SignRequest) (*pangea.PangeaResponse[SignResult], error) {
	return request.DoPost(ctx, v.Client, "v2/sign", input, &SignResult{})
}

// @summary Verify
//
// @description Verify a signature using a key.
//
// @operationId vault_post_v2_verify
//
// @example
//
//	input := &vault.VerifyRequest{
//		ID:        pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		Version:   pangea.Int(1),
//		Message:   pangea.StringValue("lJkk0gCLux+Q+rPNqLPEYw=="),
//		Signature: pangea.StringValue("FfWuT2Mq/+cxa7wIugfhzi7ktZxVf926idJNgBDCysF/knY9B7M6wxqHMMPDEBs86D8OsEGuED21y3J7IGOpCQ=="),
//	}
//
//	vr, err := vaultcli.Verify(ctx, input)
func (v *vault) Verify(ctx context.Context, input *VerifyRequest) (*pangea.PangeaResponse[VerifyResult], error) {
	return request.DoPost(ctx, v.Client, "v2/verify", input, &VerifyResult{})
}

// @summary JWT Sign
//
// @description Sign a JSON Web Token (JWT) using a key.
//
// @operationId vault_post_v2_jwt_sign
//
// @example
//
//	input := &vault.JWTSignRequest{
//		ID:      pangea.StringValue("pvi_p6g5i3gtbvqvc3u6zugab6qs6r63tqf5"),
//		Payload: pangea.StringValue("{\"sub\": \"1234567890\",\"name\": \"John Doe\",\"admin\": true}"),
//	}
//
//	jr, err := vaultcli.JWTSign(ctx, input)
func (v *vault) JWTSign(ctx context.Context, input *JWTSignRequest) (*pangea.PangeaResponse[JWTSignResult], error) {
	return request.DoPost(ctx, v.Client, "v2/jwt/sign", input, &JWTSignResult{})
}

// @summary JWT Verify
//
// @description Verify the signature of a JSON Web Token (JWT).
//
// @operationId vault_post_v2_jwt_verify
//
// @example
//
//	input := &vault.JWTVerifyRequest{
//		JWS: pangea.StringValue("ewogICJhbGciO..."),
//	}
//
//	jr, err := vaultcli.JWTVerify(ctx, input)
func (v *vault) JWTVerify(ctx context.Context, input *JWTVerifyRequest) (*pangea.PangeaResponse[JWTVerifyResult], error) {
	return request.DoPost(ctx, v.Client, "v2/jwt/verify", input, &JWTVerifyResult{})
}

// @summary Create
//
// @description Creates a folder.
//
// @operationId vault_post_v2_folder_create
//
// @example
//
//	input := &vault.FolderCreateRequest{
//	 	Name:   "folder_name",
//	 	Folder: "parent/folder/name",
//	}
//
//	enc, err := vaultcli.FolderCreate(ctx, input)
func (v *vault) FolderCreate(ctx context.Context, input *FolderCreateRequest) (*pangea.PangeaResponse[FolderCreateResult], error) {
	return request.DoPost(ctx, v.Client, "v2/folder/create", input, &FolderCreateResult{})
}

// @summary Encrypt structured
//
// @description Encrypt parts of a JSON object.
//
// @operationId vault_post_v2_encrypt_structured
//
// @example
//
//	data := map[string]interface{}{
//		"field1": [4]interface{}{1, 2, "true", "false"},
//		"field2": "true",
//	}
//
//	encryptedResponse, err := client.EncryptStructured(
//		ctx,
//		&vault.EncryptStructuredRequest{
//			ID:             "pvi_[...]",
//			StructuredData: data,
//			Filter:         "$.field1[2:4]",
//		},
//	)
func (v *vault) EncryptStructured(ctx context.Context, input *EncryptStructuredRequest) (*pangea.PangeaResponse[EncryptStructuredResult], error) {
	return request.DoPost(ctx, v.Client, "v2/encrypt_structured", input, &EncryptStructuredResult{})
}

// @summary Decrypt structured
//
// @description Decrypt parts of a JSON object.
//
// @operationId vault_post_v2_decrypt_structured
//
// @example
//
//	data := map[string]interface{}{
//		"field1": [4]interface{}{1, 2, "kxcbC9E9IlgVaSCChPWUMgUC3ko=", "6FfI/LCzatLRLNAc8SuBK/TDnGxp"},
//		"field2": "true",
//	}
//
//	decryptedResponse, err := client.DecryptStructured(
//		ctx,
//		&vault.EncryptStructuredRequest{
//			ID:             "pvi_[...]",
//			StructuredData: data,
//			Filter:         "$.field1[2:4]",
//		},
//	)
func (v *vault) DecryptStructured(ctx context.Context, input *EncryptStructuredRequest) (*pangea.PangeaResponse[EncryptStructuredResult], error) {
	return request.DoPost(ctx, v.Client, "v2/decrypt_structured", input, &EncryptStructuredResult{})
}

// @summary Encrypt transform
//
// @description Encrypt using a format-preserving algorithm (FPE).
//
// @operationId vault_post_v2_encrypt_transform
//
// @example
//
//	encryptedResponse, err := client.EncryptTransform(
//		ctx,
//		&vault.EncryptTransformRequest{
//			ID:        "pvi_[...]",
//			PlainText: "123-4567-8901",
//			Tweak:     "MTIzMTIzMT==",
//			Alphabet:  vault.TAalphanumeric,
//		},
//	)
func (v *vault) EncryptTransform(ctx context.Context, input *EncryptTransformRequest) (*pangea.PangeaResponse[EncryptTransformResult], error) {
	return request.DoPost(ctx, v.Client, "v2/encrypt_transform", input, &EncryptTransformResult{})
}

// @summary Decrypt transform
//
// @description Decrypt using a format-preserving algorithm (FPE).
//
// @operationId vault_post_v2_decrypt_transform
//
// @example
//
//	decryptedResponse, err := client.DecryptTransform(
//		ctx,
//		&vault.DecryptTransformRequest{
//			ID:         "pvi_[...]",
//			CipherText: "tZB-UKVP-MzTM",
//			Tweak:      "MTIzMTIzMT==",
//			Alphabet:   vault.TAalphanumeric,
//		},
//	)
func (v *vault) DecryptTransform(ctx context.Context, input *DecryptTransformRequest) (*pangea.PangeaResponse[DecryptTransformResult], error) {
	return request.DoPost(ctx, v.Client, "v2/decrypt_transform", input, &DecryptTransformResult{})
}

// @summary Export
//
// @description Export a symmetric or asymmetric key.
//
// @operationId vault_post_v2_export
//
// @example
//
//	// Generate an exportable key.
//	generated, err := client.AsymmetricGenerate(ctx,
//		&vault.AsymmetricGenerateRequest{
//			CommonGenerateRequest: vault.CommonGenerateRequest{
//				Name: "a-name-for-the-key",
//			},
//			Algorithm:  vault.AAed25519,
//			Purpose:    vault.KPsigning,
//			Exportable: pangea.Bool(true),
//		})
//
//	// Then it can be exported whenever needed.
//	exported, err := client.Export(ctx, &vault.ExportRequest{id: generated.Result.ID})
func (v *vault) Export(ctx context.Context, input *ExportRequest) (*pangea.PangeaResponse[ExportResult], error) {
	return request.DoPost(ctx, v.Client, "v2/export", input, &ExportResult{})
}

// @summary Encrypt structured with FPE
//
// @description Encrypt using a format preserving algorithm (FPE) parts of a JSON object.
//
// @operationId vault_post_v2_encrypt_transform_structured
//
// @example
//
//	data := map[string]interface{}{
//		"field1": [4]interface{}{"123-4567-8901", 2, "123-4567-8901", "false"},
//		"field2": "true",
//	}
//
// tweak := "MTIzMTIzMT=="
// encryptedResponse, err := client.EncryptTransformStructured(
//
//	ctx,
//	&vault.EncryptTransformStructuredRequest{
//		ID:             key,
//		StructuredData: data,
//		Filter:         "$.field1[2:4]",
//		Tweak:          &tweak,
//		Alphabet:       vault.TAalphanumeric,
//	},
//
// )
func (v *vault) EncryptTransformStructured(ctx context.Context, input *EncryptTransformStructuredRequest) (*pangea.PangeaResponse[EncryptTransformStructuredResult], error) {
	return request.DoPost(ctx, v.Client, "v2/encrypt_transform_structured", input, &EncryptTransformStructuredResult{})
}

// @summary Decrypt structured with FPE
//
// @description Decrypt using a format preserving algorithm (FPE) parts of a JSON object.
//
// @operationId vault_post_v2_decrypt_transform_structured
//
// @example
// tweak := "MTIzMTIzMT=="
//
//	data := map[string]interface{}{
//		"field1": [4]interface{}{"123-4567-8901",2,"Dbw-618a-KAty","o3ZaE"},
//		"field2": "true",
//	}
//
// decryptedResponse, err := client.DecryptTransformStructured(
//
//	ctx,
//	&vault.EncryptTransformStructuredRequest{
//		ID:             key,
//		StructuredData: encryptedResponse.Result.StructuredData,
//		Filter:         "$.field1[2:4]",
//		Tweak:          &tweak,
//		Alphabet:       vault.TAalphanumeric,
//	},
//
// )
func (v *vault) DecryptTransformStructured(ctx context.Context, input *EncryptTransformStructuredRequest) (*pangea.PangeaResponse[EncryptTransformStructuredResult], error) {
	return request.DoPost(ctx, v.Client, "v2/decrypt_transform_structured", input, &EncryptTransformStructuredResult{})
}
