package user_intel

import (
	"errors"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type PasswordStatus int

const (
	PSbreached     PasswordStatus = 0
	PSunbreached   PasswordStatus = 1
	PSinconclusive PasswordStatus = 2
)

func IsPasswordBreached(r *pangea.PangeaResponse[UserPasswordBreachedResult], h string) (PasswordStatus, error) {
	if r == nil {
		return PSinconclusive, errors.New("Response nil pointer")
	}

	if r.Result.RawData == nil {
		return PSinconclusive, errors.New("Need raw data to check if hash is breached. Send request with raw=true")
	}

	_, ok := r.Result.RawData[h]
	if ok {
		// If hash is present in raw data, it's because it was breached
		return PSbreached, nil
	} else {
		// If it's not present, should check if I have all breached hash
		// Server will return a maximum of 1000 hash, so if breached count is greater than that,
		// I can't conclude is password is or is not breached
		if len(r.Result.RawData) >= 1000 {
			return PSinconclusive, nil
		} else {
			return PSunbreached, nil
		}
	}
}
