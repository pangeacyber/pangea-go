package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/share"
)

func main() {
	var t = time.Now().Format("20060102_150405")
	var path = "/sdk_example/delete/" + t

	// Load pangea token from environment variables
	token := os.Getenv("PANGEA_SHARE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	// create a new store client with pangea token and domain
	client := share.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	// Create a FolderCreateRequest and set the path of the folder to be created
	input := &share.FolderCreateRequest{
		Path: path,
	}

	fmt.Printf("Let's create a folder: %s\n", path)
	// Send the CreateRequest
	out, err := client.FolderCreate(ctx, input)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	id := out.Result.Object.ID
	fmt.Printf("Folder created. ID: %s.\n", id)

	fmt.Printf("Let's create this folder now\n")
	// Create a DeleteRequest and set the ID of the item to be deleted
	input2 := &share.DeleteRequest{
		ID: id,
	}

	// Send the DeleteRequest
	rDel, err := client.Delete(ctx, input2)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Printf("Folder deleted. Deleted %d items.\n", rDel.Result.Count)
}
