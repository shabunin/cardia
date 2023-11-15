package authentication

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/dbx"
	"github.com/shabunin/cardia/database"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"path/filepath"
	"time"
)

type identityClaims struct {
	jwt.RegisteredClaims
	User string `json:"user"`
	Role string `json:"role"`
}

const phashCost = 14

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

	_ = initUsersTable(db)

	return &Authenticator{db: db}, nil
}

func (a *Authenticator) newTokenForUser(u User) (string, error) {
	id := identityClaims{User: u.Name}
	switch u.Role {
	case Regular:
		id.Role = roleRegular
	case Service:
		id.Role = roleService
	case Superuser:
		id.Role = roleSuperuser
	}
	// TODO : customize claims
	id.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, id)
	ss, err := token.SignedString(a.jwtSigner)
	return ss, err
}

func (a *Authenticator) AuthenticateWithPassword(username string, password string) (string, error) {
	phash, err := generateFromPassword([]byte(password), phashCost)
	if err != nil {
		return "", err
	}
	u, err := selectUserPass(a.db, username, string(phash))
	if err != nil {
		return "", err
	}

	err = compareHashAndPassword([]byte(u.password), []byte(password))
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

func (v *Verifier) VerifyToken(token string) (User, error) {
	parsed, err := jwt.ParseWithClaims(token, &identityClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("signing method not supported")
			}
			return v.jwtVerifier, nil
		})
	if claims, ok := parsed.Claims.(*identityClaims); ok && parsed.Valid {
		fmt.Printf("%v %v", claims.User, claims.RegisteredClaims.Issuer)
		// TODO create User struct
		return User{}, nil
	} else {
		return User{}, fmt.Errorf("cannot parse token: %w", err)
	}
}
