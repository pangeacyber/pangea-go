//go:build integration

package vault_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea/rsa"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/vault"
	"github.com/stretchr/testify/assert"
)

const (
	actor = "GoSDKTest"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("vault", pangeatesting.Live)

var timeNow = time.Now()
var timeStr = timeNow.Format("20060102_150405")
var KEY_ED25519_algorithm = vault.AAed25519
var KEY_ED25519_private_key = "-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIGthqegkjgddRAn0PWN2FeYC6HcCVQf/Ph9sUbeprTBO\n-----END PRIVATE KEY-----\n"
var KEY_ED25519_public_key = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAPlGrDliJXUbPc2YWEhFxlL2UbBfLHc3ed1f36FrDtTc=\n-----END PUBLIC KEY-----\n"
var KEY_AES_algorithm = vault.SYAaes
var KEY_AES_key = "oILlp2FUPHWiaqFXl4/1ww=="

func PrintPangeaAPIError(err error) {
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
			Type: vault.ITsecret,
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
	assert.NotEmpty(t, rStore.Result.ID)
	assert.Equal(t, 1, rStore.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rStore.Result.NumVersions)
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
	assert.Equal(t, 1, len(rRotate.Result.ItemVersions))
	assert.Nil(t, rRotate.Result.ItemVersions[0].Secret)
	assert.Equal(t, 2, rRotate.Result.ItemVersions[0].Version)
	assert.Equal(t, 2, rRotate.Result.NumVersions)
	assert.Equal(t, string(vault.ITsecret), rRotate.Result.Type)

	rUpdate, err := client.Update(ctx,
		&vault.UpdateRequest{
			ID:     ID,
			Folder: "updated",
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, rUpdate)
	assert.NotNil(t, rUpdate.Result)
	assert.Equal(t, ID, rUpdate.Result.ID)

	rGet, err := client.Get(ctx,
		&vault.GetRequest{
			ID: ID,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 1, len(rGet.Result.ItemVersions))
	assert.Equal(t, secretV2, *rGet.Result.ItemVersions[0].Secret)
	assert.Equal(t, string(vault.IVSactive), rGet.Result.ItemVersions[0].State)
	assert.Nil(t, rGet.Result.ItemVersions[0].PublicKey)
	assert.Nil(t, rGet.Result.ItemVersions[0].DestroyedAt)
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
	assert.Equal(t, 1, len(rGet.Result.ItemVersions))
	assert.Equal(t, secretV2, *rGet.Result.ItemVersions[0].Secret)
	assert.Nil(t, rGet.Result.ItemVersions[0].PublicKey)
	assert.Equal(t, string(vault.IVSsuspended), rGet.Result.ItemVersions[0].State)
	assert.Nil(t, rGet.Result.ItemVersions[0].DestroyedAt)
	assert.Equal(t, string(vault.ITsecret), rGet.Result.Type)
}

func AsymSigningCycle(t *testing.T, client vault.Client, ctx context.Context, id string) {
	message := "thisisamessagetosign"
	data := base64.StdEncoding.EncodeToString([]byte(message))

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
	assert.Nil(t, rRotate.Result.ItemVersions[0].PublicKey)
	assert.Equal(t, 2, rRotate.Result.ItemVersions[0].Version)
	assert.Equal(t, 2, rRotate.Result.NumVersions)
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

func JWTSymSigningCycle(t *testing.T, client vault.Client, ctx context.Context, id string) {
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
	assert.Equal(t, 2, rRotate.Result.ItemVersions[0].Version)
	assert.Equal(t, 2, rRotate.Result.NumVersions)
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

func JWTAsymSigningCycle(t *testing.T, client vault.Client, ctx context.Context, id string) {
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
	assert.Equal(t, 2, rRotate.Result.ItemVersions[0].Version)
	assert.Equal(t, 2, rRotate.Result.NumVersions)
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
	PrintPangeaAPIError(err)
	assert.NoError(t, err)
	assert.NotNil(t, rGet)
	assert.NotNil(t, rGet.Result)
	assert.Equal(t, 1, len(rGet.Result.Keys))

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
	assert.Equal(t, 1, len(rGet.Result.Keys))

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
	assert.Equal(t, 2, len(rGet.Result.Keys))

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
	assert.Equal(t, 1, len(rGet.Result.Keys))

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
	assert.Equal(t, 2, rRot1.Result.ItemVersions[0].Version)
	assert.Equal(t, 2, rRot1.Result.NumVersions)

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

func Test_Integration_AsymmetricSigningGenerate(t *testing.T) {
	algorithms := []vault.AsymmetricAlgorithm{
		vault.AAed25519,
		vault.AArsa2048_pss_sha256,
		vault.AArsa3072_pss_sha256,
		vault.AArsa4096_pss_sha256,
		vault.AArsa4096_pss_sha512,
		vault.AArsa2048_pkcs1v15_sha256,
		vault.AAes256K,
		vault.AAed25519_dilithium2_beta,
		vault.AAed488_dilithium3_beta,
		vault.AAsphincsplus_128f_shake256_simple_beta,
		vault.AAsphincsplus_128f_shake256_robust_beta,
		vault.AAsphincsplus_192f_shake256_simple_beta,
		vault.AAsphincsplus_192f_shake256_robust_beta,
		vault.AAsphincsplus_256f_shake256_simple_beta,
		vault.AAsphincsplus_256f_shake256_robust_beta,
		vault.AAsphincsplus_128f_sha256_simple_beta,
		vault.AAsphincsplus_128f_sha256_robust_beta,
		vault.AAsphincsplus_192f_sha256_simple_beta,
		vault.AAsphincsplus_192f_sha256_robust_beta,
		vault.AAsphincsplus_256f_sha256_simple_beta,
		vault.AAsphincsplus_256f_sha256_robust_beta,
		vault.AAfalcon1024_beta,
	}

	failed := false
	purpose := vault.KPsigning
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	for _, algorithm := range algorithms {
		name := GetName(fmt.Sprintf("Test_Integration_%s_%s_generate", purpose, algorithm))
		// Generate
		_, err := client.AsymmetricGenerate(ctx,
			&vault.AsymmetricGenerateRequest{
				CommonGenerateRequest: vault.CommonGenerateRequest{
					Name: name,
				},
				Algorithm: algorithm,
				Purpose:   purpose,
			})

		if err != nil {
			fmt.Printf("Failed to generate %s %s\n%s...\n\n", algorithm, purpose, err)
			failed = true
		}
	}
	assert.False(t, failed)
}

func Test_Integration_AsymmetricEncryptionGenerate(t *testing.T) {
	algorithms := []vault.AsymmetricAlgorithm{
		vault.AArsa2048_oaep_sha256,
		vault.AArsa2048_oaep_sha1,
		vault.AArsa2048_oaep_sha512,
		vault.AArsa3072_oaep_sha1,
		vault.AArsa3072_oaep_sha256,
		vault.AArsa3072_oaep_sha512,
		vault.AArsa4096_oaep_sha1,
		vault.AArsa4096_oaep_sha256,
		vault.AArsa4096_oaep_sha512,
	}

	failed := false
	purpose := vault.KPencryption
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	for _, algorithm := range algorithms {
		name := GetName(fmt.Sprintf("Test_Integration_%s_%s_generate", purpose, algorithm))
		// Generate
		_, err := client.AsymmetricGenerate(ctx,
			&vault.AsymmetricGenerateRequest{
				CommonGenerateRequest: vault.CommonGenerateRequest{
					Name: name,
				},
				Algorithm: algorithm,
				Purpose:   purpose,
			})

		if err != nil {
			fmt.Printf("Failed to generate %s %s...\n%s\n\n", algorithm, purpose, err)
			failed = true
		}
	}
	assert.False(t, failed)
}

func Test_Integration_SymmetricEncryptionGenerate(t *testing.T) {
	algorithms := []vault.SymmetricAlgorithm{
		vault.SYAaes128_cfb,
		vault.SYAaes256_cfb,
		vault.SYAaes256_gcm,
		vault.SYAaes128_cbc,
		vault.SYAaes256_cbc,
	}

	failed := false
	purpose := vault.KPencryption
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	for _, algorithm := range algorithms {
		name := GetName(fmt.Sprintf("Test_Integration_%s_%s_generate", purpose, algorithm))
		// Generate
		_, err := client.SymmetricGenerate(ctx,
			&vault.SymmetricGenerateRequest{
				CommonGenerateRequest: vault.CommonGenerateRequest{
					Name: name,
				},
				Algorithm: algorithm,
				Purpose:   purpose,
			})

		if err != nil {
			fmt.Printf("Failed to generate %s %s...\n%s\n\n", algorithm, purpose, err)
			failed = true
		}
	}
	assert.False(t, failed)
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
	assert.Equal(t, 1, len(rGen.Result.ItemVersions))
	assert.Nil(t, rGen.Result.ItemVersions[0].PublicKey)
	assert.NotEmpty(t, rGen.Result.ID)
	assert.Equal(t, 1, rGen.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rGen.Result.NumVersions)

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
	assert.Nil(t, rStore.Result.ItemVersions[0].PublicKey)
	assert.NotEmpty(t, rStore.Result.ID)
	assert.Equal(t, 1, rStore.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rStore.Result.NumVersions)

	AsymSigningCycle(t, client, ctx, rStore.Result.ID)
}

func Test_Integration_JWT_AsymSigningLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithms := []vault.AsymmetricAlgorithm{vault.AAes256, vault.AAes384, vault.AAes512}
	purpose := vault.KPjwt

	for _, algorithm := range algorithms {
		fmt.Printf("\nRunning Test_Integration_JWT_AsymSigningLifeCycle with %s\n", algorithm)
		name := GetName("Test_Integration_JWT_AsymSigningLifeCycle")
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
		assert.Equal(t, 1, len(rGen.Result.ItemVersions))
		assert.Nil(t, rGen.Result.ItemVersions[0].PublicKey)
		assert.NotEmpty(t, rGen.Result.ID)
		assert.Equal(t, 1, rGen.Result.ItemVersions[0].Version)
		assert.Equal(t, 1, rGen.Result.NumVersions)

		JWTAsymSigningCycle(t, client, ctx, rGen.Result.ID)
		fmt.Printf("Finished Test_Integration_JWT_AsymSigningLifeCycle with %s\n", algorithm)
	}
}

func Test_Integration_JWT_SymSigningLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithms := []vault.SymmetricAlgorithm{vault.SYAhs256, vault.SYAhs512, vault.SYAhs384}
	purpose := vault.KPjwt

	for _, algorithm := range algorithms {
		fmt.Printf("\nRunning Test_Integration_JWT_SymSigningLifeCycle with %s\n", algorithm)
		name := GetName("Test_Integration_JWT_SymSigningLifeCycle")
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

		JWTSymSigningCycle(t, client, ctx, rGen.Result.ID)
		fmt.Printf("Finished Test_Integration_JWT_SymSigningLifeCycle with %s\n", algorithm)
	}
}

func Test_Integration_AESencryptingLifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithms := []vault.SymmetricAlgorithm{
		vault.SYAaes128_cfb,
		vault.SYAaes128_cfb,
		vault.SYAaes256_cbc,
		vault.SYAaes256_cfb,
		vault.SYAaes256_gcm,
	}
	purpose := vault.KPencryption

	for _, algorithm := range algorithms {
		fmt.Printf("\nRunning Test_Integration_AESencryptingLifeCycle with %s\n", algorithm)
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
		fmt.Printf("Finished Test_Integration_AESencryptingLifeCycle with %s\n", algorithm)
	}
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

func Test_List_And_Delete(t *testing.T) {
	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	client := vault.New(cfg)

	filter := vault.NewFilterList()
	filter.Name().SetContains([]string{actor})
	last := ""
	item_counter := 0
	list_call_counter := 0
	start := time.Now().UnixMilli()

	for item_counter < 500 {
		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()

		lreq := &vault.ListRequest{
			Filter: filter.Filter(),
			Last:   last,
		}

		list_call_counter++
		fmt.Printf("List call %d\n", list_call_counter)
		lresp, err := client.List(ctx, lreq)
		last = lresp.Result.Last

		assert.NoError(t, err)
		assert.NotNil(t, lresp)

		assert.GreaterOrEqual(t, len(lresp.Result.Items), 0)
		for _, i := range lresp.Result.Items {
			if i.ID != "" && i.Type != "folder" && i.Folder != "/service-tokens/" {
				item_counter++
				ctx_del, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancelFn()
				dresp, err := client.Delete(ctx_del, &vault.DeleteRequest{
					ID: i.ID,
				})

				assert.NoError(t, err)
				assert.NotNil(t, dresp)
				assert.NotNil(t, dresp.Result)
			}
		}

		if len(lresp.Result.Items) == 0 {
			break
		}

	}

	end := time.Now().UnixMilli()
	fmt.Printf("Deleted %d keys\n", item_counter)
	fmt.Printf("Deleted %d keys in %d s\n", item_counter, (end-start)/1000)
	fmt.Printf("Average delete time: %f ms\n", float64(end-start)/float64(item_counter))
}

func Test_List_And_Delete_Folders(t *testing.T) {
	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	client := vault.New(cfg)
	list_call_counter := 0

	for _, name := range []string{"test_folder_name_new", "test_folder_name", "test_parent_folder_"} {
		fmt.Printf("Deleting folders with name '%s'\n", name)
		last := ""
		item_counter := 0
		start := time.Now().UnixMilli()
		filter := vault.NewFilterList()
		filter.Name().SetContains([]string{name})

		for item_counter < 100 {
			ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancelFn()
			fmt.Printf("List call %d. Deleted items: %d\n", list_call_counter, item_counter)

			lreq := &vault.ListRequest{
				Filter: filter.Filter(),
				Last:   last,
			}
			lresp, err := client.List(ctx, lreq)
			assert.NoError(t, err)
			assert.NotNil(t, lresp)
			list_call_counter++
			last = lresp.Result.Last

			assert.GreaterOrEqual(t, len(lresp.Result.Items), 0)
			for _, i := range lresp.Result.Items {
				if i.ID != "" && i.Type == "folder" && i.Folder != "/service-tokens/" {
					item_counter++
					ctx_del, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
					defer cancelFn()
					dresp, err := client.Delete(ctx_del, &vault.DeleteRequest{
						ID: i.ID,
					})

					assert.NoError(t, err)
					assert.NotNil(t, dresp)
					assert.NotNil(t, dresp.Result)
				}
			}

			if len(lresp.Result.Items) == 0 {
				break
			}
		}

		end := time.Now().UnixMilli()
		fmt.Printf("Deleted folders in %d s\n", (end-start)/1000)
		fmt.Printf("Average delete time: %f ms\n", float64(end-start)/float64(item_counter))
	}
}

func Test_Integration_Folders(t *testing.T) {
	var FOLDER_PARENT = "test_parent_folder_" + timeStr
	var FOLDER_NAME = "test_folder_name"
	var FOLDER_NAME_NEW = "test_folder_name_new"

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	client := vault.New(cfg)

	// Create parent
	pcr, err := client.FolderCreate(
		ctx,
		&vault.FolderCreateRequest{
			Name:   FOLDER_PARENT,
			Folder: "/",
		},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, pcr.Result.ID)

	// Create folder
	fcr, err := client.FolderCreate(
		ctx,
		&vault.FolderCreateRequest{
			Name:   FOLDER_NAME,
			Folder: FOLDER_PARENT,
		},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, fcr.Result.ID)

	// Update
	ur, err := client.Update(
		ctx,
		&vault.UpdateRequest{
			ID:   fcr.Result.ID,
			Name: FOLDER_NAME_NEW,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, fcr.Result.ID, ur.Result.ID)

	filter := vault.NewFilterList()
	filter.Folder().Set(pangea.String(FOLDER_PARENT))

	// List
	lr, err := client.List(
		ctx,
		&vault.ListRequest{
			Filter: filter.Filter(),
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(lr.Result.Items))
	assert.Equal(t, FOLDER_NAME_NEW, lr.Result.Items[0].Name)
	assert.Equal(t, fcr.Result.ID, lr.Result.Items[0].ID)

	// Delete folder
	dfr, err := client.Delete(
		ctx,
		&vault.DeleteRequest{
			ID: fcr.Result.ID,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, dfr.Result.ID, fcr.Result.ID)

	// Delete parent
	dpr, err := client.Delete(
		ctx,
		&vault.DeleteRequest{
			ID: pcr.Result.ID,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, pcr.Result.ID, dpr.Result.ID)
}

func Test_Integration_EncryptStructured(t *testing.T) {
	// Test data.
	data := map[string]interface{}{
		"field1": [4]interface{}{1, 2, "true", "false"},
		"field2": "true",
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	// Generate an encryption key.
	rGen, err := client.SymmetricGenerate(
		ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: GetName("Test_Integration_EncryptStructured"),
			},
			Algorithm: vault.SYAaes256_cfb,
			Purpose:   vault.KPencryption,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)
	key := rGen.Result.ID

	// Encrypt.
	encryptedResponse, err := client.EncryptStructured(
		ctx,
		&vault.EncryptStructuredRequest{
			ID:             key,
			StructuredData: data,
			Filter:         "$.field1[2:4]",
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, encryptedResponse)
	assert.NotNil(t, encryptedResponse.Result)
	assert.Equal(t, key, encryptedResponse.Result.ID)
	assert.Len(t, encryptedResponse.Result.StructuredData["field1"], 4)
	assert.Equal(t, data["field2"], encryptedResponse.Result.StructuredData["field2"])

	// Decrypt what we encrypted.
	decryptedResponse, err := client.DecryptStructured(
		ctx,
		&vault.EncryptStructuredRequest{
			ID:             key,
			StructuredData: encryptedResponse.Result.StructuredData,
			Filter:         "$.field1[2:4]",
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedResponse)
	assert.NotNil(t, decryptedResponse.Result)
	assert.Equal(t, key, decryptedResponse.Result.ID)
	assert.Len(t, decryptedResponse.Result.StructuredData["field1"], 4)
	assert.Equal(t, data["field2"], decryptedResponse.Result.StructuredData["field2"])
}

func Test_Integration_EncryptTransform(t *testing.T) {
	// Test data.
	plainText := "123-4567-8901"
	tweak := "MTIzMTIzMT=="

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	// Generate an encryption key.
	rGen, err := client.SymmetricGenerate(
		ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: GetName("Test_Integration_EncryptTransform"),
			},
			Algorithm: vault.SYAaes_ff3_1_256,
			Purpose:   vault.KPfpe,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)
	key := rGen.Result.ID

	// Encrypt.
	encryptedResponse, err := client.EncryptTransform(
		ctx,
		&vault.EncryptTransformRequest{
			ID:        key,
			PlainText: plainText,
			Tweak:     &tweak,
			Alphabet:  vault.TAalphanumeric,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, encryptedResponse)
	assert.NotNil(t, encryptedResponse.Result)
	assert.Equal(t, key, encryptedResponse.Result.ID)
	assert.Len(t, encryptedResponse.Result.CipherText, len(plainText))

	// Decrypt what we encrypted.
	decryptedResponse, err := client.DecryptTransform(
		ctx,
		&vault.DecryptTransformRequest{
			ID:         key,
			CipherText: encryptedResponse.Result.CipherText,
			Tweak:      tweak,
			Alphabet:   vault.TAalphanumeric,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedResponse)
	assert.NotNil(t, decryptedResponse.Result)
	assert.Equal(t, key, decryptedResponse.Result.ID)
	assert.Equal(t, plainText, decryptedResponse.Result.PlainText)
}

func Test_Integration_ExportGenerateAsymmetric(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.AAed25519
	purpose := vault.KPsigning
	name := GetName("Test_Integration_ExportGenerateAsymmetric")

	// Generate
	rGen, err := client.AsymmetricGenerate(ctx,
		&vault.AsymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: name,
			},
			Algorithm:  algorithm,
			Purpose:    purpose,
			Exportable: pangea.Bool(true),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.Equal(t, 1, len(rGen.Result.ItemVersions))
	assert.Nil(t, rGen.Result.ItemVersions[0].PublicKey)
	assert.NotEmpty(t, rGen.Result.ID)
	id := rGen.Result.ID
	assert.Equal(t, 1, rGen.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rGen.Result.NumVersions)

	// Export with no encryption
	rExp, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:      id,
			Version: pangea.Int(1),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rExp.Result.PrivateKey)
	assert.NotEmpty(t, *rExp.Result.PrivateKey)
	assert.NotNil(t, rExp.Result.PublicKey)
	assert.NotEmpty(t, *rExp.Result.PublicKey)
	assert.Nil(t, rExp.Result.Key)
	assert.Equal(t, rExp.Result.EncryptionType, "none")

	// Export with asymmetric encryption
	rsaPubKey, rsaPrivKey, err := rsa.GenerateKeyPair(4096)
	assert.NoError(t, err)

	rsaPubKeyPEM, err := rsa.EncodePEMPublicKey(rsaPubKey)
	assert.NoError(t, err)

	ea := vault.EEArsa4096_oaep_sha512

	rExpEnc, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rExpEnc.Result.PrivateKey)
	assert.NotEmpty(t, *rExpEnc.Result.PrivateKey)
	assert.NotNil(t, rExpEnc.Result.PublicKey)
	assert.NotEmpty(t, *rExpEnc.Result.PublicKey)
	assert.Nil(t, rExpEnc.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "asymmetric")

	expPrivKeyDec, err := base64.StdEncoding.DecodeString(*rExpEnc.Result.PrivateKey)
	assert.NoError(t, err)
	expPrivKey, err := rsa.DecryptSHA512(rsaPrivKey, expPrivKeyDec)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.PrivateKey, string(expPrivKey))
	assert.Equal(t, *rExp.Result.PublicKey, *rExpEnc.Result.PublicKey)

	// Export encrypted kem
	ea = vault.EEArsa4096_no_padding_kem
	password := "mypassword"

	rExpEnc, err = client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
			KEMPassword:         &password,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rExpEnc.Result.PrivateKey)
	assert.NotEmpty(t, *rExpEnc.Result.PrivateKey)
	assert.NotNil(t, rExpEnc.Result.PublicKey)
	assert.NotEmpty(t, *rExpEnc.Result.PublicKey)
	assert.Nil(t, rExpEnc.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "kem")

	expPrivKeyDec, err = base64.StdEncoding.DecodeString(*rExpEnc.Result.PrivateKey)
	assert.NoError(t, err)

	input, err := vault.NewKEMDecryptInput(*rExpEnc.Result, password, *rsaPrivKey)
	assert.NoError(t, err)
	assert.NotNil(t, input)

	expPrivKey, err = vault.KEMDecrypt(*input)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.PrivateKey, string(expPrivKey))
	assert.Equal(t, *rExp.Result.PublicKey, *rExpEnc.Result.PublicKey)
}

func Test_Integration_ExportStoreAsymmetric(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	name := GetName("Test_Integration_ExportStoreAsymmetric")

	// Generate
	rGen, err := client.AsymmetricStore(ctx,
		&vault.AsymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{
				Name: name,
			},
			Exportable: pangea.Bool(true),
			Algorithm:  KEY_ED25519_algorithm,
			Purpose:    vault.KPsigning,
			PublicKey:  vault.EncodedPublicKey(KEY_ED25519_public_key),
			PrivateKey: vault.EncodedPrivateKey(KEY_ED25519_private_key),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.Equal(t, 1, len(rGen.Result.ItemVersions))
	assert.Nil(t, rGen.Result.ItemVersions[0].PublicKey)
	assert.NotEmpty(t, rGen.Result.ID)
	id := rGen.Result.ID
	assert.Equal(t, 1, rGen.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rGen.Result.NumVersions)

	// Export with no encryption
	rExp, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:      id,
			Version: pangea.Int(1),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rExp.Result.PrivateKey)
	assert.NotEmpty(t, *rExp.Result.PrivateKey)
	assert.NotNil(t, rExp.Result.PublicKey)
	assert.NotEmpty(t, *rExp.Result.PublicKey)
	assert.Nil(t, rExp.Result.Key)
	assert.Equal(t, rExp.Result.EncryptionType, "none")

	// Export with encryption
	rsaPubKey, rsaPrivKey, err := rsa.GenerateKeyPair(4096)
	assert.NoError(t, err)
	assert.NotNil(t, rsaPrivKey)
	assert.NotNil(t, rsaPubKey)

	rsaPubKeyPEM, err := rsa.EncodePEMPublicKey(rsaPubKey)
	assert.NoError(t, err)

	// Export encrypted asymmetric
	ea := vault.EEArsa4096_oaep_sha512

	rExpEnc, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rExpEnc.Result.PrivateKey)
	assert.NotEmpty(t, *rExpEnc.Result.PrivateKey)
	assert.NotNil(t, rExpEnc.Result.PublicKey)
	assert.NotEmpty(t, *rExpEnc.Result.PublicKey)
	assert.Nil(t, rExpEnc.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "asymmetric")

	expPrivKeyDec, err := base64.StdEncoding.DecodeString(*rExpEnc.Result.PrivateKey)
	assert.NoError(t, err)
	expPrivKey, err := rsa.DecryptSHA512(rsaPrivKey, expPrivKeyDec)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.PrivateKey, string(expPrivKey))
	assert.Equal(t, *rExp.Result.PublicKey, *rExpEnc.Result.PublicKey)

	// Export encrypted kem
	ea = vault.EEArsa4096_no_padding_kem
	password := "mypassword"

	rExpEnc, err = client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
			KEMPassword:         &password,
		})

	assert.NoError(t, err)
	assert.NotNil(t, rExpEnc.Result.PrivateKey)
	assert.NotEmpty(t, *rExpEnc.Result.PrivateKey)
	assert.NotNil(t, rExpEnc.Result.PublicKey)
	assert.NotEmpty(t, *rExpEnc.Result.PublicKey)
	assert.Nil(t, rExpEnc.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "kem")

	expPrivKeyDec, err = base64.StdEncoding.DecodeString(*rExpEnc.Result.PrivateKey)
	assert.NoError(t, err)

	input, err := vault.NewKEMDecryptInput(*rExpEnc.Result, password, *rsaPrivKey)
	assert.NoError(t, err)
	assert.NotNil(t, input)

	expPrivKey, err = vault.KEMDecrypt(*input)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.PrivateKey, string(expPrivKey))
	assert.Equal(t, *rExp.Result.PublicKey, *rExpEnc.Result.PublicKey)
}

func Test_Integration_ExportGenerateSymmetric(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	algorithm := vault.SYAaes128_cbc
	purpose := vault.KPencryption
	name := GetName("Test_Integration_ExportGenerateSymmetric")

	// Generate
	rGen, err := client.SymmetricGenerate(ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: name,
			},
			Algorithm:  algorithm,
			Purpose:    purpose,
			Exportable: pangea.Bool(true),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)
	id := rGen.Result.ID
	assert.Equal(t, 1, len(rGen.Result.ItemVersions))
	assert.Equal(t, 1, rGen.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rGen.Result.NumVersions)

	// Export with no encryption
	rExp, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:      id,
			Version: pangea.Int(1),
		})

	assert.NoError(t, err)
	assert.Nil(t, rExp.Result.PrivateKey)
	assert.Nil(t, rExp.Result.PublicKey)
	assert.NotNil(t, rExp.Result.Key)
	assert.NotEmpty(t, *rExp.Result.Key)
	assert.Equal(t, rExp.Result.EncryptionType, "none")

	// Export with asymmetric encryption
	rsaPubKey, rsaPrivKey, err := rsa.GenerateKeyPair(4096)
	assert.NoError(t, err)

	rsaPubKeyPEM, err := rsa.EncodePEMPublicKey(rsaPubKey)
	assert.NoError(t, err)

	ea := vault.EEArsa4096_oaep_sha512

	rExpEnc, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
		})

	assert.NoError(t, err)
	assert.Nil(t, rExp.Result.PrivateKey)
	assert.Nil(t, rExp.Result.PublicKey)
	assert.NotNil(t, rExp.Result.Key)
	assert.NotEmpty(t, *rExp.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "asymmetric")

	expKeyDec, err := base64.StdEncoding.DecodeString(*rExpEnc.Result.Key)
	assert.NoError(t, err)
	expKey, err := rsa.DecryptSHA512(rsaPrivKey, expKeyDec)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.Key, string(expKey))

	// Export encrypted kem
	ea = vault.EEArsa4096_no_padding_kem
	password := "mypassword"

	rExpEnc, err = client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
			KEMPassword:         &password,
		})

	assert.NoError(t, err)
	assert.Nil(t, rExpEnc.Result.PrivateKey)
	assert.Nil(t, rExpEnc.Result.PublicKey)
	assert.NotNil(t, rExpEnc.Result.Key)
	assert.NotEmpty(t, *rExpEnc.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "kem")

	expKeyDec, err = base64.StdEncoding.DecodeString(*rExpEnc.Result.Key)
	assert.NoError(t, err)

	input, err := vault.NewKEMDecryptInput(*rExpEnc.Result, password, *rsaPrivKey)
	assert.NoError(t, err)
	assert.NotNil(t, input)

	expKey, err = vault.KEMDecrypt(*input)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.Key, string(expKey))
}

func Test_Integration_ExportStoreSymmetric(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	name := GetName("Test_Integration_ExportStoreSymmetric")

	// Generate
	rGen, err := client.SymmetricStore(ctx,
		&vault.SymmetricStoreRequest{
			CommonStoreRequest: vault.CommonStoreRequest{
				Name: name,
			},
			Exportable: pangea.Bool(true),
			Algorithm:  KEY_AES_algorithm,
			Purpose:    vault.KPencryption,
			Key:        vault.EncodedSymmetricKey(KEY_AES_key),
		})

	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)
	id := rGen.Result.ID
	assert.Equal(t, 1, len(rGen.Result.ItemVersions))
	assert.Equal(t, 1, rGen.Result.ItemVersions[0].Version)
	assert.Equal(t, 1, rGen.Result.NumVersions)

	// Export with no encryption
	rExp, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:      id,
			Version: pangea.Int(1),
		})

	assert.NoError(t, err)
	assert.Nil(t, rExp.Result.PrivateKey)
	assert.Nil(t, rExp.Result.PublicKey)
	assert.NotNil(t, rExp.Result.Key)
	assert.NotEmpty(t, *rExp.Result.Key)
	assert.Equal(t, rExp.Result.EncryptionType, "none")

	// Export with asymmetric encryption
	rsaPubKey, rsaPrivKey, err := rsa.GenerateKeyPair(4096)
	assert.NoError(t, err)

	rsaPubKeyPEM, err := rsa.EncodePEMPublicKey(rsaPubKey)
	assert.NoError(t, err)

	ea := vault.EEArsa4096_oaep_sha512

	rExpEnc, err := client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
		})

	assert.NoError(t, err)
	assert.Nil(t, rExp.Result.PrivateKey)
	assert.Nil(t, rExp.Result.PublicKey)
	assert.NotNil(t, rExp.Result.Key)
	assert.NotEmpty(t, *rExp.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "asymmetric")

	expKeyDec, err := base64.StdEncoding.DecodeString(*rExpEnc.Result.Key)
	assert.NoError(t, err)
	expKey, err := rsa.DecryptSHA512(rsaPrivKey, expKeyDec)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.Key, string(expKey))

	// Export encrypted kem
	ea = vault.EEArsa4096_no_padding_kem
	password := "mypassword"

	rExpEnc, err = client.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
			KEMPassword:         &password,
		})

	assert.NoError(t, err)
	assert.Nil(t, rExpEnc.Result.PrivateKey)
	assert.Nil(t, rExpEnc.Result.PublicKey)
	assert.NotNil(t, rExpEnc.Result.Key)
	assert.NotEmpty(t, *rExpEnc.Result.Key)
	assert.Equal(t, rExpEnc.Result.EncryptionType, "kem")

	expKeyDec, err = base64.StdEncoding.DecodeString(*rExpEnc.Result.Key)
	assert.NoError(t, err)

	input, err := vault.NewKEMDecryptInput(*rExpEnc.Result, password, *rsaPrivKey)
	assert.NoError(t, err)
	assert.NotNil(t, input)

	expKey, err = vault.KEMDecrypt(*input)
	assert.NoError(t, err)
	assert.Equal(t, *rExp.Result.Key, string(expKey))
}

func Test_Integration_GetBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	filter := vault.NewFilterList()
	filter.Folder().Set(pangea.String("/"))
	size := 5

	// List
	gbr, err := client.GetBulk(
		ctx,
		&vault.GetBulkRequest{
			Filter: filter.Filter(),
			Size:   size,
		},
	)

	assert.NoError(t, err)
	// assert.Equal(t, gbr.Result.Count, size)  // FIXME: uncomment when fixed on backend
	assert.Equal(t, size, len(gbr.Result.Items))

}

func Test_Integration_EncryptTransformStructured(t *testing.T) {
	// Test data.
	plainText := "123-4567-8901"
	tweak := "MTIzMTIzMT=="

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	client := vault.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	// Generate an encryption key.
	rGen, err := client.SymmetricGenerate(
		ctx,
		&vault.SymmetricGenerateRequest{
			CommonGenerateRequest: vault.CommonGenerateRequest{
				Name: GetName("Test_Integration_EncryptStructured"),
			},
			Algorithm: vault.SYAaes_ff3_1_256,
			Purpose:   vault.KPfpe,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, rGen)
	assert.NotNil(t, rGen.Result)
	assert.NotEmpty(t, rGen.Result.ID)
	key := rGen.Result.ID

	// Encrypt.
	data := map[string]interface{}{
		"field1": [4]interface{}{plainText, 2, plainText, "false"},
		"field2": "true",
	}

	encryptedResponse, err := client.EncryptTransformStructured(
		ctx,
		&vault.EncryptTransformStructuredRequest{
			ID:             key,
			StructuredData: data,
			Filter:         "$.field1[2:4]",
			Tweak:          &tweak,
			Alphabet:       vault.TAalphanumeric,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, encryptedResponse)
	assert.NotNil(t, encryptedResponse.Result)
	assert.Equal(t, key, encryptedResponse.Result.ID)
	assert.Len(t, encryptedResponse.Result.StructuredData["field1"], 4)
	assert.Equal(t, data["field2"], encryptedResponse.Result.StructuredData["field2"])

	// Decrypt what we encrypted.
	decryptedResponse, err := client.DecryptTransformStructured(
		ctx,
		&vault.EncryptTransformStructuredRequest{
			ID:             key,
			StructuredData: encryptedResponse.Result.StructuredData,
			Filter:         "$.field1[2:4]",
			Tweak:          &tweak,
			Alphabet:       vault.TAalphanumeric,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedResponse)
	assert.NotNil(t, decryptedResponse.Result)
	assert.Equal(t, key, decryptedResponse.Result.ID)
	assert.Len(t, decryptedResponse.Result.StructuredData["field1"], 4)
	assert.Equal(t, data["field2"], decryptedResponse.Result.StructuredData["field2"])
}
