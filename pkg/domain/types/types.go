package types

type ResourceType string

const (
	Provider   ResourceType = "provider"
	VPC        ResourceType = "vpc"
	Subnetwork ResourceType = "subnet"
	Network    ResourceType = "network"
	VM         ResourceType = "vm"
	Firewall   ResourceType = "firewall"
)

type ProviderType string

const (
	AWS   ProviderType = "aws"
	Azure ProviderType = "azure"
	GCP   ProviderType = "gcp"
)
