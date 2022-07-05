# Pangea Go SDK

# Usage
```go
// embargo check is an example of how to use the check method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
)

func main() {
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("EMBARGO_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	embargocli, err := embargo.New(&pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "aws",
		},
		CfgToken: configID,
	})
	if err != nil {
		log.Fatal("failed to create embargo client")
	}

	ctx := context.Background()
	input := &embargo.CheckInput{
		ISOCode: pangea.String("CU"),
	}

	checkOutput, _, err := embargocli.Check(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(checkOutput))
}
```

# Contributing

Currently, the setup scripts only have support for Mac/ZSH environments.
Future support is incoming.

To install our linters, simply run `./dev/setup_repo.sh`
These linters will run on every `git commit` operation.
