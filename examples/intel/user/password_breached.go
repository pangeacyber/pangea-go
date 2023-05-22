// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/user_intel"
)

func main() {
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := user_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	// Set the password you would like to check
	password := "mypassword"
	// Calculate its hash, it could be sha256 or sha1
	hash := pangea.HashSHA256(password)
	// get the hash prefix, right know it should be just 5 characters
	hashPrefix := pangea.GetHashPrefix(hash, 5)

	ctx := context.Background()
	input := &user_intel.UserPasswordBreachedRequest{
		// should setup right hash_type here, sha256 or sha1
		HashType:   user_intel.HTsha265,
		HashPrefix: hashPrefix,
		Raw:        true,
		Verbose:    true,
		Provider:   "spycloud",
	}

	r, err := intelcli.PasswordBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	// This auxiliary function analyze service provider raw data to search for full hash in their registers
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
