module examples/embargo

go 1.23

require github.com/pangeacyber/pangea-go/pangea-sdk/v5 v5.3.0

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/pangeacyber/pangea-go/pangea-sdk/v5 v5.0.0 => ../../pangea-sdk
