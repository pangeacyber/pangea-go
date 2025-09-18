package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea/rsa"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/vault"
	"github.com/spf13/cobra"
)

func init() {
	vaultCmd := &cobra.Command{
		Use:   "vault",
		Short: "Vault examples",
	}

	encryptCmd := &cobra.Command{
		Use:  "encrypt",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return encrypt(cmd)
		},
	}
	vaultCmd.AddCommand(encryptCmd)

	exportCmd := &cobra.Command{
		Use:  "export",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return export(cmd)
		},
	}
	vaultCmd.AddCommand(exportCmd)

	fpeCmd := &cobra.Command{
		Use:  "fpe",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fpe(cmd)
		},
	}
	vaultCmd.AddCommand(fpeCmd)

	getCmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return get(cmd)
		},
	}
	vaultCmd.AddCommand(getCmd)

	rotateCmd := &cobra.Command{
		Use:  "rotate",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return rotate(cmd)
		},
	}
	vaultCmd.AddCommand(rotateCmd)

	signCmd := &cobra.Command{
		Use:  "sign",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sign(cmd)
		},
	}
	vaultCmd.AddCommand(signCmd)

	structuredEncryptCmd := &cobra.Command{
		Use:  "structured_encrypt",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return structuredEncrypt(cmd)
		},
	}
	vaultCmd.AddCommand(structuredEncryptCmd)

	ExamplesCmd.AddCommand(vaultCmd)
}

func export(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	vaultcli := vault.New(config)

	ctx := context.Background()

	fmt.Println("Generate key with exportable field on true...")
	name := "Go encrypt example " + time.Now().Format(time.RFC3339)
	generateInput := &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes256_cbc,
		Purpose:   vault.KPencryption,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: name,
		},
		Exportable: pangea.Bool(true),
	}
	generateResponse, err := vaultcli.SymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(generateResponse.Result))

	id := generateResponse.Result.ID

	// Export with no encryption
	fmt.Println("Exporting key without encryption...")
	rExp, err := vaultcli.Export(ctx,
		&vault.ExportRequest{
			ID: id,
		})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(rExp.Result))

	// Export with encryption
	fmt.Println("Export key encrypted...")

	// Generate a RSA key pair
	rsaPubKey, rsaPrivKey, err := rsa.GenerateKeyPair(4096)
	if err != nil {
		log.Fatal(err)
	}

	// Should send public key in PEM format to encrypt exported key
	rsaPubKeyPEM, err := rsa.EncodePEMPublicKey(rsaPubKey)
	if err != nil {
		log.Fatal(err)
	}

	ea := vault.EEArsa4096_oaep_sha512
	rExpEnc, err := vaultcli.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
		})
	if err != nil {
		log.Fatal(err)
	}

	// Decode base64 key field
	expKeyDec, err := base64.StdEncoding.DecodeString(*rExpEnc.Result.Key)
	if err != nil {
		log.Fatal(err)
	}

	// Decrypt decoded field
	expKey, err := rsa.DecryptSHA512(rsaPrivKey, expKeyDec)
	if err != nil {
		log.Fatal(err)
	}

	// Use decrypted key
	fmt.Println("Decrypted key:")
	fmt.Println(string(expKey))
	return nil
}

func fpe(cmd *cobra.Command) error {
	// Set up a Pangea Vault client.
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Missing `PANGEA_VAULT_TOKEN` environment variable.")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	vaultClient := vault.New(config)

	ctx := context.Background()

	// Plain text that we'll encrypt.
	plainText := "123-4567-8901"

	// Optional tweak string.
	tweak := "MTIzMTIzMT=="

	// Generate an encryption key.
	generated, err := vaultClient.SymmetricGenerate(ctx, &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes_ff3_1_256,
		Purpose:   vault.KPfpe,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: "go-fpe-example-" + time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	keyId := generated.Result.ID

	// Encrypt the plain text.
	encrypted, err := vaultClient.EncryptTransform(ctx, &vault.EncryptTransformRequest{
		ID:        keyId,
		PlainText: plainText,
		Tweak:     &tweak,
		Alphabet:  vault.TAnumeric,
	})
	if err != nil {
		log.Fatal(err)
	}
	encryptedText := encrypted.Result.CipherText
	fmt.Printf("Plain text: %s. Encrypted text: %s.\n", plainText, encryptedText)

	// Decrypt the result to get back the text we started with.
	decrypted, err := vaultClient.DecryptTransform(ctx, &vault.DecryptTransformRequest{
		ID:         keyId,
		CipherText: encryptedText,
		Tweak:      tweak,
		Alphabet:   vault.TAnumeric,
	})
	if err != nil {
		log.Fatal(err)
	}
	decryptedText := decrypted.Result.PlainText
	fmt.Printf("Original text: %s. Decrypted text: %s.\n", plainText, decryptedText)
	return nil
}

func get(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Error: No token present")
	}

	auditTokenID := os.Getenv("PANGEA_AUDIT_TOKEN_VAULT_ID")
	if auditTokenID == "" {
		log.Fatal("Error: No audit token id present")
	}
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}
	vaultClient := vault.New(config)

	ctx := context.Background()

	fmt.Println("Fetch the audit token...")
	getRequest := &vault.GetRequest{
		ID: auditTokenID,
	}
	storeResponse, err := vaultClient.Get(ctx, getRequest)
	if err != nil {
		log.Fatal(err)
	}
	auditToken := storeResponse.Result.ItemVersions[0].Token
	if auditToken == nil {
		log.Fatal("Unexpected: token not present")
	}

	fmt.Println("Initialize Log...")
	auditConfig, err := pangea.NewConfig(
		option.WithToken(*auditToken),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}
	auditClient, err := audit.New(auditConfig)
	if err != nil {
		log.Fatal("failed to create audit client")
	}
	event := &audit.StandardEvent{
		Message: "Hello, World!",
	}
	lr, err := auditClient.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*audit.StandardEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
	return nil
}

func rotate(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	vaultcli := vault.New(config)

	const (
		secretV1 = "my first secret"
		secretV2 = "my second secret"
	)

	ctx := context.Background()

	fmt.Println("Store secret...")
	name := "Go rotate example " + time.Now().Format(time.RFC3339)
	storeInput := &vault.SecretStoreRequest{
		Secret: secretV1,
		CommonStoreRequest: vault.CommonStoreRequest{
			Name: name,
			Type: vault.ITsecret,
		},
	}
	storeResponse, err := vaultcli.SecretStore(ctx, storeInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(storeResponse.Result))

	fmt.Println("Rotate secret...")
	rotateInput := &vault.SecretRotateRequest{
		CommonRotateRequest: vault.CommonRotateRequest{
			ID: storeResponse.Result.ID,
		},
		Secret: secretV1,
	}

	rotateResponse, err := vaultcli.SecretRotate(ctx, rotateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(rotateResponse.Result))

	fmt.Println("Get last version")
	getInput := &vault.GetRequest{
		ID: storeResponse.Result.ID,
	}

	getResponse, err := vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(getResponse.Result))

	fmt.Println("Get version 1")
	getInput = &vault.GetRequest{
		ID:      storeResponse.Result.ID,
		Version: "1",
	}

	getResponse, err = vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(getResponse.Result))
	return nil
}

func sign(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	vaultcli := vault.New(config)

	ctx := context.Background()

	fmt.Println("Generate key...")
	name := "Go sign example " + time.Now().Format(time.RFC3339)
	generateInput := &vault.AsymmetricGenerateRequest{
		Algorithm: vault.AAed25519,
		Purpose:   vault.KPsigning,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: name,
		},
	}
	generateResponse, err := vaultcli.AsymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(generateResponse.Result))

	fmt.Println("sign...")
	message := "messagetosign"
	data := base64.StdEncoding.EncodeToString([]byte(message))

	signInput := &vault.SignRequest{
		ID:      generateResponse.Result.ID,
		Message: data,
	}

	signResponse, err := vaultcli.Sign(ctx, signInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(signResponse.Result))

	fmt.Println("Verify...")
	verifyInput := &vault.VerifyRequest{
		ID:        generateResponse.Result.ID,
		Message:   data,
		Signature: signResponse.Result.Signature,
	}

	verifyResponse, err := vaultcli.Verify(ctx, verifyInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(verifyResponse.Result))

	if verifyResponse.Result.ValidSignature {
		fmt.Println("Verify success")
	} else {
		fmt.Println("Verify failed")
	}
	return nil
}

func structuredEncrypt(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("missing PANGEA_VAULT_TOKEN environment variable")
	}

	domain := os.Getenv("PANGEA_DOMAIN")
	if domain == "" {
		log.Fatal("missing PANGEA_DOMAIN environment variable")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(domain),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	vaultClient := vault.New(config)

	ctx := context.Background()

	// First create an encryption key, either from the Pangea Console or
	// programmatically as below.
	generateInput := &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes256_cfb,
		Purpose:   vault.KPencryption,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: "Go structured encrypt example " + time.Now().Format(time.RFC3339),
		},
	}
	generateResponse, err := vaultClient.SymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	encryptionKeyId := generateResponse.Result.ID

	// Structured data that we'll encrypt.
	data := map[string]interface{}{
		"foo":  [4]interface{}{1, 2, "true", "false"},
		"some": "thing",
	}

	encryptInput := &vault.EncryptStructuredRequest{
		ID:             encryptionKeyId,
		StructuredData: data,
		Filter:         "$.foo[2:4]",
	}
	encryptResponse, err := vaultClient.EncryptStructured(ctx, encryptInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Encrypted result:", pangea.Stringify(encryptResponse.Result.StructuredData))
	return nil
}

func encrypt(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	vaultcli := vault.New(config)

	ctx := context.Background()

	fmt.Println("Generate key...")
	name := "Go encrypt example " + time.Now().Format(time.RFC3339)
	generateInput := &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes128_cfb,
		Purpose:   vault.KPencryption,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: name,
		},
	}
	generateResponse, err := vaultcli.SymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(generateResponse.Result))

	fmt.Println("Encrypt...")
	message := "messagetoencrypt"
	data := base64.StdEncoding.EncodeToString([]byte(message))

	encryptInput := &vault.EncryptRequest{
		ID:        generateResponse.Result.ID,
		PlainText: data,
	}

	encryptResponse, err := vaultcli.Encrypt(ctx, encryptInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(encryptResponse.Result))

	fmt.Println("Decrypt...")
	decryptInput := &vault.DecryptRequest{
		ID:         generateResponse.Result.ID,
		CipherText: encryptResponse.Result.CipherText,
	}

	decryptResponse, err := vaultcli.Decrypt(ctx, decryptInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(decryptResponse.Result))

	if decryptResponse.Result.PlainText == data {
		fmt.Println("Encrypt/Decrypt success")
	} else {
		fmt.Println("Encrypt/Decrypt failed")
	}
	return nil
}
