// go:build integration
package vault_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/vault"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Develop
)

func PrintPangeAPIError(err error) {
	if err != nil {
		apiErr := err.(*pangea.APIError)
		fmt.Println(apiErr.Err.Error())
		for _, ef := range apiErr.PangeaErrors.Errors {
			fmt.Println(ef.Detail)
		}
	}
}

func Test_Integration_SecretLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	const (
		secretV1 = "mysecret"
		secretV2 = "newsecret"
	)

	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &vault.SecretStoreRequest{
		Secret: secretV1,
	}
	rStore, err := client.SecretStore(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.Equal(t, secretV1, rStore.Result.Secret)
	assert.NotEmpty(t, rStore.Result.ID)
	assert.Equal(t, 1, rStore.Result.Version)
	assert.Equal(t, string(vault.ITsecret), rStore.Result.Type)

	ID := rStore.Result.ID
	rRotate, err := client.SecretRotate(ctx,
		&vault.SecretRotateRequest{
			CommonRotateRequest: vault.CommonRotateRequest{
				ID: ID,
			},
			Secret: secretV2,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rRotate)
	assert.NotNil(t, rRotate.Result)
	assert.Equal(t, secretV2, rRotate.Result.Secret)
	assert.Equal(t, 2, rRotate.Result.Version)
	assert.Equal(t, string(vault.ITsecret), rRotate.Result.Type)

	rGet, err := client.Get(ctx,
		&vault.GetRequest{
			ID: ID,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.NotNil(t, rGet.Result.Secret)
	assert.Equal(t, secretV2, *rGet.Result.Secret)
	assert.Nil(t, rGet.Result.PrivateKey)
	assert.Nil(t, rGet.Result.PublicKey)
	assert.Nil(t, rGet.Result.Key)
	assert.Empty(t, rGet.Result.RevokedAt)
	assert.Equal(t, string(vault.ITsecret), rGet.Result.Type)

	rRevoke, err := client.Revoke(ctx,
		&vault.RevokeRequest{
			ID: ID,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rRevoke)
	assert.NotNil(t, rRevoke.Result)
	assert.Equal(t, ID, rRevoke.Result.ID)

	rGet, err = client.Get(ctx,
		&vault.GetRequest{
			ID: ID,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.NotNil(t, rGet.Result.Secret)
	assert.Equal(t, secretV2, *rGet.Result.Secret)
	assert.Nil(t, rGet.Result.PrivateKey)
	assert.Nil(t, rGet.Result.PublicKey)
	assert.Nil(t, rGet.Result.Key)
	assert.Empty(t, rGet.Result.RevokedAt)
	assert.Equal(t, string(vault.ITsecret), rGet.Result.Type)
}

func AsymSigningCycle(t *testing.T, client vault.Client, ctx context.Context, id string) {
	data := "thisisamessagetosign"

	// Sign 1
	rSign1, err := client.Sign(ctx,
		&vault.SignRequest{
			ID:      id,
			Message: data,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rSign1)
	assert.NotNil(t, rSign1.Result)
	assert.Equal(t, 1, rSign1.Result.Version)
	assert.NotEmpty(t, rSign1.Result.Signature)

	rRotate, err := client.KeyRotate(ctx,
		&vault.KeyRotateRequest{
			CommonRotateRequest: vault.CommonRotateRequest{
				ID: id,
			},
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rRotate)
	assert.NotNil(t, rRotate.Result)
	assert.NotNil(t, rRotate.Result.PublicKey)
	assert.Equal(t, 2, rRotate.Result.Version)
	assert.Equal(t, id, rRotate.Result.ID)

	// Sign 2
	rSign2, err := client.Sign(ctx,
		&vault.SignRequest{
			ID:      id,
			Message: data,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rSign2)
	assert.NotNil(t, rSign2.Result)
	assert.Equal(t, 2, rSign2.Result.Version)
	assert.NotEmpty(t, rSign2.Result.Signature)

	// Verify 2
	rVerify2, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign2.Result.Signature,
			Version:   pangea.Int(2),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerify2)
	assert.NotNil(t, rVerify2.Result)
	assert.True(t, rVerify2.Result.ValidSignature)

	// Verify 1
	rVerify1, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign1.Result.Signature,
			Version:   pangea.Int(1),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerify1)
	assert.NotNil(t, rVerify1.Result)
	assert.True(t, rVerify1.Result.ValidSignature)

	// Verify default
	rVerifyDef, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign2.Result.Signature,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerifyDef)
	assert.NotNil(t, rVerifyDef.Result)
	assert.True(t, rVerifyDef.Result.ValidSignature)

	// Verify wrong ID
	rVerifyBad1, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        "notavalidid",
			Message:   data,
			Signature: rSign2.Result.Signature,
		},
	)

	assert.Error(t, err)
	assert.Nil(t, rVerifyBad1)
	_ = err.(*pangea.APIError)

	// Verify wrong signature
	rVerifyBad2, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: "thisisnotasignature",
		},
	)

	assert.Error(t, err)
	assert.Nil(t, rVerifyBad2)
	_ = err.(*pangea.APIError)

	// Verify wrong pair signature/version
	rVerifyBad3, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign2.Result.Signature,
			Version:   pangea.Int(1),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerifyBad3)
	assert.NotNil(t, rVerifyBad3.Result)
	assert.False(t, rVerifyBad3.Result.ValidSignature)

	// Verify wrong data
	rVerifyBad4, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   "thisisnottheoriginaldata",
			Signature: rSign2.Result.Signature,
			Version:   pangea.Int(2),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerifyBad4)
	assert.NotNil(t, rVerifyBad4.Result)
	assert.False(t, rVerifyBad4.Result.ValidSignature)

	rRevoke, err := client.Revoke(ctx,
		&vault.RevokeRequest{
			ID: id,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rRevoke)
	assert.NotNil(t, rRevoke.Result)

	// Verify Revoked 2
	rVerifyRevoked2, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign2.Result.Signature,
			Version:   pangea.Int(2),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerifyRevoked2)
	assert.NotNil(t, rVerifyRevoked2.Result)
	assert.True(t, rVerifyRevoked2.Result.ValidSignature)

	// Verify Revoked 1
	rVerifyRevoked1, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign1.Result.Signature,
			Version:   pangea.Int(1),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerifyRevoked1)
	assert.NotNil(t, rVerifyRevoked1.Result)
	assert.True(t, rVerifyRevoked1.Result.ValidSignature)

}

func JWTSigningCycle(t *testing.T, client vault.Client, ctx context.Context, id string) {
	data := map[string]string{
		"message": "message to sign",
		"data":    "Some extra data",
	}

	b, err := json.Marshal(data)
	assert.NoError(t, err)

	payload := string(b)

	// Sign 1
	rSign1, err := client.JWTSign(ctx,
		&vault.JWTSignRequest{
			ID:      id,
			Payload: payload,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rSign1)
	assert.NotNil(t, rSign1.Result)
	assert.NotEmpty(t, rSign1.Result.JWS)

	rRotate, err := client.KeyRotate(ctx,
		&vault.KeyRotateRequest{
			CommonRotateRequest: vault.CommonRotateRequest{
				ID: id,
			},
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rRotate)
	assert.NotNil(t, rRotate.Result)
	assert.Equal(t, 2, rRotate.Result.Version)
	assert.Equal(t, id, rRotate.Result.ID)

	// Sign 2
	rSign2, err := client.JWTSign(ctx,
		&vault.JWTSignRequest{
			ID:      id,
			Payload: payload,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rSign2)
	assert.NotNil(t, rSign2.Result)
	assert.NotEmpty(t, rSign2.Result.JWS)

	// Verify 2
	rVerify2, err := client.JWTVerify(ctx,
		&vault.JWTVerifyRequest{
			JWS: rSign2.Result.JWS,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerify2)
	assert.NotNil(t, rVerify2.Result)
	assert.True(t, rVerify2.Result.ValidSignature)

	// Verify 1
	rVerify1, err := client.JWTVerify(ctx,
		&vault.JWTVerifyRequest{
			JWS: rSign1.Result.JWS,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerify1)
	assert.NotNil(t, rVerify1.Result)
	assert.True(t, rVerify1.Result.ValidSignature)

	// Get default
	rGet, err := client.JWTGet(ctx,
		&vault.JWTGetRequest{
			ID: id,
		},
	)
	PrintPangeAPIError(err)
	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 1, len(rGet.Result.JWK.Keys))

	// Get version 1
	rGet, err = client.JWTGet(ctx,
		&vault.JWTGetRequest{
			ID:      id,
			Version: pangea.String("1"),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 1, len(rGet.Result.JWK.Keys))

	// Get all
	rGet, err = client.JWTGet(ctx,
		&vault.JWTGetRequest{
			ID:      id,
			Version: pangea.String("all"),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 2, len(rGet.Result.JWK.Keys))

	// Get version -1
	rGet, err = client.JWTGet(ctx,
		&vault.JWTGetRequest{
			ID:      id,
			Version: pangea.String("-1"),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 2, len(rGet.Result.JWK.Keys))

	rRevoke, err := client.Revoke(ctx,
		&vault.RevokeRequest{
			ID: id,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rRevoke)
	assert.NotNil(t, rRevoke.Result)

	// Verify Revoked 2
	rVerifyRevoked2, err := client.JWTVerify(ctx,
		&vault.JWTVerifyRequest{
			JWS: rSign2.Result.JWS,
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rVerifyRevoked2)
	assert.NotNil(t, rVerifyRevoked2.Result)
	assert.True(t, rVerifyRevoked2.Result.ValidSignature)

}

func EncryptionCycle(t *testing.T, client vault.Client, ctx context.Context, id string) {
	const (
		msg = "thisisamessagetoencrypt"
	)
	dataB64 := pangea.StrToB64(msg)

	// Encode 1
	rEnc1, err := client.Encrypt(ctx,
		&vault.EncryptRequest{
			ID:        id,
			PlainText: dataB64,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rEnc1)
	assert.NotNil(t, rEnc1.Result)
	assert.NotEmpty(t, rEnc1.Result.CipherText)
	assert.Equal(t, 1, rEnc1.Result.Version)

	// Rotate 1
	rRot1, err := client.KeyRotate(ctx,
		&vault.KeyRotateRequest{
			CommonRotateRequest: vault.CommonRotateRequest{
				ID: id,
			},
		})

	assert.NoError(t, err)
	assert.NotNil(t, rRot1)
	assert.NotNil(t, rRot1.Result)
	assert.Equal(t, 2, rRot1.Result.Version)

	// Encode 2
	rEnc2, err := client.Encrypt(ctx,
		&vault.EncryptRequest{
			ID:        id,
			PlainText: dataB64,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rEnc2)
	assert.NotNil(t, rEnc2.Result)
	assert.NotEmpty(t, rEnc2.Result.CipherText)
	assert.Equal(t, 2, rEnc2.Result.Version)

	// Decrypt 1
	rDec1, err := client.Decrypt(ctx,
		&vault.DecryptRequest{
			ID:         id,
			CipherText: rEnc1.Result.CipherText,
			Version:    pangea.Int(1),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rDec1)
	assert.NotNil(t, rDec1.Result)
	assert.Equal(t, dataB64, rDec1.Result.PlainText)

	// Decrypt 2
	rDec2, err := client.Decrypt(ctx,
		&vault.DecryptRequest{
			ID:         id,
			CipherText: rEnc2.Result.CipherText,
			Version:    pangea.Int(2),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rDec2)
	assert.NotNil(t, rDec2.Result)
	assert.Equal(t, dataB64, rDec2.Result.PlainText)

	// Decrypt default
	rDecDef, err := client.Decrypt(ctx,
		&vault.DecryptRequest{
			ID:         id,
			CipherText: rEnc2.Result.CipherText,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rDecDef)
	assert.NotNil(t, rDecDef.Result)
	assert.Equal(t, dataB64, rDecDef.Result.PlainText)

	// Decrypt wrong version
	rDecBad1, err := client.Decrypt(ctx,
		&vault.DecryptRequest{
			ID:         id,
			CipherText: rEnc1.Result.CipherText,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rDecBad1)
	assert.NotNil(t, rDecBad1.Result)
	assert.NotEqual(t, dataB64, rDecBad1.Result.PlainText)

	// Error not and ID
	resp, err := client.Decrypt(ctx,
		&vault.DecryptRequest{
			ID:         "notanid",
			CipherText: rEnc2.Result.CipherText,
		})

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Revoke key
	rRevoke, err := client.Revoke(ctx,
		&vault.RevokeRequest{
			ID: id,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rRevoke)
	assert.NotNil(t, rRevoke.Result)
	assert.Equal(t, id, rRevoke.Result.ID)

	// Decrypt after revoked
	rDecRev1, err := client.Decrypt(ctx,
		&vault.DecryptRequest{
			ID:         id,
			CipherText: rEnc1.Result.CipherText,
			Version:    pangea.Int(1),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rDecRev1)
	assert.NotNil(t, rDecRev1.Result)
	assert.Equal(t, dataB64, rDecRev1.Result.PlainText)

}

func Test_Integration_Ed25519SigningLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.AAed25519
	purpose := vault.KPsigning

	// Generate
	rGen, err := client.AsymmetricGenerate(ctx,
		&vault.AsymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Managed: pangea.Bool(false),
				Store:   pangea.Bool(false),
			},
			Algorithm: &algorithm,
			Purpose:   &purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.PublicKey)
	assert.NotNil(t, rGen.Result.PublicKey)
	assert.NotEmpty(t, *rGen.Result.PrivateKey)
	assert.Empty(t, rGen.Result.ID)
	assert.Nil(t, rGen.Result.Version)

	// Store
	rStore, err := client.AsymmetricStore(ctx,
		&vault.AsymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{},
			Managed:            pangea.Bool(true),
			Algorithm:          algorithm,
			PublicKey:          rGen.Result.PublicKey,
			PrivateKey:         *rGen.Result.PrivateKey,
			Purpose:            &purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.NotEmpty(t, rStore.Result.PublicKey)
	assert.Nil(t, rStore.Result.PrivateKey)
	assert.NotEmpty(t, rStore.Result.ID)
	assert.Equal(t, 1, rStore.Result.Version)

	AsymSigningCycle(t, client, ctx, rStore.Result.ID)
}

func Test_Integration_JWT_ES256SigningLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.AAes256
	purpose := vault.KPjwt

	// Generate
	rGen, err := client.AsymmetricGenerate(ctx,
		&vault.AsymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Managed: pangea.Bool(false),
				Store:   pangea.Bool(false),
			},
			Algorithm: &algorithm,
			Purpose:   &purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.PublicKey)
	assert.NotNil(t, rGen.Result.PublicKey)
	assert.NotEmpty(t, *rGen.Result.PrivateKey)
	assert.Empty(t, rGen.Result.ID)
	assert.Nil(t, rGen.Result.Version)

	// Store
	rStore, err := client.AsymmetricStore(ctx,
		&vault.AsymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{},
			Managed:            pangea.Bool(true),
			Algorithm:          algorithm,
			PublicKey:          rGen.Result.PublicKey,
			PrivateKey:         *rGen.Result.PrivateKey,
			Purpose:            &purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.NotEmpty(t, rStore.Result.PublicKey)
	assert.Nil(t, rStore.Result.PrivateKey)
	assert.NotEmpty(t, rStore.Result.ID)
	assert.Equal(t, 1, rStore.Result.Version)

	JWTSigningCycle(t, client, ctx, rStore.Result.ID)
}

func Test_Integration_JWT_HS256SigningLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.SYAhs256
	purpose := vault.KPjwt

	rGen, err := client.SymmetricGenerate(ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Managed: pangea.Bool(false),
				Store:   pangea.Bool(false),
			},
			Algorithm: &algorithm,
			Purpose:   &purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.Empty(t, rGen.Result.ID)
	assert.NotNil(t, rGen.Result.Key)
	assert.NotEmpty(t, rGen.Result.Key)

	rStore, err := client.SymmetricStore(ctx,
		&vault.SymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{},
			Managed:            pangea.Bool(true),
			Key:                vault.EncodedSymmetricKey(*rGen.Result.Key),
			Algorithm:          algorithm,
		})

	PrintPangeAPIError(err)

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.NotEmpty(t, rStore.Result.ID)
	JWTSigningCycle(t, client, ctx, rStore.Result.ID)
}

func Test_Integration_AESencryptingLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.SYAaes

	rGen, err := client.SymmetricGenerate(ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Managed: pangea.Bool(false),
				Store:   pangea.Bool(false),
			},
			Algorithm: &algorithm,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.Empty(t, rGen.Result.ID)
	assert.NotNil(t, rGen.Result.Key)
	assert.NotEmpty(t, rGen.Result.Key)

	rStore, err := client.SymmetricStore(ctx,
		&vault.SymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{},
			Managed:            pangea.Bool(true),
			Key:                vault.EncodedSymmetricKey(*rGen.Result.Key),
			Algorithm:          algorithm,
		})

	PrintPangeAPIError(err)

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.NotEmpty(t, rStore.Result.ID)
	EncryptionCycle(t, client, ctx, rStore.Result.ID)
}

// 	const respStore = await vault.symmetricStore(algorithm, String(respGen.result.key));

// 	const id = respStore.result.id;
// 	expect(id).toBeDefined();
// 	expect(respStore.result.version).toBe(1);
// 	await encryptingCycle(id);
//   });

func Test_Integration_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	client := vault.New(cfg)

	input := &vault.SecretStoreRequest{
		Secret: "somesecret",
	}

	out, err := client.SecretStore(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}
