package types

import "strings"

type Resource string

const (
	VPC      Resource = "vpc"
	Net      Resource = "net"
	Subnet   Resource = "subnet"
	Firewall Resource = "firewall"
	VM       Resource = "vm"
)

func (r Resource) FromString(str string) Resource {
	strShaped := strings.ReplaceAll(strings.ToLower(str), " ", "")
	switch strShaped {
	case string(VPC):
		return VPC
	case string(Net):
		return Net
	case string(Subnet):
		return Subnet
	case string(Firewall):
		return Firewall
	case string(VM):
		return VM
	default:
		return ""
	}
}
