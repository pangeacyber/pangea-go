module examples/vault

go 1.24.1

require github.com/pangeacyber/pangea-go/pangea-sdk/v4 v4.4.0

require (
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/pangeacyber/pangea-go/pangea-sdk/v4 v4.0.0 => ../../pangea-sdk
