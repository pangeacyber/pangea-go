package imports

import (
	"context"
	"errors"
	"extensions/authn/internal/readers"
	"fmt"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/authn"
	"io"
	"strings"
	"time"
)

const (
	EmailField = "Email"
	FirstName  = ""
)

type Report struct {
	success map[string]interface{}
	failed  map[string]error
}

func convertMapToCreateUserRequest(rawUser map[string]interface{}) (*authn.UserCreateRequest, error) {
	user := new(authn.UserCreateRequest)
	// TODO - Fix this. Convert to builder
	user.Email = strings.Trim(rawUser[EmailField].(string), "'")
	//user.Profile = &authn.ProfileData{
	//	"NickName":   "",
	//	"FamilyName": "",
	//	"Name":       "",
	//}
	//bTrue := true
	//user.Verified = &bTrue
	fmt.Println("Email", user.Email)
	user.IDProvider = authn.IDPGoogle    // authn.IDPGoogle
	user.Authenticator = "My1s+Password" // testMy1s+Password
	return user, nil
}

func ImportUsers(token string, domain string, filePath string) error {
	if token == "" || domain == "" {
		return errors.New("token or domain is empty")
	}

	csvReader, err := readers.NewCSVReader(filePath, nil)
	if err != nil {
		return err
	}

	client := authn.New(&pangea.Config{
		Token:  token,
		Domain: domain,
	})
	report := Report{
		success: make(map[string]interface{}),
		failed:  make(map[string]error),
	}
	for {
		rawUser, err := csvReader.Next()
		if err == io.EOF {
			// Successfully process all users
			break
		}
		if err != nil {
			fmt.Printf("failed to read file=%s \n", err)
			continue
		}
		user, err := convertMapToCreateUserRequest(rawUser)
		if err != nil {
			fmt.Println("failed to convert map to user object")
			report.failed[user.Email] = err
			break
		}
		if err != nil {
			fmt.Println("Failed to read user from file, trying other users")
			report.failed[user.Email] = err
			continue
		}
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		fmt.Printf("creating user=%s \n", user.Email)
		resp, err := client.User.Create(ctx, *user)
		// TODO: Add error summary
		if err != nil {
			fmt.Printf("Failed to create user, trying other users, err=%s \n", err)
			fmt.Println(resp)
			report.failed[user.Email] = err
			continue
		}
		fmt.Println("Create user success. Result: " + pangea.Stringify(resp.Result))
		report.success[user.Email] = resp.Result
		break
	}
	fmt.Printf("completed import workflow, stats success:%d, failed:%d \n", len(report.success), len(report.failed))
	fmt.Println(report)
	return nil
}
