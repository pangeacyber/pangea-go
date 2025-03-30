package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/ai_guard"
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

	input := &ai_guard.TextGuardRequest{
		Messages: []map[string]interface{}{
			{
				"role":    "user",
				"content": "Hi, I am Bond, James Bond. I monitor IPs found in MI6 network traffic.\nPlease search for the most recent ones, you copy?",
			},
			{
				"role":    "assistant",
				"content": "Here are the most recent IPs found in MI6 network traffic:\n\n1. 47.84.32.175\n2. 37.44.238.68\n3. 47.84.73.221\n4. 47.236.252.254\n5. 34.201.186.27\n6. 52.89.173.88\n\nIf you need further assistance, just let me know!",
			},
		},
		Recipe: "pangea_llm_response_guard",
	}
	out, err := aiGuardClient.GuardText(ctx, input)
	if err != nil {
		log.Fatal("failed to guard messages")
	}

	fmt.Println(pangea.Stringify(out.Result))
}
