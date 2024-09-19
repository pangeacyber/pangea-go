package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/share"
	"github.com/spf13/cobra"
)

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

func main() {
	var input string
	var dest string
	var email string
	var phone string
	var password string
	var linkType string

	rootCmd := &cobra.Command{
		Long: "An example command line utility that creates an email+code, " +
			"SMS+code, or password-secured download/upload/editor share-link for " +
			"a given file or for each file in a given directory.",
		Version: "0.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			if email == "" && phone == "" && password == "" {
				fmt.Fprintln(os.Stderr, "At least one of --email, --phone, or --password must be provided.")
				os.Exit(1)
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

			client := share.New(&pangea.Config{
				Token:  token,
				Domain: os.Getenv("PANGEA_DOMAIN"),
			})

			files, err := getFiles(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			var objectIDs []string
			for _, file := range files {
				f, err := os.Open(file)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				defer f.Close()

				uploadResponse, err := client.Put(ctx, &share.PutRequest{Folder: fmt.Sprintf("%s/%s", dest, filepath.Base(file))}, f)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				objectIDs = append(objectIDs, uploadResponse.Result.Object.ID)
			}

			linkResponse, err := client.ShareLinkCreate(ctx, &share.ShareLinkCreateRequest{
				Links: []share.ShareLinkCreateItem{share.ShareLinkCreateItem{
					Targets:        objectIDs,
					LinkType:       share.LinkType(linkType),
					Authenticators: authenticators,
				}},
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println(linkResponse.Result.ShareLinkObjects[0].Link)
		},
	}

	rootCmd.Flags().StringVarP(&input, "input", "i", "", "Local path to upload.")
	rootCmd.MarkFlagRequired("input")

	rootCmd.Flags().StringVar(&dest, "dest", "/", "Destination path in Share.")

	rootCmd.Flags().StringVar(&email, "email", "", "Email address to protect the share link with.")
	rootCmd.Flags().StringVar(&phone, "phone", "", "Phone number to protect the share link with.")
	rootCmd.Flags().StringVar(&password, "password", "", "Password to protect the share link with.")

	rootCmd.Flags().StringVar(&linkType, "Link type", string(share.LTdownload), "Type of link.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
