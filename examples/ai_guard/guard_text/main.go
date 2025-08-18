package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ai_guard"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("PANGEA_AI_GUARD_TOKEN")
	if token == "" {
		log.Fatal("missing AI Guard API token")
	}

	config, err := pangea.NewConfig(option.WithToken(token), option.WithDomain(os.Getenv("PANGEA_DOMAIN")))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	aiGuardClient := ai_guard.New(config)

	input := &ai_guard.TextGuardRequest{Text: "what was pangea?"}
	out, err := aiGuardClient.GuardText(ctx, input)
	if err != nil {
		log.Fatal("failed to guard text")
	}

	fmt.Println(pangea.Stringify(out.Result))
}
