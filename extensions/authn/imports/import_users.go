package imports

import (
	"context"
	"errors"
	"extensions/authn/internal/importio"
	"extensions/authn/internal/models"
	"fmt"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/authn"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
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
	user.IDProvider = authn.IDPGoogle // authn.IDPGoogle

	randomPass, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return nil, err
	}

	user.Authenticator = randomPass //"My1s+Password" // testMy1s+Password
	return user, nil
}

func ImportUsers(token string, domain string, filePath string, mappingFile string, isDryRun bool) error {
	// Get global logger
	logger := zap.L()
	if token == "" || domain == "" {
		return errors.New("token or domain is empty")
	}

	// Read mapping file
	var mappings *models.Mappings
	var err error
	if mappingFile == "" {
		mappings, err = models.NewMappings(mappingFile)
		if err != nil {
			logger.Error("failed to open mapping file", zap.String("err", err.Error()))
			return err
		}
	}

	csvReader, err := importio.NewCSVReader(filePath, mappings)
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

	t := time.Now()
	uniqueID := t.UnixNano()
	outputFileName := fmt.Sprintf("success_userinfo_%d", uniqueID)
	csvWriter, err := importio.NewCSVWriter(outputFileName, []string{"Id", "Email", "Password"})
	if err != nil {
		logger.Error("Failed to write password to a file", zap.Error(err))
		return err
	}
	defer csvWriter.Close()

	for {
		rawUser, err := csvReader.Next()
		if err == io.EOF {
			// Successfully process all users
			break
		}
		if err != nil {
			logger.Error("failed to read user record", zap.Error(err))
			continue
		}
		user, err := convertMapToCreateUserRequest(rawUser)
		if err != nil {
			logger.Error("failed to build user profile object from raw format", zap.Error(err))
			report.failed[user.Email] = err
			break
		}
		if isDryRun {
			report.success[user.Email] = "Dry Run"
			continue
		}
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		logger.Info("creating user", zap.String("email", user.Email))
		resp, err := client.User.Create(ctx, *user)
		if err != nil {
			logger.Info("Failed to create user, trying other users", zap.Error(err))
			report.failed[user.Email] = err
			continue
		}
		logger.Info("successfully created a user", zap.String("Result", pangea.Stringify(resp.Result)))
		report.success[user.Email] = resp.Result
		csvWriter.Write([]string{resp.Result.ID, resp.Result.Email, user.Authenticator})
	}
	logger.Info("completed import workflow, stats", zap.Int("success", len(report.success)),
		zap.Int("failed", len(report.failed)))
	return nil
}
