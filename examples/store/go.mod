module examples/redact

go 1.19

require github.com/pangeacyber/pangea-go/pangea-sdk/v3 v3.6.0

require (
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

replace github.com/pangeacyber/pangea-go/pangea-sdk/v3 v3.6.0 => ../../pangea-sdk/v3
