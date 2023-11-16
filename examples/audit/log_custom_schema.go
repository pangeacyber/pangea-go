package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"examples/audit/util"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_AUDIT_CUSTOM_SCHEMA_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditcli, err := audit.New(
		&pangea.Config{
			Token:  token,
			Domain: os.Getenv("PANGEA_DOMAIN"),
		},
		audit.WithCustomSchema(util.CustomSchemaEvent{}),
	)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	msg := "go-sdk-custom-schema-no-signed"
	var event = &util.CustomSchemaEvent{
		Message:       msg,
		FieldInt:      1,
		FieldBool:     true,
		FieldStrShort: "no-signed",
		FieldStrLong:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed lacinia, orci eget commodo commodo non.",
		FieldTime:     pangea.Time(time.Now()),
	}

	fmt.Printf("Logging: %s\n", event.Message)

	lr, err := auditcli.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*util.CustomSchemaEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
}
