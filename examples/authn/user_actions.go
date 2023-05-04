// vault sign is an example of how to use the sign/verify methods
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/authn"
)

func main() {
	// Set up variables to be used
	rand.Seed(time.Now().UnixNano())
	RANDOM_VALUE := strconv.Itoa(rand.Intn(10000000))
	USER_EMAIL := "user.email+goexample" + RANDOM_VALUE + "@pangea.cloud"
	PROFILE_INITIAL := &authn.ProfileData{
		"name":    "User name",
		"country": "Argentina",
	}
	PROFILE_UPDATE := authn.ProfileData{"age": "18"}
	PASSWORD_OLD := "My1s+Password"
	PASSWORD_NEW := "My1s+Password_new"
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
	// User create...
	fmt.Println("Creating user...")
	input := authn.UserCreateRequest{
		Email:         USER_EMAIL,
		Authenticator: PASSWORD_OLD,
		IDProvider:    authn.IDPPassword,
		Profile:       PROFILE_INITIAL,
	}
	resp, err := client.User.Create(ctx, input)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	USER_ID := resp.Result.ID
	fmt.Println("Create user success. Result: " + pangea.Stringify(resp.Result))

	// User login
	fmt.Println("\n\nUser login...")
	input2 := authn.UserLoginPasswordRequest{
		Email:    USER_EMAIL,
		Password: PASSWORD_OLD,
	}
	resp2, err := client.User.Login.Password(ctx, input2)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("User login success. Result: " + pangea.Stringify(resp2.Result))

	// User password change
	fmt.Println("\n\nUser password change..")
	input3 := authn.ClientPasswordChangeRequest{
		Token:       resp2.Result.ActiveToken.Token,
		OldPassword: PASSWORD_OLD,
		NewPassword: PASSWORD_NEW,
	}
	_, err = client.Client.Password.Change(ctx, input3)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("User password change success")

	// User profile get by email
	fmt.Println("User profile get by email...")
	input4 := authn.UserProfileGetRequest{
		Email: pangea.String(USER_EMAIL),
	}
	resp4, err := client.User.Profile.Get(ctx, input4)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Get profile success. Profile: " + pangea.Stringify(resp4.Result.Profile))

	// User profile get by id
	fmt.Println("Get profile by ID... %s", USER_ID)
	input5 := authn.UserProfileGetRequest{
		ID: pangea.String(USER_ID),
	}
	resp5, err := client.User.Profile.Get(ctx, input5)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Get profile success. Profile: " + pangea.Stringify(resp5.Result.Profile))

	// User profile update
	fmt.Println("User profile update...")
	input6 := authn.UserProfileUpdateRequest{
		Email:   pangea.String(USER_EMAIL),
		Profile: PROFILE_UPDATE,
	}
	resp6, err := client.User.Profile.Update(ctx, input6)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Update profile success. Profile: " + pangea.Stringify(resp6.Result.Profile))

	// User update
	fmt.Println("User update...")
	input7 := authn.UserUpdateRequest{
		Email:      pangea.String(USER_EMAIL),
		Disabled:   pangea.Bool(false),
		RequireMFA: pangea.Bool(false),
	}
	resp7, err := client.User.Update(ctx, input7)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("Update user success. Result: " + pangea.Stringify(resp7.Result))

	// User list
	fmt.Println("User list...")
	input8 := authn.UserListRequest{}
	resp8, err := client.User.List(ctx, input8)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("User list success. There is %d users", resp8.Result.Count)

	// User delete
	fmt.Println("User delete...")
	input9 := authn.UserDeleteRequest{
		Email: USER_EMAIL,
	}
	_, err = client.User.Delete(ctx, input9)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("User delete success")

	// User list
	fmt.Println("User list...")
	input10 := authn.UserListRequest{}
	resp10, err := client.User.List(ctx, input10)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return
	}
	fmt.Println("User list success. There is %d users", resp10.Result.Count)
}
