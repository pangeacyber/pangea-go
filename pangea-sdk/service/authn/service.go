package authn

import (
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// TODO: Do we need this interface?
// type Client interface {
// 	Lookup(ctx context.Context, input *DomainLookupInput) (*pangea.PangeaResponse[DomainLookupOutput], error)
// }

type Password struct {
	*pangea.Client
}

type Profile struct {
	*pangea.Client
}

type Invites struct {
	*pangea.Client
}

type User struct {
	*pangea.Client
	Profile *Profile
	Invites *Invites
}

type AuthN struct {
	*pangea.Client
	Password *Password
	User     *User
}

func newPassword(cli *pangea.Client) *Password {
	return &Password{
		Client: cli,
	}
}

func newProfile(cli *pangea.Client) *Profile {
	return &Profile{
		Client: cli,
	}
}

func newInvites(cli *pangea.Client) *Invites {
	return &Invites{
		Client: cli,
	}
}

func newUser(cli *pangea.Client) *User {
	return &User{
		Client:  cli,
		Profile: newProfile(cli),
		Invites: newInvites(cli),
	}
}

func New(cfg *pangea.Config) *AuthN {
	pc := pangea.NewClient("authn", cfg)
	cli := &AuthN{
		Client:   pc,
		Password: newPassword(pc),
		User:     newUser(pc),
	}
	return cli
}
