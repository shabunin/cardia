package authentication

import (
	"github.com/pocketbase/dbx"
)

type Role int

const (
	Regular Role = iota
	Service
	Superuser
)

type User struct {
	Name  string
	Role  Role
	Home  string
	Email string
}

type user struct {
	enabled  bool
	username string
	password string
	role     string
	email    string
	home     string
}

const (
	roleRegular   string = "u"
	roleService   string = "s"
	roleSuperuser string = "r"
)

func (u user) Export() User {
	var r Role
	switch u.role {
	case roleSuperuser:
		r = Superuser
	case roleService:
		r = Service
	case roleRegular:
		r = Regular
	default:
		r = Regular
	}
	return User{
		Name:  u.username,
		Role:  r,
		Home:  u.home,
		Email: u.email,
	}
}

const (
	tableUsers        = "users"
	fieldUserEnabled  = "enabled"
	fieldUserUsername = "username"
	fieldUserPassword = "password"
	fieldUserRole     = "role"
	fieldUserEmail    = "email"
	fieldUserHome     = "home"
	indexUserUsername = "username_idx"
	indexUserEmail    = "email_idx"
	indexUserHome     = "home_idx"
)

func initUsersTable(db *dbx.DB) error {
	users := make(map[string]string)
	users[fieldUserEnabled] = "BOOLEAN DEFAULT TRUE NOT NULL"
	users[fieldUserUsername] = "TEXT PRIMARY KEY NOT NULL"
	users[fieldUserPassword] = "TEXT NOT NULL"
	users[fieldUserRole] = "TEXT DEFAULT 'u' NOT NULL"
	users[fieldUserEmail] = "TEXT DEFAULT '' NOT NULL"
	users[fieldUserHome] = "TEXT PRIMARY KEY NOT NULL"

	query := db.CreateTable(tableUsers, users)
	_, err := query.Execute()
	if err != nil {
		return err
	}

	query = db.CreateUniqueIndex(tableUsers, indexUserUsername, fieldUserUsername)
	_, err = query.Execute()
	if err != nil {
		return err
	}

	query = db.CreateUniqueIndex(tableUsers, indexUserEmail, fieldUserEmail)
	_, err = query.Execute()
	if err != nil {
		return err
	}

	query = db.CreateUniqueIndex(tableUsers, indexUserHome, fieldUserHome)
	_, err = query.Execute()
	if err != nil {
		return err
	}

	return err
}

func selectUser(db *dbx.DB, username string) (user, error) {
	var u user
	e := db.Select(
		fieldUserEnabled,
		fieldUserUsername,
		fieldUserPassword,
		fieldUserEmail,
		fieldUserHome).
		From(tableUsers).
		Where(dbx.HashExp{
			fieldUserUsername: username,
		}).
		One(&u)
	return u, e
}

func createUser(db *dbx.DB, u user) error {
	_, e := db.Insert(tableUsers,
		dbx.Params{
			fieldUserEnabled:  u.enabled,
			fieldUserUsername: u.username,
			fieldUserPassword: u.password,
			fieldUserEmail:    u.email,
			fieldUserHome:     u.home,
		}).Execute()
	return e
}

func updateUser(db *dbx.DB, username string, values dbx.Params) error {
	_, e := db.Update(tableUsers, values,
		dbx.HashExp{
			fieldUserUsername: username,
		}).Execute()
	return e
}

func updatePassword(db *dbx.DB, username string, phash string) error {
	return updateUser(db, username, dbx.Params{fieldUserPassword: phash})
}

func enableUser(db *dbx.DB, username string) error {
	return updateUser(db, username, dbx.Params{fieldUserEnabled: true})
}
func disableUser(db *dbx.DB, username string) error {
	return updateUser(db, username, dbx.Params{fieldUserEnabled: false})
}

func deleteUser(db *dbx.DB, username string) error {
	_, e := db.Delete(tableUsers,
		dbx.HashExp{
			fieldUserUsername: username,
		}).Execute()
	return e
}

func listUsers(db *dbx.DB, q *dbx.AndOrExp, offset int64, limit int64, sortBy string) (int, []user, error) {

	var e error

	var c int
	e = db.Select("COUNT (*)").From(tableUsers).One(&c)

	var u []user
	e = db.Select(
		fieldUserEnabled,
		fieldUserUsername,
		fieldUserPassword,
		fieldUserEmail,
		fieldUserHome).
		From(tableUsers).
		Where(q).
		Limit(limit).
		Offset(offset).
		OrderBy(sortBy).
		All(&u)

	if e != nil {
		return 0, nil, e
	}

	return c, u, nil
}
