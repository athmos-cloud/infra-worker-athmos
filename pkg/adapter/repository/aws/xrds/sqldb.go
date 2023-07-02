package xrds

import (
	"github.com/upbound/provider-aws/apis/rds/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SQLDatabaseParameters struct {
	MachineType       *string  `json:"machineType"`
	NetworkRef        *string  `json:"networkRef"`
	PasswordNamespace *string  `json:"passwordNamespace"`
	PasswordRef       *string  `json:"passwordRef"`
	ProviderRef       *string  `json:"providerRef"`
	Region            *string  `json:"region"`
	ResourceName      *string  `json:"resourceName"`
	SqlType           *string  `json:"sqlType"`
	SqlVersion        *string  `json:"sqlVersion"`
	StorageGB         *float64 `json:"storageGB"`
	StorageGBLimit    *float64 `json:"storageGBLimit"`
	Subnet1           *string  `json:"subnet1"`
	Subnet1IpRange    *string  `json:"subnet1_ip_range"`
	Subnet2           *string  `json:"subnet2"`
	Subnet2IpRange    *string  `json:"subnet2_ip_range"`
	SubnetGroupName   *string  `json:"subnetGroupName"`
}

type SQLDatabaseSpec struct {
	Parameters SQLDatabaseParameters `json:"parameters"`
}

type SQLDatabase struct {
	metav1.TypeMeta   `json:",inline"'`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SQLDatabaseSpec   `json:"spec"`
	Status            SQLDatabaseStatus `json:"status,omitempty"`
}

type SQLDatabaseStatus struct {
	DatabaseStatus v1beta1.InstanceStatus `json:"databaseStatus"`
}

type SQLDatabaseList struct {
	metav1.TypeMeta `json:",inline"'`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQLDatabase `json:"items"`
}
