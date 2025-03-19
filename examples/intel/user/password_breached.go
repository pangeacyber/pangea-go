// Example of how to check if a password has been exposed/breached
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/user_intel"
)

func main() {
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := user_intel.New(&pangea.Config{
		Token:           token,
		BaseURLTemplate: os.Getenv("PANGEA_URL_TEMPLATE"),
	})

	// Set the password you would like to check
	// Observe proper safety with passwords, do not check them into source control etc.
	password := "mypassword"
	// Calculate its hash, it could be sha256, sha512 or sha1
	hash := pangea.HashSHA256(password)
	// get the hash prefix, just the first 5 characters
	hashPrefix := pangea.GetHashPrefix(hash, 5)

	ctx := context.Background()
	input := &user_intel.UserPasswordBreachedRequest{
		// set the right hash_type here, sha256, sha512 or sha1
		HashType:   user_intel.HTsha265,
		HashPrefix: hashPrefix,
		Raw:        pangea.Bool(true),
		Verbose:    pangea.Bool(true),
		Provider:   "spycloud",
	}

	r, err := intelcli.PasswordBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	// IsPasswordBreached is a helper function that can simplify searching the response's raw data for the full hash
	s, err := user_intel.IsPasswordBreached(r, hash)
	if err != nil {
		log.Fatal(err)
	}

	if s == user_intel.PSbreached {
		fmt.Printf("Password '%s' has been breached.\n", password)
	} else if s == user_intel.PSunbreached {
		fmt.Printf("Password '%s' has not been breached.\n", password)
	} else if s == user_intel.PSinconclusive {
		fmt.Printf("Not enough information to confirm if password '%s' has been or has not been breached.\n", password)
	} else {
		fmt.Printf("Unknown status: %d.\n", s)
	}
}
