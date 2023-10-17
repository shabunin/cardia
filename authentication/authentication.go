package authentication

import "github.com/shabunin/cardia/user"

type Authenticator struct {
	// TODO: jwt keypair
	// TODO: database table
}

func (a *Authenticator) NewTokenForUser(u user.User) string {
	return ""
}

func (a *Authenticator) AuthenticateWithPassword(username string, password string) (string, error) {
	return "", nil
}

func (a *Authenticator) AuthenticateWithPubkey(
	username string,
	pubkeyPayload []byte,
	signCallback func(request []byte) []byte) (string, error) {

	return "", nil
}

type Verifier struct {
	// TODO: jwt pubkey
}

func (v *Verifier) VerifyToken(token string) (user.User, error) {
	return user.User{}, nil
}
