// go:build integration
package vault_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/vault"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
	actor              = "GoSDKTest"
)

var timeNow = time.Now()
var timeStr = timeNow.Format("yyyyMMdd_HHmmss")
var KEY_ED25519_algorithm = vault.AAed25519
var KEY_ED25519_private_key = "-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIGthqegkjgddRAn0PWN2FeYC6HcCVQf/Ph9sUbeprTBO\n-----END PRIVATE KEY-----\n"
var KEY_ED25519_public_key = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAPlGrDliJXUbPc2YWEhFxlL2UbBfLHc3ed1f36FrDtTc=\n-----END PUBLIC KEY-----\n"
var KEY_AES_algorithm = vault.SYAaes
var KEY_AES_key = "oILlp2FUPHWiaqFXl4/1ww=="

func PrintPangeAPIError(err error) {
	if err != nil {
		apiErr := err.(*pangea.APIError)
		fmt.Println(apiErr.Err.Error())
		for _, ef := range apiErr.PangeaErrors.Errors {
			fmt.Println(ef.Detail)
		}
	}
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetName(name string) string {
	return fmt.Sprintf("%s_%s_%s_%s", timeStr, actor, name, GetRandID())
}

func GetRandID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprint(rand.Intn(1000000))
}

func Test_Integration_SecretLifeCycle(t *testing.T) {
	name := GetName("Test_Integration_SecretLifeCycle")
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	const (
		secretV1 = "mysecret"
		secretV2 = "newsecret"
	)

	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &vault.SecretStoreRequest{
		CommonStoreRequest: vault.CommonStoreRequest{
			Name: name,
		},
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
				ID:            ID,
				RotationState: vault.IVSsuspended,
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
	assert.Equal(t, 0, len(rGet.Result.Versions))
	assert.Equal(t, secretV2, *rGet.Result.CurrentVersion.Secret)
	assert.Equal(t, string(vault.IVSactive), rGet.Result.CurrentVersion.State)
	assert.Nil(t, rGet.Result.CurrentVersion.PublicKey)
	assert.Nil(t, rGet.Result.CurrentVersion.DestroyAt)
	assert.Equal(t, string(vault.ITsecret), rGet.Result.Type)

	rStateChange, err := client.StateChange(ctx,
		&vault.StateChangeRequest{
			ID:      ID,
			Version: pangea.Int(2),
			State:   vault.IVSsuspended,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rStateChange)
	assert.NotNil(t, rStateChange.Result)
	assert.Equal(t, ID, rStateChange.Result.ID)

	rGet, err = client.Get(ctx,
		&vault.GetRequest{
			ID: ID,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 0, len(rGet.Result.Versions))
	assert.Equal(t, secretV2, *rGet.Result.CurrentVersion.Secret)
	assert.Nil(t, rGet.Result.CurrentVersion.PublicKey)
	assert.Equal(t, string(vault.IVSsuspended), rGet.Result.CurrentVersion.State)
	assert.Nil(t, rGet.Result.CurrentVersion.DestroyAt)
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
				ID:            id,
				RotationState: vault.IVSsuspended,
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

	rStateChange, err := client.StateChange(ctx,
		&vault.StateChangeRequest{
			ID:      id,
			Version: pangea.Int(1),
			State:   vault.IVSdeactivated,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rStateChange)
	assert.NotNil(t, rStateChange.Result)

	// Verify deactivated 1
	rVerifyDeactivated1, err := client.Verify(ctx,
		&vault.VerifyRequest{
			ID:        id,
			Message:   data,
			Signature: rSign1.Result.Signature,
			Version:   pangea.Int(1),
		},
	)

	// FIXME: Should be an error
	assert.NoError(t, err)
	assert.NotNil(t, rVerifyDeactivated1)
	assert.NotNil(t, rVerifyDeactivated1.Result)
	assert.True(t, rVerifyDeactivated1.Result.ValidSignature)
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
				ID:            id,
				RotationState: vault.IVSsuspended,
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
	rGet, err := client.JWKGet(ctx,
		&vault.JWKGetRequest{
			ID: id,
		},
	)
	PrintPangeAPIError(err)
	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 1, len(rGet.Result.JWK.Keys))

	// Get version 1
	rGet, err = client.JWKGet(ctx,
		&vault.JWKGetRequest{
			ID:      id,
			Version: pangea.String("1"),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 1, len(rGet.Result.JWK.Keys))

	// Get all
	rGet, err = client.JWKGet(ctx,
		&vault.JWKGetRequest{
			ID:      id,
			Version: pangea.String("all"),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 2, len(rGet.Result.JWK.Keys))

	// Get version -1
	rGet, err = client.JWKGet(ctx,
		&vault.JWKGetRequest{
			ID:      id,
			Version: pangea.String("-1"),
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 2, len(rGet.Result.JWK.Keys))

	rStateChange, err := client.StateChange(ctx,
		&vault.StateChangeRequest{
			ID:      id,
			State:   vault.IVSsuspended,
			Version: pangea.Int(2),
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rStateChange)
	assert.NotNil(t, rStateChange.Result)

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
				ID:            id,
				RotationState: vault.IVSsuspended,
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
	rStateChange, err := client.StateChange(ctx,
		&vault.StateChangeRequest{
			ID:      id,
			Version: pangea.Int(2),
			State:   vault.IVSsuspended,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rStateChange)
	assert.NotNil(t, rStateChange.Result)
	assert.Equal(t, id, rStateChange.Result.ID)

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
	name := GetName("Test_Integration_Ed25519SigningLifeCycle")

	// Generate
	rGen, err := client.AsymmetricGenerate(ctx,
		&vault.AsymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: name,
			},
			Algorithm: algorithm,
			Purpose:   purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.PublicKey)
	assert.NotEmpty(t, rGen.Result.ID)
	assert.Equal(t, 1, rGen.Result.Version)

	AsymSigningCycle(t, client, ctx, rGen.Result.ID)
}

func Test_Integration_Ed25519StoreLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	purpose := vault.KPsigning
	name := GetName("Test_Integration_Ed25519StoreLifeCycle")

	// Store
	rStore, err := client.AsymmetricStore(ctx,
		&vault.AsymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{
				Name: name,
			},
			Algorithm:  KEY_ED25519_algorithm,
			Purpose:    purpose,
			PublicKey:  vault.EncodedPublicKey(KEY_ED25519_public_key),
			PrivateKey: vault.EncodedPrivateKey(KEY_ED25519_private_key),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.NotEmpty(t, rStore.Result.PublicKey)
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
	name := GetName("Test_Integration_JWT_ES256SigningLifeCycle")
	// Generate
	rGen, err := client.AsymmetricGenerate(ctx,
		&vault.AsymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: name,
			},
			Algorithm: algorithm,
			Purpose:   purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.PublicKey)
	assert.NotEmpty(t, rGen.Result.ID)
	assert.Equal(t, 1, rGen.Result.Version)

	JWTSigningCycle(t, client, ctx, rGen.Result.ID)
}

func Test_Integration_JWT_HS256SigningLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.SYAhs256
	purpose := vault.KPjwt
	name := GetName("Test_Integration_JWT_HS256SigningLifeCycle")

	rGen, err := client.SymmetricGenerate(ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: name,
			},
			Algorithm: algorithm,
			Purpose:   purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)

	JWTSigningCycle(t, client, ctx, rGen.Result.ID)
}

func Test_Integration_AESencryptingLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.SYAaes
	purpose := vault.KPencryption
	name := GetName("Test_Integration_AESencryptingLifeCycle")

	rGen, err := client.SymmetricGenerate(ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: name,
			},
			Algorithm: algorithm,
			Purpose:   purpose,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)

	EncryptionCycle(t, client, ctx, rGen.Result.ID)
}

func Test_Integration_AESstoreLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	purpose := vault.KPencryption
	name := GetName("Test_Integration_AESstoreLifeCycle")

	rStore, err := client.SymmetricStore(ctx,
		&vault.SymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{
				Name: name,
			},
			Algorithm: KEY_AES_algorithm,
			Purpose:   purpose,
			Key:       vault.EncodedSymmetricKey(KEY_AES_key),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rStore)
	assert.NotNil(t, rStore.Result)
	assert.NotEmpty(t, rStore.Result.ID)

	EncryptionCycle(t, client, ctx, rStore.Result.ID)
}

func Test_Integration_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	client := vault.New(cfg)
	name := GetName("Test_Integration_AESstoreLifeCycle")

	input := &vault.SecretStoreRequest{
		CommonStoreRequest: vault.CommonStoreRequest{
			Name: name,
		},
		Secret: "somesecret",
	}

	out, err := client.SecretStore(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}
