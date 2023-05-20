package resource

type ProviderResource interface {
	Provider
	VPC
	Network
	Subnetwork
	Firewall
	VM
}
