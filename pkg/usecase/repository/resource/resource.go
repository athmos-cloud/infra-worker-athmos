package resource

type Resource interface {
	Provider
	VPC
	Network
	Subnetwork
	Firewall
	VM
}

type FindResourceOption struct {
	Name      string
	Namespace string
}

type FindAllResourceOption struct {
	Labels    map[string]string
	Namespace string
}

type ResourceExistsOption struct {
	Labels    map[string]string
	Namespace string
}
