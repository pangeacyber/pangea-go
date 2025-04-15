package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ai_guard"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("PANGEA_AI_GUARD_TOKEN")
	if token == "" {
		log.Fatal("missing AI Guard API token")
	}

	aiGuardClient := ai_guard.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	input := &ai_guard.TextGuardRequest{Text: "what was pangea?"}
	out, err := aiGuardClient.GuardText(ctx, input)
	if err != nil {
		log.Fatal("failed to guard text")
	}

	fmt.Println(pangea.Stringify(out.Result))
}
