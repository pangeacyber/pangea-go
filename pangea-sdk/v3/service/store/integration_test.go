//go:build integration

package store_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/store"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Develop
)

var PDF_FILEPATH = "./testdata/testfile.pdf"
var timeNow = time.Now()
var TIME = timeNow.Format("20060102_150405")
var FOLDER_DELETE = "/sdk_tests/delete/" + TIME
var FOLDER_FILES = "/sdk_tests/files/" + TIME
var METADATA = map[string]string{"field1": "value1", "field2": "value2"}
var ADD_METADATA = map[string]string{"field3": "value3"}
var TAGS = []string{"tag1", "tag2"}
var ADD_TAGS = []string{"tag3"}

func storeIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func Test_Integration_Folder(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	cfg := storeIntegrationCfg(t)
	client := store.New(cfg)

	input := &store.FolderCreateRequest{
		Path: FOLDER_DELETE,
	}

	out, err := client.FolderCreate(ctx, input)
	if err != nil {
		fmt.Println(reflect.TypeOf(err))
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Object.ID)
	assert.NotEmpty(t, out.Result.Object.Name)
	assert.NotEmpty(t, out.Result.Object.CreatedAt)
	assert.NotEmpty(t, out.Result.Object.UpdatedAt)
	assert.Equal(t, out.Result.Object.Type, "folder")
	id := out.Result.Object.ID

	input2 := &store.DeleteRequest{
		ID: id,
	}
	rDel, err := client.Delete(ctx, input2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, *rDel.Status, "Success")
	assert.Equal(t, rDel.Result.Count, 1)
}

func Test_Integration_PutTransferMethodPostURL(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := storeIntegrationCfg(t)
	client := store.New(cfg)

	name := TIME + "_file_post_url"

	input := &store.PutRequest{
		Name: name,
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
	}

	file, err := os.Open(PDF_FILEPATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := client.Put(ctx, input, file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Object.ID)
	assert.NotEmpty(t, out.Result.Object.Name)
}

func Test_Integration_PutTransferMethodMultipart(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	cfg := storeIntegrationCfg(t)
	client := store.New(cfg)

	name := TIME + "_file_multipart"

	input := &store.PutRequest{
		Name: name,
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMmultipart,
		},
	}

	file, err := os.Open(PDF_FILEPATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := client.Put(ctx, input, file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Object.ID)
	assert.NotEmpty(t, out.Result.Object.Name)
}

func Test_Integration_SplitUpload_Put(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	client := store.New(cfg)

	name := TIME + "_file_split_put_url"

	input := &store.PutRequest{
		Name: name,
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMputURL,
		},
	}

	file, err := os.Open(PDF_FILEPATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.RequestUploadURL(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err.Error())
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.AcceptedResult)
	assert.NotEmpty(t, resp.AcceptedResult.PutURL)
	assert.Empty(t, resp.AcceptedResult.PostURL)
	assert.Empty(t, resp.AcceptedResult.PostFormData)

	url := resp.AcceptedResult.PutURL

	fd := pangea.FileData{
		File: file,
		Name: "someName",
	}

	uploader := store.NewFileUploader()
	err = uploader.UploadFile(ctx, url, pangea.TMputURL, fd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var pr *pangea.PangeaResponse[any]
	i := 0

	for i < 24 {
		// Wait until result should be ready
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByID(ctx, *resp.RequestID, &store.PutResult{})
		if err == nil {
			break
		}
		i++
	}
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.NotNil(t, pr.Result)

	_, ok := (*pr.Result).(*store.PutResult)
	assert.True(t, ok)

}

func Test_Integration_LifeCycle(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	client := store.New(cfg)

	// Create a folder
	respCreate, err := client.FolderCreate(ctx, &store.FolderCreateRequest{
		Path: FOLDER_FILES,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, respCreate)
	assert.NotNil(t, respCreate.Result)
	assert.NotEmpty(t, respCreate.Result.Object.ID)
	folderID := respCreate.Result.Object.ID

	// Upload a file with path as unique param
	file, err := os.Open(PDF_FILEPATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	respPut, err := client.Put(ctx,
		&store.PutRequest{
			Path: FOLDER_FILES + "/" + TIME + "_file_multipart_1",
			TransferRequest: pangea.TransferRequest{
				TransferMethod: pangea.TMmultipart,
			},
		},
		file)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, respPut)
	assert.NotNil(t, respPut.Result)
	assert.Equal(t, folderID, respPut.Result.Object.ParentID)
	assert.Empty(t, respPut.Result.Object.Metadata)
	assert.Empty(t, respPut.Result.Object.Tags)
	assert.Empty(t, respPut.Result.Object.MD5)
	assert.Empty(t, respPut.Result.Object.SHA512)
	assert.NotEmpty(t, respPut.Result.Object.SHA256)

	// Upload a file with parent id and name
	file, err = os.Open(PDF_FILEPATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	respPut2, err := client.Put(ctx,
		&store.PutRequest{
			Name:     TIME + "_file_multipart_2",
			ParentID: folderID,
			TransferRequest: pangea.TransferRequest{
				TransferMethod: pangea.TMmultipart,
			},
			Metadata: METADATA,
			Tags:     TAGS,
		},
		file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, respPut2)
	assert.NotNil(t, respPut2.Result)
	assert.Equal(t, folderID, respPut2.Result.Object.ParentID)
	assert.Equal(t, store.Metadata(METADATA), respPut2.Result.Object.Metadata)
	assert.Equal(t, store.Tags(TAGS), respPut2.Result.Object.Tags)
	assert.Empty(t, respPut2.Result.Object.MD5)
	assert.Empty(t, respPut2.Result.Object.SHA512)
	assert.NotEmpty(t, respPut2.Result.Object.SHA256)

	// Update file with full metadata and tags
	respUpdate, err := client.Update(ctx, &store.UpdateRequest{
		ID:       respPut.Result.Object.ID,
		Metadata: METADATA,
		Tags:     TAGS,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, respUpdate)
	assert.NotNil(t, respUpdate.Result)
	assert.Equal(t, store.Metadata(METADATA), respUpdate.Result.Object.Metadata)
	assert.Equal(t, store.Tags(TAGS), respUpdate.Result.Object.Tags)

	// Update file with add metadata and tags
	respUpdate2, err := client.Update(ctx, &store.UpdateRequest{
		ID:          respPut2.Result.Object.ID,
		AddMetadata: ADD_METADATA,
		AddTags:     ADD_TAGS,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, respUpdate2)
	assert.NotNil(t, respUpdate2.Result)

	// Get archive
	respGetArchive, err := client.GetArchive(ctx, &store.GetArchiveRequest{
		Ids:            []string{folderID},
		Format:         store.AFzip,
		TransferMethod: pangea.TMmultipart,
	})

	assert.NoError(t, err)
	assert.NotNil(t, respGetArchive)
	assert.NotNil(t, respGetArchive.Result)
	assert.Nil(t, respGetArchive.Result.DestURL)
	assert.Equal(t, len(respGetArchive.AttachedFiles), 1)

	for _, af := range respGetArchive.AttachedFiles {
		err := af.Save(pangea.AttachedFileSaveInfo{
			Folder: "./download/archive/",
		})
		assert.NoError(t, err)
	}

	respGetArchive2, err := client.GetArchive(ctx, &store.GetArchiveRequest{
		Ids:            []string{folderID},
		Format:         store.AFzip,
		TransferMethod: pangea.TMdestURL,
	})

	assert.NoError(t, err)
	assert.NotNil(t, respGetArchive2)
	assert.NotNil(t, respGetArchive2.Result)
	assert.NotNil(t, respGetArchive2.Result.DestURL)
	assert.Equal(t, len(respGetArchive2.AttachedFiles), 0)

	// Download file
	attachedFile, err := client.DownloadFile(ctx, *respGetArchive2.Result.DestURL)
	assert.NoError(t, err)
	assert.NotNil(t, attachedFile)
	assert.NotEmpty(t, attachedFile.File)
	assert.NotEmpty(t, attachedFile.Filename)
	assert.NotEmpty(t, attachedFile.ContentType)

	// Create share link
	authenticators := []store.Authenticator{store.Authenticator{
		AuthType:    store.ATpassword,
		AuthContext: "somepassword",
	}}
	ll := []store.ShareLinkCreateItem{store.ShareLinkCreateItem{
		Targets:        []string{folderID},
		LinkType:       store.LTeditor,
		Authenticators: authenticators,
		MaxAccessCount: pangea.Int(3),
		Message:        "share message",
		Title:          "share title",
	}}
	respCreateLink, err := client.ShareLinkCreate(ctx, &store.ShareLinkCreateRequest{
		Links: ll,
	})

	assert.NoError(t, err)
	assert.NotNil(t, respCreateLink)
	assert.NotNil(t, respCreateLink.Result)

	links := respCreateLink.Result.ShareLinkObjects
	assert.Equal(t, len(links), 1)

	link := links[0]
	assert.Equal(t, link.AccessCount, 0)
	assert.Equal(t, link.MaxAccessCount, 3)
	assert.Equal(t, len(link.Authenticators), 1)
	assert.Equal(t, string(link.Authenticators[0].AuthType), string(store.ATpassword))
	assert.NotEmpty(t, link.Link)
	assert.NotEmpty(t, link.ID)
	assert.Equal(t, len(link.Targets), 1)

	// Get share link
	respGetLink, err := client.ShareLinkGet(ctx, &store.ShareLinkGetRequest{
		ID: link.ID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, respGetLink)
	assert.NotNil(t, respGetLink.Result)
	assert.Equal(t, respGetLink.Result.ShareLinkObject, link)

	// List share link
	respListLink, err := client.ShareLinkList(ctx, &store.ShareLinkListRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, respListLink)
	assert.True(t, respListLink.Result.Count > 0)
	assert.True(t, len(respListLink.Result.ShareLinkObjects) > 0)

	// Delete share link
	respDeleteLink, err := client.ShareLinkDelete(ctx, &store.ShareLinkDeleteRequest{
		Ids: []string{link.ID},
	})

	assert.NoError(t, err)
	assert.NotNil(t, respDeleteLink)
	assert.Equal(t, len(respDeleteLink.Result.ShareLinkObjects), 1)

	// List files in folder
	listFilter := store.NewFilterList()
	listFilter.Folder().Set(pangea.String(FOLDER_FILES))

	respList, err := client.List(ctx, &store.ListRequest{
		Filter: listFilter.Filter(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, respList)
	assert.NotNil(t, respList.Result)
	assert.Equal(t, respList.Result.Count, 2)
	assert.Equal(t, len(respList.Result.Objects), 2)
}
