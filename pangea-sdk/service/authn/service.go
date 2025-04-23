package authn

import (
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Agreements struct {
	*pangea.Client
}

// Group is an API client for AuthN's `/group/` endpoints.
type Group struct {
	*pangea.Client
}

// UserGroup is an API client for AuthN's `/user/group/` endpoints.
type UserGroup struct {
	*pangea.Client
}

type UserProfile struct {
	*pangea.Client
}

type UserInvite struct {
	*pangea.Client
}

type Flow struct {
	*pangea.Client
}

type User struct {
	*pangea.Client
	Profile        *UserProfile
	Invites        *UserInvite
	Authenticators *UserAuthenticators
	Group          *UserGroup
}

type Session struct {
	*pangea.Client
}

type Client struct {
	client   *pangea.Client
	Session  *ClientSession
	Password *ClientPassword
	Token    *ClientToken
}

type ClientSession struct {
	*pangea.Client
}

type ClientPassword struct {
	*pangea.Client
}

type ClientToken struct {
	*pangea.Client
}

type UserAuthenticators struct {
	*pangea.Client
}

type AuthN struct {
	client     *pangea.Client
	User       *User
	Flow       *Flow
	Client     *Client
	Session    *Session
	Agreements *Agreements
	Group      *Group
}

func newAgreements(cli *pangea.Client) *Agreements {
	return &Agreements{
		Client: cli,
	}
}

// newGroup creates a new Group client.
func newGroup(cli *pangea.Client) *Group {
	return &Group{
		Client: cli,
	}
}

func newFlow(cli *pangea.Client) *Flow {
	return &Flow{
		Client: cli,
	}
}

func newClient(cli *pangea.Client) *Client {
	return &Client{
		client:   cli,
		Session:  newClientSession(cli),
		Password: newClientPassword(cli),
		Token:    newClientToken(cli),
	}
}

func newSession(cli *pangea.Client) *Session {
	return &Session{
		Client: cli,
	}
}

func newClientSession(cli *pangea.Client) *ClientSession {
	return &ClientSession{
		Client: cli,
	}
}

func newClientPassword(cli *pangea.Client) *ClientPassword {
	return &ClientPassword{
		Client: cli,
	}
}

func newClientToken(cli *pangea.Client) *ClientToken {
	return &ClientToken{
		Client: cli,
	}
}

func newUserAuthenticators(cli *pangea.Client) *UserAuthenticators {
	return &UserAuthenticators{
		Client: cli,
	}
}

// newUserGroup creates a new UserGroup client.
func newUserGroup(cli *pangea.Client) *UserGroup {
	return &UserGroup{
		Client: cli,
	}
}

func newUserProfile(cli *pangea.Client) *UserProfile {
	return &UserProfile{
		Client: cli,
	}
}

func newUserInvites(cli *pangea.Client) *UserInvite {
	return &UserInvite{
		Client: cli,
	}
}

func newUser(cli *pangea.Client) *User {
	return &User{
		Client:         cli,
		Profile:        newUserProfile(cli),
		Invites:        newUserInvites(cli),
		Authenticators: newUserAuthenticators(cli),
		Group:          newUserGroup(cli),
	}
}

func New(cfg *pangea.Config) *AuthN {
	pc := pangea.NewClient("authn", cfg)
	cli := &AuthN{
		client:     pc,
		User:       newUser(pc),
		Flow:       newFlow(pc),
		Client:     newClient(pc),
		Session:    newSession(pc),
		Agreements: newAgreements(pc),
		Group:      newGroup(pc),
	}
	return cli
}
