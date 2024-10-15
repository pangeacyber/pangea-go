//go:build unit

package authn_test

import (
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/authn"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	// Create a new FilterUserList
	filterUserList := authn.NewFilterUserList()
	assert.Equal(t, 0, len(filterUserList.Filter()))

	// Set values using setters
	filterUserList.AcceptedEulaID().Set(pangea.String("new_accepted_eula_id"))
	filterUserList.CreatedAt().Set(pangea.String("new_created_at"))
	filterUserList.Disabled().Set(pangea.Bool(true))
	filterUserList.Email().Set(pangea.String("new_email"))
	filterUserList.ID().Set(pangea.String("new_id"))
	filterUserList.LastLoginAt().Set(pangea.String("new_last_login_at"))
	filterUserList.LastLoginIP().Set(pangea.String("new_last_login_ip"))
	filterUserList.LastLoginCity().Set(pangea.String("new_last_login_city"))
	filterUserList.LastLoginCountry().Set(pangea.String("new_last_login_country"))
	filterUserList.LoginCount().Set(pangea.Int(42))
	filterUserList.RequireMFA().Set(pangea.Bool(false))
	l := []string{"scope1", "scope2"}
	filterUserList.Scopes().Set(&l)
	filterUserList.Verified().Set(pangea.Bool(true))

	// Use assert to compare set values with getter values
	assert.Equal(t, "new_accepted_eula_id", *filterUserList.AcceptedEulaID().Get())
	assert.Equal(t, "new_created_at", *filterUserList.CreatedAt().Get())
	assert.Equal(t, true, *filterUserList.Disabled().Get())
	assert.Equal(t, "new_email", *filterUserList.Email().Get())
	assert.Equal(t, "new_id", *filterUserList.ID().Get())
	assert.Equal(t, "new_last_login_at", *filterUserList.LastLoginAt().Get())
	assert.Equal(t, "new_last_login_ip", *filterUserList.LastLoginIP().Get())
	assert.Equal(t, "new_last_login_city", *filterUserList.LastLoginCity().Get())
	assert.Equal(t, "new_last_login_country", *filterUserList.LastLoginCountry().Get())
	assert.Equal(t, 42, *filterUserList.LoginCount().Get())
	assert.Equal(t, false, *filterUserList.RequireMFA().Get())
	assert.Equal(t, []string{"scope1", "scope2"}, *filterUserList.Scopes().Get())
	assert.Equal(t, true, *filterUserList.Verified().Get())

	assert.Equal(t, 13, len(filterUserList.Filter()))

	// remove values
	filterUserList.AcceptedEulaID().Set(nil)
	filterUserList.CreatedAt().Set(nil)
	filterUserList.Disabled().Set(nil)
	filterUserList.Email().Set(nil)
	filterUserList.ID().Set(nil)
	filterUserList.LastLoginAt().Set(nil)
	filterUserList.LastLoginIP().Set(nil)
	filterUserList.LastLoginCity().Set(nil)
	filterUserList.LastLoginCountry().Set(nil)
	filterUserList.LoginCount().Set(nil)
	filterUserList.RequireMFA().Set(nil)
	filterUserList.Scopes().Set(nil)
	filterUserList.Verified().Set(nil)

	assert.Equal(t, 0, len(filterUserList.Filter()))
}
