package xrds

import (
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/upbound/provider-aws/apis/ec2/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VMInstanceParameters struct {
	AssignPublicIp   *bool                               `json:"assignPublicIp"`
	DeletionPolicy   v1.DeletionPolicy                   `json:"deletionPolicy"`
	Disks            []v1beta1.RootBlockDeviceParameters `json:"disks"`
	KeyPairRef       *string                             `json:"keyPairId"`
	MachineType      *string                             `json:"machineType"`
	NetworkRef       *string                             `json:"networkRef"`
	Os               *string                             `json:"os"`
	ProviderRef      *string                             `json:"providerRef"`
	Region           *string                             `json:"region"`
	SecurityGroupRef *string                             `json:"securityGroupRef"`
	SubnetworkRef    *string                             `json:"subnetworkRef"`
	VmId             *string                             `json:"vmId"`
}

type VMInstanceSpec struct {
	Parameters VMInstanceParameters `json:"parameters"`
}

type VMInstanceStatus struct {
	v1.ResourceStatus `json:",inline"`
	PublicIp          *string `json:"publicIp,omitempty"`
}

type VMInstance struct {
	metav1.TypeMeta   `json:",inline"'`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VMInstanceSpec   `json:"spec"`
	Status            VMInstanceStatus `json:"status,omitempty"`
}

type VMInstanceList struct {
	metav1.TypeMeta `json:",inline"'`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMInstance `json:"items"`
}
