package domain

type Firewall struct {
	ID    string
	Name  string
	Allow RuleList
	Deny  RuleList
}

type Rule struct {
	Protocol string
	Ports    []int
}

type RuleList []Rule
