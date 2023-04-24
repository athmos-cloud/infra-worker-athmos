package domain

type Project struct {
	ID        string
	Name      string
	Owner     string
	Providers map[string]Provider
}
