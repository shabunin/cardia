package user

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
