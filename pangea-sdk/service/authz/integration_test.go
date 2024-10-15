//go:build integration

package authz_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/authz"
	"github.com/stretchr/testify/assert"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("authz", pangeatesting.Live)

var timeNow = time.Now()
var timeStr = timeNow.Format("20060102_150405")
var folder1 = "folder_1_" + timeStr
var folder2 = "folder_2_" + timeStr
var user1 = "user_1_" + timeStr
var user2 = "user_2_" + timeStr

const (
	typeFolder     = "folder"
	typeUser       = "user"
	relationOwner  = "owner"
	relationEditor = "editor"
	relationReader = "reader"
)

func Test_Integration(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cli := authz.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	// Create tuples
	rCreate, err := cli.TupleCreate(ctx, &authz.TupleCreateRequest{
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

	assert.NoError(t, err)
	assert.NotNil(t, rCreate)
	assert.NotNil(t, rCreate.Result)

	// Tuple list with resource
	filter := authz.NewFilterUserList()
	filter.ResourceType().Set(pangea.String(typeFolder))
	filter.ResourceID().Set(pangea.String(folder1))

	rListWithResource, err := cli.TupleList(ctx, &authz.TupleListRequest{
		Filter: filter.Filter(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, rListWithResource)
	assert.NotNil(t, rListWithResource.Result)
	assert.Equal(t, len(rListWithResource.Result.Tuples), 2)
	assert.Equal(t, rListWithResource.Result.Count, 2)

	// Tuple list with subject
	filter = authz.NewFilterUserList()
	filter.SubjectType().Set(pangea.String(typeUser))
	filter.SubjectID().Set(pangea.String(user1))

	rListWithSubject, err := cli.TupleList(ctx, &authz.TupleListRequest{
		Filter: filter.Filter(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, rListWithSubject)
	assert.NotNil(t, rListWithSubject.Result)
	assert.Equal(t, len(rListWithSubject.Result.Tuples), 2)
	assert.Equal(t, rListWithSubject.Result.Count, 2)

	// Tuple delete
	rDelete, err := cli.TupleDelete(ctx, &authz.TupleDeleteRequest{
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

	assert.NoError(t, err)
	assert.NotNil(t, rDelete)
	assert.NotNil(t, rDelete.Result)
	assert.Empty(t, rDelete.Result)

	// Check no debug
	rCheck, err := cli.Check(ctx, &authz.CheckRequest{
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

	assert.NoError(t, err)
	assert.NotNil(t, rCheck)
	assert.NotNil(t, rCheck.Result)
	assert.False(t, rCheck.Result.Allowed)
	assert.Nil(t, rCheck.Result.Debug)
	assert.NotEmpty(t, rCheck.Result.SchemaID)
	assert.NotEmpty(t, rCheck.Result.SchemaVersion)

	// Check debug
	rCheck, err = cli.Check(ctx, &authz.CheckRequest{
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

	assert.NoError(t, err)
	assert.NotNil(t, rCheck)
	assert.NotNil(t, rCheck.Result)
	assert.True(t, rCheck.Result.Allowed)
	assert.NotNil(t, rCheck.Result.Debug)
	assert.NotEmpty(t, *rCheck.Result.Debug)
	assert.NotEmpty(t, rCheck.Result.SchemaID)
	assert.NotEmpty(t, rCheck.Result.SchemaVersion)

	// List resources
	rListResources, err := cli.ListResources(ctx, &authz.ListResourcesRequest{
		Type:   typeFolder,
		Action: relationEditor,
		Subject: authz.Subject{
			Type: typeUser,
			ID:   user2,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, rListResources)
	assert.NotNil(t, rListResources.Result)
	assert.Equal(t, len(rListResources.Result.IDs), 1)

	// List subjects
	rListSubjects, err := cli.ListSubjects(ctx, &authz.ListSubjectsRequest{
		Resource: authz.Resource{
			Type: typeFolder,
			ID:   folder2,
		},
		Action: relationEditor,
	})

	assert.NoError(t, err)
	assert.NotNil(t, rListSubjects)
	assert.NotNil(t, rListSubjects.Result)
	assert.Equal(t, len(rListSubjects.Result.Subjects), 1)

}
