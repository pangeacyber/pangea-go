package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/authz"
)

const (
	typeFolder     = "folder"
	typeUser       = "user"
	relationOwner  = "owner"
	relationEditor = "editor"
	relationReader = "reader"
)

func main() {
	var timeNow = time.Now()
	var timeStr = timeNow.Format("20060102_150405")
	var folder1 = "folder_1_" + timeStr
	var folder2 = "folder_2_" + timeStr
	var user1 = "user_1_" + timeStr
	var user2 = "user_2_" + timeStr

	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_AUTHZ_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// create a new AuthZ client with pangea token and domain
	client := authz.New(&pangea.Config{
		Token:              token,
		BaseURLTemplate:    os.Getenv("PANGEA_URL_TEMPLATE"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  120 * time.Second,
		Retry:              true,
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	})

	// Create tuples
	fmt.Println("Creating tuples...")
	_, err := client.TupleCreate(ctx, &authz.TupleCreateRequest{
		Tuples: []authz.Tuple{
			{
				Resource: authz.Resource{
					Type: typeFolder,
					ID:   folder1,
				},
				Relation: relationReader,
				Subject: authz.Subject{
					Type: typeUser,
					ID:   user1,
				},
			},
			{
				Resource: authz.Resource{
					Type: typeFolder,
					ID:   folder1,
				},
				Relation: relationEditor,
				Subject: authz.Subject{
					Type: typeUser,
					ID:   user2,
				},
			},
			{
				Resource: authz.Resource{
					Type: typeFolder,
					ID:   folder2,
				},
				Relation: relationEditor,
				Subject: authz.Subject{
					Type: typeUser,
					ID:   user1,
				},
			},
			{
				Resource: authz.Resource{
					Type: typeFolder,
					ID:   folder2,
				},
				Relation: relationOwner,
				Subject: authz.Subject{
					Type: typeUser,
					ID:   user2,
				},
			},
		},
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}
	fmt.Println("Tuples created.")

	// Tuple list with resource
	fmt.Println("Listing tuples with resource...")
	filter := authz.NewFilterUserList()
	filter.ResourceType().Set(pangea.String(typeFolder))
	filter.ResourceID().Set(pangea.String(folder1))

	rListWithResource, err := client.TupleList(ctx, &authz.TupleListRequest{
		Filter: filter.Filter(),
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}

	fmt.Printf("Got %d tuples.\n", rListWithResource.Result.Count)
	for i, tuple := range rListWithResource.Result.Tuples {
		fmt.Printf("Tuple #%d\n", i)
		fmt.Printf("\tType: %s\n", tuple.Subject.Type)
		fmt.Printf("\tID: %s\n", tuple.Subject.ID)
	}

	// Tuple list with subject
	filter = authz.NewFilterUserList()
	fmt.Println("Listing tuples with subject...")
	filter.SubjectType().Set(pangea.String(typeUser))
	filter.SubjectID().Set(pangea.String(user1))

	rListWithSubject, err := client.TupleList(ctx, &authz.TupleListRequest{
		Filter: filter.Filter(),
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}

	fmt.Printf("Got %d tuples.\n", rListWithSubject.Result.Count)
	for i, tuple := range rListWithResource.Result.Tuples {
		fmt.Printf("Tuple #%d\n", i)
		fmt.Printf("\tType: %s\n", tuple.Subject.Type)
		fmt.Printf("\tID: %s\n", tuple.Subject.ID)
	}

	// Tuple delete
	fmt.Println("Deleting tuples...")
	_, err = client.TupleDelete(ctx, &authz.TupleDeleteRequest{
		Tuples: []authz.Tuple{
			{
				Resource: authz.Resource{
					Type: typeFolder,
					ID:   folder1,
				},
				Relation: relationReader,
				Subject: authz.Subject{
					Type: typeUser,
					ID:   user1,
				},
			},
		},
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}
	fmt.Println("Delete success.")

	// Check no debug
	fmt.Println("Checking tuple...")
	rCheck, err := client.Check(ctx, &authz.CheckRequest{
		Resource: authz.Resource{
			Type: typeFolder,
			ID:   folder1,
		},
		Action: "reader",
		Subject: authz.Subject{
			Type: typeUser,
			ID:   user2,
		},
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}

	if rCheck.Result.Allowed {
		fmt.Println("Subject IS allowed to read resource")
	} else {
		fmt.Println("Subject is NOT allowed to read resource")
	}

	// Check debug
	fmt.Println("Checking tuple with debug enabled...")
	rCheck, err = client.Check(ctx, &authz.CheckRequest{
		Resource: authz.Resource{
			Type: typeFolder,
			ID:   folder1,
		},
		Action: "editor",
		Subject: authz.Subject{
			Type: typeUser,
			ID:   user2,
		},
		Debug: pangea.Bool(true),
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}

	if rCheck.Result.Allowed {
		fmt.Println("Subject IS allowed to edit resource")
	} else {
		fmt.Println("Subject is NOT allowed to edit resource")
	}
	if rCheck.Result.Debug != nil {
		fmt.Printf("Debug data: %s\n", pangea.Stringify(*rCheck.Result.Debug))
	}

	// List resources
	fmt.Println("Listing resources...")
	rListResources, err := client.ListResources(ctx, &authz.ListResourcesRequest{
		Type:   typeFolder,
		Action: relationEditor,
		Subject: authz.Subject{
			Type: typeUser,
			ID:   user2,
		},
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}

	fmt.Printf("Got %d resources.\n", len(rListResources.Result.IDs))
	for i, id := range rListResources.Result.IDs {
		fmt.Printf("Resource #%d. ID: %s\n", i, id)
	}

	// List subjects
	fmt.Println("Listing subjects...")
	rListSubjects, err := client.ListSubjects(ctx, &authz.ListSubjectsRequest{
		Resource: authz.Resource{
			Type: typeFolder,
			ID:   folder2,
		},
		Action: relationEditor,
	})
	if err != nil {
		log.Fatal("Unexpected error.", err)
	}

	fmt.Printf("Got %d subjects.\n", len(rListSubjects.Result.Subjects))
	for i, subject := range rListSubjects.Result.Subjects {
		fmt.Printf("Tuple #%d\n", i)
		fmt.Printf("\tType: %s\n", subject.Type)
		fmt.Printf("\tID: %s\n", subject.ID)
	}

}
