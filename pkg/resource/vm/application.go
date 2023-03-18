package vm

import "github.com/PaulBarrie/infra-worker/pkg/resource"

type Application struct {
	ResourceReference resource.Reference
	VPC               string   `bson:"vpc"`
	Network           string   `bson:"network"`
	Subnetwork        string   `bson:"subnetwork"`
	Zone              string   `bson:"zone"`
	MachineType       string   `bson:"machineType"`
	Auths             []VMAuth `bson:"auths"`
	Disks             []Disk   `bson:"disks"`
	OS                OS       `bson:"os"`
}
