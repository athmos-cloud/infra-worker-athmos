package resource

type Resource interface {
	Provider
	Network
	Subnetwork
	Firewall
	VM
	SqlDB
}

type FindResourceOption struct {
	Name string
}

type FindAllResourceOption struct {
	Labels map[string]string
}
