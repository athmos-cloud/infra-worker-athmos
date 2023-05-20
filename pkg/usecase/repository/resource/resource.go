package resource

type Resource interface {
	Provider
	VPC
	Network
	Subnetwork
	Firewall
	VM
}
