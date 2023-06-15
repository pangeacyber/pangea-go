package utils

import (
	"errors"
	"fmt"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/authn"
	"github.com/sethvargo/go-password/password"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	EMAILFIELD     = "email"
	FULLNAMEFIELD  = "profile.full_name"
	FIRSTNAMEFIELD = "profile.first_name"

	LASTNAMEFIELD = "profile.last_name"

	NICKNAMEFIELD = "profile.nickname"

	PICTUREFIELD = "profile.image_url"

	EMAILVERIFIEDFIELD = "verified"

	IDPROVIDERFIELD = "id_provider"
)

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("file %v does not exist\n", filePath)
		return false
	}
	return true
}

func ValidateAndOpen(fileName string) (*os.File, error) {
	if !IsFileExist(fileName) {
		return nil, errors.New("mapping file does not exist")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ConvertValue(value interface{}, sourceType string, destType string) (interface{}, error) {
	// If source field is not defined then return as it is
	if sourceType == "" {
		return value, nil
	}
	// conversion
	// TODO - Add more conversion
	switch sourceType {
	case "string":
		switch destType {
		case "integer":
			return strconv.Atoi(value.(string))
		// TODO: regex is supported only for string
		case "regex":
			re := regexp.MustCompile(`(?i)(password|google|github)`)
			match := re.FindString(value.(string))
			return match, nil
		default:
			return value, nil
		}
	case "integer":
		switch destType {
		case "string":
			return strconv.Itoa(value.(int)), nil
		default:
			return value, nil
		}
	default:
		return nil, fmt.Errorf("unsupported type conversion: %s to %s", sourceType, destType)
	}
}

func ConvertMapToCreateUserRequest(rawUser map[string]interface{}) (*authn.UserCreateRequest, error) {
	user := new(authn.UserCreateRequest)
	// AuthO csv import has extra ' in the email field. So, we have to trim it
	user.Email = strings.Trim(convertInterfaceToString(rawUser[EMAILFIELD]), "'")
	user.Profile = &authn.ProfileData{
		"first_name": convertInterfaceToString(rawUser[FIRSTNAMEFIELD]),
		"last_name":  convertInterfaceToString(rawUser[LASTNAMEFIELD]),
		"image_url":  convertInterfaceToString(rawUser[PICTUREFIELD]),
		"nick_name":  convertInterfaceToString(rawUser[NICKNAMEFIELD]),
	}
	switch convertInterfaceToString(rawUser[IDPROVIDERFIELD]) {
	case authn.IDPGoogle:
		// TODO: add authenticator
		user.IDProvider = authn.IDPGoogle
	case authn.IDPGithub:
		// TODO: add authenticator
		user.IDProvider = authn.IDPGithub
	default:
		// default is password
		pass, err := createRandomPassword()
		if err != nil {
			return nil, err
		}
		user.Authenticator = pass

	}
	return user, nil
}

func convertInterfaceToString(data interface{}) string {
	if data == nil {
		return ""
	}
	switch data.(type) {
	case string:
		return data.(string)
	default:
		// unsupported type
		return ""
	}
}

func createRandomPassword() (string, error) {
	randomPass, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return "", err
	}
	return randomPass, nil
}
