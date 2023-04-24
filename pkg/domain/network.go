package domain

type Network struct {
	Name        string
	Subnetworks map[string]Subnetwork
}
