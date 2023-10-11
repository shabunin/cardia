package user

type Role int

const (
	Regular Role = iota
	Service
	Superuser
)

type User struct {
	Name  string
	Email string
	Role  Role
}
