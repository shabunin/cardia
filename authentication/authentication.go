package authentication

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/dbx"
	"github.com/shabunin/cardia/database"
	"github.com/shabunin/cardia/user"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"path/filepath"
)

type identity struct {
	jwt.RegisteredClaims
	User string `json:"username"`
	Role string `json:"role"`
}

const (
	roleRegular   string = "regular"
	roleService   string = "service"
	roleSuperuser string = "superuser"
)

const (
	minCost     = bcrypt.MinCost
	defaultCost = 12 // ballpark: 250 msec on a modern Intel CPU
)

func generateFromPassword(password []byte, cost int) (result []byte, err error) {
	sum := sha3.Sum512(password)
	return bcrypt.GenerateFromPassword(sum[:], cost)
}
func compareHashAndPassword(hashedPassword, password []byte) error {
	sum := sha3.Sum512(password)
	return bcrypt.CompareHashAndPassword(hashedPassword, sum[:])
}

type Authenticator struct {
	db        *dbx.DB
	jwtSigner *rsa.PrivateKey
}

func NewAuthenticator(dbpath string) (*Authenticator, error) {
	if !filepath.IsAbs(dbpath) {
		base, _ := os.Getwd()
		dbpath = path.Join(base, dbpath)
	}
	db, err := database.ConnectDB(dbpath)
	if err != nil {
		return nil, err
	}

	return &Authenticator{db: db}, nil
}

func (a *Authenticator) newTokenForUser(u user.User) string {
	id := identity{User: u.Name}
	switch u.Role {
	case user.Regular:
		id.Role = roleRegular
	case user.Service:
		id.Role = roleService
	case user.Superuser:
		id.Role = roleSuperuser
	}
	// TODO sign
	return ""
}

func (a *Authenticator) AuthenticateWithPassword(username string, password string) (string, error) {
	var err error
	// TODO find user in db

	var fromDb []byte
	err = compareHashAndPassword(fromDb, []byte(password))
	if err != nil {
		return "", errors.New("wrong credentials")
	}

	// TODO return newTokenForUser

	return "", nil
}

func (a *Authenticator) AuthenticateWithPubkey(
	username string,
	pubkeyPayload []byte,
	signCallback func(request []byte) []byte) (string, error) {

	return "", nil
}

type Verifier struct {
	jwtVerifier *rsa.PublicKey
}

func (v *Verifier) VerifyToken(token string) (user.User, error) {
	parsed, err := jwt.ParseWithClaims(token, &identity{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("signing method not supported")
			}
			return v.jwtVerifier, nil
		})
	if claims, ok := parsed.Claims.(*identity); ok && parsed.Valid {
		fmt.Printf("%v %v", claims.User, claims.RegisteredClaims.Issuer)
		// TODO create User struct
		return user.User{}, nil
	} else {
		return user.User{}, fmt.Errorf("cannot parse token: %w", err)
	}
}
