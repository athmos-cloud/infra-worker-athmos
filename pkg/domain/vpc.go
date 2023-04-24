package domain

type VPC struct {
	ID          string
	Name        string
	Monitored   bool
	Subnetworks map[string]Subnetwork
}
