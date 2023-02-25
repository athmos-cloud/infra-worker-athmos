package auth

type Auth struct{}

type IAuth interface {
	Generate() error
}

type Owner struct {
	User User
}
