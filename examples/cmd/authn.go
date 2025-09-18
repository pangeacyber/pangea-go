package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/authn"
	"github.com/spf13/cobra"
)

func init() {
	authnCmd := &cobra.Command{
		Use:   "authn",
		Short: "AuthN examples",
	}

	userActionsCmd := &cobra.Command{
		Use:  "user_actions",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return authnUserActions(cmd)
		},
	}
	authnCmd.AddCommand(userActionsCmd)

	inviteActionsCmd := &cobra.Command{
		Use:  "invite_actions",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return authnInviteActions(cmd)
		},
	}
	authnCmd.AddCommand(inviteActionsCmd)

	ExamplesCmd.AddCommand(authnCmd)
}

func authnInviteActions(cmd *cobra.Command) error {
	// Set up variables to be used
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
	config, err := pangea.NewConfig(option.WithToken(token), option.WithDomain(os.Getenv("PANGEA_DOMAIN")))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := authn.New(config)

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
		return err
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
		return err
	}
	fmt.Println("Invite user 2 success. Result: " + pangea.Stringify(resp2.Result))

	fmt.Println("List invites...")
	resp4, err := client.User.Invites.List(ctx, authn.UserInviteListRequest{})
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
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
		return err
	}
	fmt.Println("Delete invite user 2 success")

	fmt.Println("List invites...")
	resp4, err = client.User.Invites.List(ctx, authn.UserInviteListRequest{})
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Println("List success. There is ", len(resp4.Result.Invites), " invites")
	return nil
}

var CB_URI = "https://someurl.com/callbacklink"

func flowHandlePasswordPhase(ctx context.Context, client *authn.AuthN, flow_id, password string) *authn.FlowUpdateResult {
	fmt.Println("Handling password phase...")
	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flow_id,
		Choice: authn.FCPassword,
		Data: authn.FlowUpdateDataPassword{
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return resp.Result
}

func flowHandleProfilePhase(ctx context.Context, client *authn.AuthN, flow_id string, profile *authn.ProfileData) *authn.FlowUpdateResult {
	fmt.Println("Handling profile phase...")
	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flow_id,
		Choice: authn.FCProfile,
		Data: authn.FlowUpdateDataProfile{
			Profile: *profile,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return resp.Result
}

func flowHandleAgreementsPhase(ctx context.Context, client *authn.AuthN, flow_id string, result *authn.FlowUpdateResult) *authn.FlowUpdateResult {
	// Iterate over flow_choices in response.result
	fmt.Println("Handling agreements phase...")
	agreed := []string{}
	for _, flowChoice := range result.FlowChoices {
		// Check if the choice is AGREEMENTS
		if flowChoice.Choice == string(authn.FCAgreements) {
			// Assuming flowChoice.Data["agreements"] is a map[string]interface{}
			agreements, ok := flowChoice.Data["agreements"].(map[string]interface{})
			if ok {
				// Iterate over agreements and append the "id" values to agreed slice
				for _, v := range agreements {
					agreement, ok := v.(map[string]interface{})
					if ok {
						id, ok := agreement["id"].(string)
						if ok {
							agreed = append(agreed, id)
						}
					}
				}
			}
		}
	}

	resp, err := client.Flow.Update(ctx, authn.FlowUpdateRequest{
		FlowID: flow_id,
		Choice: authn.FCAgreements,
		Data: authn.FlowUpdateDataAgreements{
			Agreed: agreed,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return resp.Result
}

func choiceIsAvailable(choices []authn.FlowChoiceItem, choice string) bool {
	for _, fc := range choices {
		if fc.Choice == choice {
			return true

		}
	}
	return false
}

func CreateAndLogin(client *authn.AuthN, email, password string, profile *authn.ProfileData) *authn.FlowCompleteResult {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	fmt.Println("Flow starting...")
	fsresp, err := client.Flow.Start(ctx,
		authn.FlowStartRequest{
			Email:     email,
			FlowTypes: []authn.FlowType{authn.FTsignup, authn.FTsignin},
			CBURI:     CB_URI,
		})

	if err != nil {
		log.Fatal(err)
	}
	flowID := fsresp.Result.FlowID
	var result *authn.FlowUpdateResult = nil
	flowPhase := "initial"
	choices := fsresp.Result.FlowChoices

	for flowPhase != "phase_completed" {
		if choiceIsAvailable(choices, string(authn.FCPassword)) {
			result = flowHandlePasswordPhase(ctx, client, flowID, password)
		} else if choiceIsAvailable(choices, string(authn.FCProfile)) {
			result = flowHandleProfilePhase(ctx, client, flowID, profile)
		} else if choiceIsAvailable(choices, string(authn.FCAgreements)) {
			result = flowHandleAgreementsPhase(ctx, client, flowID, result)
		} else {
			fmt.Printf("Phase %s not handled", result.FlowPhase)
		}
		flowPhase = result.FlowPhase
		choices = result.FlowChoices
	}

	fcresp, err := client.Flow.Complete(ctx,
		authn.FlowCompleteRequest{
			FlowID: flowID,
		})

	if err != nil {
		log.Fatal(err)
	}
	return fcresp.Result
}

func authnUserActions(cmd *cobra.Command) error {
	// Set up variables to be used
	RANDOM_VALUE := strconv.Itoa(rand.Intn(10000000))
	USER_EMAIL := "user.email+goexample" + RANDOM_VALUE + "@pangea.cloud"
	PROFILE_INITIAL := &authn.ProfileData{
		"first_name": "Name",
		"last_name":  "User",
	}
	PROFILE_UPDATE := authn.ProfileData{"first_name": "NameUpdate"}
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
	config, err := pangea.NewConfig(option.WithToken(token), option.WithDomain(os.Getenv("PANGEA_DOMAIN")))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := authn.New(config)

	// Requests examples...
	// User create...
	fmt.Println("Creating user...")
	result := CreateAndLogin(client, USER_EMAIL, PASSWORD_OLD, PROFILE_INITIAL)
	fmt.Println("Create user success. Result: " + pangea.Stringify(result))

	// User password change
	fmt.Println("\n\nUser password change..")
	input3 := authn.ClientPasswordChangeRequest{
		Token:       result.ActiveToken.Token,
		OldPassword: PASSWORD_OLD,
		NewPassword: PASSWORD_NEW,
	}
	_, err = client.Client.Password.Change(ctx, input3)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Println("User password change success")

	// User profile get by email
	fmt.Println("User profile get by email...")
	input4 := authn.UserProfileGetRequest{
		Email: pangea.String(USER_EMAIL),
	}
	resp4, err := client.User.Profile.Get(ctx, input4)

	USER_ID := resp4.Result.ID
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Println("Get profile success. Profile: " + pangea.Stringify(resp4.Result.Profile))

	// User profile get by id
	fmt.Printf("Get profile by ID... %s\n", USER_ID)
	input5 := authn.UserProfileGetRequest{
		ID: pangea.String(USER_ID),
	}
	resp5, err := client.User.Profile.Get(ctx, input5)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
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
		return err
	}
	fmt.Println("Update profile success. Profile: " + pangea.Stringify(resp6.Result.Profile))

	// User update
	fmt.Println("User update...")
	input7 := authn.UserUpdateRequest{
		Email:    pangea.String(USER_EMAIL),
		Disabled: pangea.Bool(false),
	}
	resp7, err := client.User.Update(ctx, input7)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Println("Update user success. Result: " + pangea.Stringify(resp7.Result))

	// User list
	fmt.Println("User list...")
	input8 := authn.UserListRequest{}
	resp8, err := client.User.List(ctx, input8)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Printf("User list success. There is %d users\n", resp8.Result.Count)

	// User delete
	fmt.Println("User delete...")
	input9 := authn.UserDeleteRequest{
		Email: USER_EMAIL,
	}
	_, err = client.User.Delete(ctx, input9)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Println("User delete success")

	// User list
	fmt.Println("User list...")
	input10 := authn.UserListRequest{}
	resp10, err := client.User.List(ctx, input10)
	if err != nil {
		fmt.Println("Something went wrong...")
		fmt.Println(err)
		return err
	}
	fmt.Printf("User list success. There is %d users\n", resp10.Result.Count)
	return nil
}
