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

type UserProfile struct {
	*pangea.Client
}

type UserInvite struct {
	*pangea.Client
}

type Flow struct {
	*pangea.Client
	Signup *FlowSignup
	Verify *FlowVerify
	Enroll *FlowEnroll
	Reset  *FlowReset
}

type FlowSignup struct {
	*pangea.Client
}

type FlowReset struct {
	*pangea.Client
}

type FlowVerify struct {
	*pangea.Client
	MFA *FlowVerifyMFA
}

type FlowVerifyMFA struct {
	*pangea.Client
}

type FlowEnroll struct {
	*pangea.Client
	MFA *FlowEnrollMFA
}

type FlowEnrollMFA struct {
	*pangea.Client
}

type User struct {
	*pangea.Client
	Profile *UserProfile
	Invites *UserInvite
	MFA     *UserMFA
	Login   *UserLogin
}

type Session struct {
	*pangea.Client
}

type Client struct {
	*pangea.Client
	Session *ClientSession
}

type ClientSession struct {
	*pangea.Client
}

type UserMFA struct {
	*pangea.Client
}

type UserLogin struct {
	*pangea.Client
}

type AuthN struct {
	*pangea.Client
	Password *Password
	User     *User
	Flow     *Flow
}

func newPassword(cli *pangea.Client) *Password {
	return &Password{
		Client: cli,
	}
}

func newFlowEnrollMFA(cli *pangea.Client) *FlowEnrollMFA {
	return &FlowEnrollMFA{
		Client: cli,
	}
}

func newFlowEnroll(cli *pangea.Client) *FlowEnroll {
	return &FlowEnroll{
		Client: cli,
	}
}

func newFlowReset(cli *pangea.Client) *FlowReset {
	return &FlowReset{
		Client: cli,
	}
}

func newFlowVerifyMFA(cli *pangea.Client) *FlowVerifyMFA {
	return &FlowVerifyMFA{
		Client: cli,
	}
}

func newFlowVerify(cli *pangea.Client) *FlowVerify {
	return &FlowVerify{
		Client: cli,
		MFA:    newFlowVerifyMFA(cli),
	}
}

func newFlowSignup(cli *pangea.Client) *FlowSignup {
	return &FlowSignup{
		Client: cli,
	}
}

func newFlow(cli *pangea.Client) *Flow {
	return &Flow{
		Client: cli,
		Enroll: newFlowEnroll(cli),
		Verify: newFlowVerify(cli),
		Signup: newFlowSignup(cli),
		Reset:  newFlowReset(cli),
	}
}

func newClient(cli *pangea.Client) *Client {
	return &Client{
		Client:  cli,
		Session: newClientSession(cli),
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

func newUserMFA(cli *pangea.Client) *UserMFA {
	return &UserMFA{
		Client: cli,
	}
}

func newUserLogin(cli *pangea.Client) *UserLogin {
	return &UserLogin{
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
		Client:  cli,
		Profile: newUserProfile(cli),
		Invites: newUserInvites(cli),
		MFA:     newUserMFA(cli),
		Login:   newUserLogin(cli),
	}
}

func New(cfg *pangea.Config) *AuthN {
	pc := pangea.NewClient("authn", cfg)
	cli := &AuthN{
		Client:   pc,
		Password: newPassword(pc),
		User:     newUser(pc),
		Flow:     newFlow(pc),
	}
	return cli
}
