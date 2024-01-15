package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/store"
)

func main() {
	var t = time.Now().Format("20060102_150405")
	const filePath = "./testdata/testfile.pdf"
	var folder = "/examples/files/" + t

	// Load pangea token from environment variables
	token := os.Getenv("PANGEA_STORE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// create a new store client with pangea token and domain
	fmt.Println("Creating new folder...")
	client := store.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  120 * time.Second,
		Retry:              true,
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	})

	// Create a folder
	respCreate, err := client.FolderCreate(ctx, &store.FolderCreateRequest{
		Path: folder,
	})
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	folderID := respCreate.Result.Object.ID
	fmt.Printf("Folder create success. Folder ID: %s\n", folderID)

	// Upload a file with path as unique param
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("Uploading file with Path field...")
	respPut, err := client.Put(ctx,
		&store.PutRequest{
			Path:           path.Join(folder, "file_multipart_1"),
			TransferMethod: pangea.TMmultipart,
		},
		file)

	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Printf("Put file success. Object ID: %s\n", respPut.Result.Object.ID)
	fmt.Printf("Parent ID: %s\n", respPut.Result.Object.ParentID)

	// Upload a file with parent id and name and adding metadata and tags
	file, err = os.Open(filePath)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	var metadata = map[string]string{"field1": "value1", "field2": "value2"}
	var tags = []string{"tag1", "tag2"}

	fmt.Println("Uploading file with Name and ParentID...")
	respPut2, err := client.Put(ctx,
		&store.PutRequest{
			Name:           "file_multipart_2",
			ParentID:       folderID,
			TransferMethod: pangea.TMmultipart,
			Metadata:       metadata,
			Tags:           tags,
		},
		file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Printf("Put file success. Object ID: %s\n", respPut2.Result.Object.ID)

	// Update file with full metadata and tags
	fmt.Println("Updating object with metadata and tags...")
	respUpdate, err := client.Update(ctx, &store.UpdateRequest{
		ID:       respPut.Result.Object.ID,
		Metadata: metadata,
		Tags:     tags,
	})
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Printf("Updated item id %s successfully\n", respUpdate.Result.Object.ID)

	// Update file with add metadata and tags
	fmt.Println("Adding metadata and tags to a object...")
	var addMetadata = map[string]string{"field3": "value3"}
	var addTags = []string{"tag3"}

	respUpdate2, err := client.Update(ctx, &store.UpdateRequest{
		ID:          respPut2.Result.Object.ID,
		AddMetadata: addMetadata,
		AddTags:     addTags,
	})
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("Updated item id %s successfully\n", respUpdate2.Result.Object.ID)

	// Get archive as a multipart response
	fmt.Println("Getting archive as multipart...")
	respGetArchive, err := client.GetArchive(ctx, &store.GetArchiveRequest{
		Ids:            []string{folderID},
		Format:         store.AFzip,
		TransferMethod: pangea.TMmultipart,
	})

	fmt.Printf("Archive download has %d file(s)\n", len(respGetArchive.AttachedFiles))
	for _, af := range respGetArchive.AttachedFiles {
		// Save file. In this case should be just one archive anyway
		err := af.Save(pangea.AttachedFileSaveInfo{
			Folder: "./download/archive/",
		})
		if err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}

	// Get archive as a download url
	fmt.Println("Getting archive as dest-url...")
	respGetArchive2, err := client.GetArchive(ctx, &store.GetArchiveRequest{
		Ids:            []string{folderID},
		Format:         store.AFzip,
		TransferMethod: pangea.TMdestURL,
	})

	fmt.Printf("Archive download has %d file(s)\n", len(respGetArchive2.AttachedFiles))

	// Download file
	fmt.Println("Download archive file from url...")
	attachedFile, err := client.DownloadFile(ctx, *respGetArchive2.Result.DestURL)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("Download success. Saving file...")

	err = attachedFile.Save(pangea.AttachedFileSaveInfo{
		Folder: "./download/archive/",
	})

	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("Save success")

	// Create share link...
	fmt.Println("Creating share link...")
	// Create authenticator methods to access the share link
	authenticators := []store.Authenticator{store.Authenticator{
		AuthType:    store.ATpassword,
		AuthContext: "somepassword",
	}}

	ll := []store.ShareLinkCreateItem{store.ShareLinkCreateItem{
		// Set targets to the share link
		Targets:        []string{folderID},
		LinkType:       store.LTeditor,
		Authenticators: authenticators,
		MaxAccessCount: pangea.Int(3),
	}}
	respCreateLink, err := client.ShareLinkCreate(ctx, &store.ShareLinkCreateRequest{
		Links: ll,
	})

	links := respCreateLink.Result.ShareLinkObjects
	link := links[0]

	fmt.Printf("Share link created: %s\n", link.Link)

	// Get share link
	fmt.Println("Getting an already created share link...")
	respGetLink, err := client.ShareLinkGet(ctx, &store.ShareLinkGetRequest{
		ID: link.ID,
	})
	fmt.Printf("Get success: %s\n", respGetLink.Result.ShareLinkObject.Link)

	// List share link
	fmt.Println("Getting a list of links...")
	respListLink, err := client.ShareLinkList(ctx, &store.ShareLinkListRequest{})
	fmt.Printf("Got %d link(s)\n", respListLink.Result.Count)

	// Delete share link
	fmt.Println("Deleting share link...")
	respDeleteLink, err := client.ShareLinkDelete(ctx, &store.ShareLinkDeleteRequest{
		Ids: []string{link.ID},
	})

	fmt.Printf("Deleted %d link(s)\n", len(respDeleteLink.Result.ShareLinkObjects))

	// List files in folder
	fmt.Println("Listing objects in folder...")

	// Create a ListFilter an set its possible values
	listFilter := store.NewFilterList()
	listFilter.Folder().Set(pangea.String(folder))

	respList, err := client.List(ctx, &store.ListRequest{
		Filter: listFilter.Filter(),
	})

	fmt.Printf("Got %d object(s)\n", len(respList.Result.Objects))

}
