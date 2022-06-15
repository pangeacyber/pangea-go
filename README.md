# go-pangea


# Usage 
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
)

func main() {
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	embargocli := embargo.New(pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "dev",
		},
	})

	input := &embargo.CheckInput{
		ISOCode: pangea.String("CU"),
	}
	ctx := context.Background()
	checkOutput, _, err := embargocli.Check(ctx, input)
	if err != nil {
		fmt.Println(err.Error())
	}
	if checkOutput != nil {
		fmt.Println(spew(checkOutput))
	}
}

// spew prints the struct into a json string with tabs
func spew(jsonObj interface{}) string {
	b, _ := json.MarshalIndent(jsonObj, "", "\t")
	return string(b)
}
```