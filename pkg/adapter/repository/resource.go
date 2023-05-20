package repository

type FindResourceOption struct {
	Name      string
	Namespace string
}

type FindAllResourceOption struct {
	Labels    map[string]string
	Namespace string
}
