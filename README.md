# go-pangea


# Usage 
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
    "os"
	"go-pangea/pangea"
	"go-pangea/service/embargo"
)

func main() {
    token := os.Getenv("PANGEA_AUTH_TOKEN")
	embargocli := embargo.New(token, nil)

	input := &embargo.CheckInput{
		ISOCode: pangea.String("CU"),
	}
	checkOutput, _, err := embargocli.Check(context.Background(), input)

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(spew(checkOutput))
}

// spew prints the struct into a json string with tabs
func spew(jsonObj interface{}) string {
	b, _ := json.MarshalIndent(jsonObj, "", "\t")
	return string(b)
}  
```