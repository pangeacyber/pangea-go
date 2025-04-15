// Here we'll see how to manage user invites
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/authn"
)

func main() {
	// Set up variables to be used
	rand.Seed(time.Now().UnixNano())
	RANDOM_VALUE := strconv.Itoa(rand.Intn(10000000))
	USER_EMAIL := "user.email+goexample" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_INVITE_1 := "user.invite1+goexample" + RANDOM_VALUE + "@pangea.cloud"
	EMAIL_INVITE_2 := "user.invite2+goexample" + RANDOM_VALUE + "@pangea.cloud"
	// End set up variables

	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	// Get token
	token := os.Getenv("PANGEA_AUTHN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	// Create config and client
	client := authn.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	// Requests examples...
	fmt.Println("Invite user 1...")
	input := authn.UserInviteRequest{
		Inviter:  USER_EMAIL,
		Email:    EMAIL_INVITE_1,
		Callback: "https://someurl.com/callbacklink",
		State:    "Somestate",
	}
	resp, err := client.User.Invite(ctx, input)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Invite user 1 success. Result: " + pangea.Stringify(resp.Result))

	fmt.Println("Invite user 2...")
	input2 := authn.UserInviteRequest{
		Inviter:  USER_EMAIL,
		Email:    EMAIL_INVITE_2,
		Callback: "https://someurl.com/callbacklink",
		State:    "Somestate",
	}
	resp2, err := client.User.Invite(ctx, input2)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Invite user 2 success. Result: " + pangea.Stringify(resp2.Result))

	fmt.Println("List invites...")
	resp4, err := client.User.Invites.List(ctx, authn.UserInviteListRequest{})
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("List success. There is ", len(resp4.Result.Invites), " invites")

	fmt.Println("Delete invite user 2...")
	input3 := authn.UserInviteDeleteRequest{
		ID: resp2.Result.ID,
	}
	_, err = client.User.Invites.Delete(ctx, input3)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Delete invite user 2 success")

	fmt.Println("List invites...")
	resp4, err = client.User.Invites.List(ctx, authn.UserInviteListRequest{})
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("List success. There is ", len(resp4.Result.Invites), " invites")

}
