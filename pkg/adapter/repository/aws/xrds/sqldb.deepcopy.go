package xrds

import "k8s.io/apimachinery/pkg/runtime"

func (in *SQLDatabaseParameters) DeepCopyInto(out *SQLDatabaseParameters) {
	out.MachineType = new(string)
	*out.MachineType = *in.MachineType

	out.NetworkRef = new(string)
	*out.NetworkRef = *in.NetworkRef

	out.PasswordNamespace = new(string)
	*out.PasswordNamespace = *in.PasswordNamespace

	out.PasswordRef = new(string)
	*out.PasswordRef = *in.PasswordRef

	out.ProviderRef = new(string)
	*out.ProviderRef = *in.ProviderRef

	out.Region = new(string)
	*out.Region = *in.Region

	out.ResourceName = new(string)
	*out.ResourceName = *in.ResourceName

	out.SqlType = new(string)
	*out.SqlType = *in.SqlType

	out.SqlVersion = new(string)
	*out.SqlVersion = *in.SqlVersion

	out.StorageGB = new(float64)
	*out.StorageGB = *in.StorageGB

	out.SubnetGroupName = new(string)
	*out.SubnetGroupName = *in.SubnetGroupName

	out.Subnet1 = new(string)
	*out.Subnet1 = *in.Subnet1

	out.Subnet1CidrBlock = new(string)
	*out.Subnet1CidrBlock = *in.Subnet1CidrBlock

	out.Subnet2 = new(string)
	*out.Subnet2 = *in.Subnet2

	out.Subnet2CidrBlock = new(string)
	*out.Subnet2CidrBlock = *in.Subnet2CidrBlock
}

func (in *SQLDatabaseParameters) DeepCopy() *SQLDatabaseParameters {
	if in == nil {
		return nil
	}
	out := &SQLDatabaseParameters{}
	in.DeepCopyInto(out)
	return out
}

func (in *SQLDatabaseSpec) DeepCopy() *SQLDatabaseSpec {
	return &SQLDatabaseSpec{
		Parameters: *in.Parameters.DeepCopy(),
	}
}

func (in *SQLDatabase) DeepCopy() *SQLDatabase {
	if in == nil {
		return nil
	}

	out := &SQLDatabase{}
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = *in.Spec.DeepCopy()

	temp := new(SQLDatabaseStatus)
	in.Status.ResourceStatus.DeepCopyInto(&temp.ResourceStatus)
	out.Status = *temp

	return out
}

func (in *SQLDatabase) DeepCopyObject() runtime.Object {
	if out := in.DeepCopy(); out != nil {
		return out
	}
	return nil
}

func (in *SQLDatabaseList) DeepCopy() *SQLDatabaseList {
	if in == nil {
		return nil
	}

	out := &SQLDatabaseList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	var sqlDatabases []SQLDatabase
	for _, sqlDatabase := range in.Items {
		sqlDatabases = append(sqlDatabases, *sqlDatabase.DeepCopy())
	}
	out.Items = sqlDatabases

	return out
}

func (in *SQLDatabaseList) DeepCopyObject() runtime.Object {
	if out := in.DeepCopy(); out != nil {
		return out
	}
	return nil
}
