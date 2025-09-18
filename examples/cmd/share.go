package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/share"
	"github.com/spf13/cobra"
)

func init() {
	shareCmd := &cobra.Command{
		Use:   "share",
		Short: "Share examples",
	}

	shareCliCmd := &cobra.Command{
		Use:   "cli",
		Short: "A command line utility to create share links",
		RunE:  shareCli,
	}
	shareCliCmd.Flags().StringP("input", "i", "", "Local path to upload.")
	err := shareCliCmd.MarkFlagRequired("input")
	if err != nil {
		log.Fatalf("failed to mark flag as required: %v", err)
	}
	shareCliCmd.Flags().String("dest", "/", "Destination path in Share.")
	shareCliCmd.Flags().String("email", "", "Email address to protect the share link with.")
	shareCliCmd.Flags().String("phone", "", "Phone number to protect the share link with.")
	shareCliCmd.Flags().String("password", "", "Password to protect the share link with.")
	shareCliCmd.Flags().String("link-type", string(share.LTdownload), "Type of link.")
	shareCmd.AddCommand(shareCliCmd)

	folderCreateAndDeleteCmd := &cobra.Command{
		Use:  "folder_create_and_delete",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return folderCreateAndDelete(cmd)
		},
	}
	shareCmd.AddCommand(folderCreateAndDeleteCmd)

	ExamplesCmd.AddCommand(shareCmd)
}

func folderCreateAndDelete(cmd *cobra.Command) error {
	var t = time.Now().Format("20060102_150405")
	var path = "/sdk_example/delete/" + t

	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SHARE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	// create a new store client with pangea token and domain
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := share.New(config)

	// Create a FolderCreateRequest and set the path of the folder to be created
	input := &share.FolderCreateRequest{
		Folder: path,
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
	return nil
}

func getFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var files []string
	if info.IsDir() {
		fileInfos, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, fileInfo := range fileInfos {
			files = append(files, filepath.Join(path, fileInfo.Name()))
		}
	} else {
		files = append(files, path)
	}

	return files, nil
}

func shareCli(cmd *cobra.Command, args []string) error {
	input, _ := cmd.Flags().GetString("input")
	dest, _ := cmd.Flags().GetString("dest")
	email, _ := cmd.Flags().GetString("email")
	phone, _ := cmd.Flags().GetString("phone")
	password, _ := cmd.Flags().GetString("password")
	linkType, _ := cmd.Flags().GetString("link-type")

	if email == "" && phone == "" && password == "" {
		return fmt.Errorf("at least one of --email, --phone, or --password must be provided")
	}

	authenticators := []share.Authenticator{}
	if email != "" {
		authenticators = append(authenticators, share.Authenticator{
			AuthType:    share.ATemailOTP,
			AuthContext: email,
		})
	}
	if phone != "" {
		authenticators = append(authenticators, share.Authenticator{
			AuthType:    share.ATsmsOTP,
			AuthContext: phone,
		})
	}
	if password != "" {
		authenticators = append(authenticators, share.Authenticator{
			AuthType:    share.ATpassword,
			AuthContext: password,
		})
	}

	token := os.Getenv("PANGEA_SHARE_TOKEN")
	if token == "" {
		log.Fatal("missing `PANGEA_SHARE_TOKEN` environment variable")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := share.New(config)

	files, err := getFiles(input)
	if err != nil {
		return err
	}

	var objectIDs []string
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		uploadResponse, err := client.Put(ctx, &share.PutRequest{Folder: fmt.Sprintf("%s/%s", dest, filepath.Base(file))}, f)
		if err != nil {
			// We are returning an error, so we can ignore the error on Close.
			_ = f.Close()
			return err
		}
		objectIDs = append(objectIDs, uploadResponse.Result.Object.ID)

		if err = f.Close(); err != nil {
			return err
		}
	}

	linkResponse, err := client.ShareLinkCreate(ctx, &share.ShareLinkCreateRequest{
		Links: []share.ShareLinkCreateItem{
			{
				Targets:        objectIDs,
				LinkType:       share.LinkType(linkType),
				Authenticators: authenticators,
			},
		},
	})
	if err != nil {
		return err
	}

	fmt.Println(linkResponse.Result.ShareLinkObjects[0].Link)
	return nil
}
