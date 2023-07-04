package xrds

import "k8s.io/apimachinery/pkg/runtime"

func (in *VMInstanceParameters) DeepCopyInto(out *VMInstanceParameters) {
	out.AssignPublicIp = new(bool)
	*out.AssignPublicIp = *in.AssignPublicIp

	out.DeletionPolicy = new(string)
	*out.DeletionPolicy = *in.DeletionPolicy

	out.Disks = *in.Disks.DeepCopy()

	out.KeyPairId = new(string)
	*out.KeyPairId = *in.KeyPairId

	out.MachineType = new(string)
	*out.MachineType = *in.MachineType

	out.NetworkRef = new(string)
	*out.NetworkRef = *in.NetworkRef

	out.Os = new(string)
	*out.Os = *in.Os

	out.ProviderRef = new(string)
	*out.ProviderRef = *in.ProviderRef

	out.Region = new(string)
	*out.Region = *in.Region

	out.SecurityGroupRef = new(string)
	*out.SecurityGroupRef = *in.SecurityGroupRef

	out.SubnetworkRef = new(string)
	*out.SubnetworkRef = *in.SubnetworkRef

	out.VmId = new(string)
	*out.VmId = *in.VmId
}

func (in *VMInstanceParameters) DeepCopy() *VMInstanceParameters {
	if in == nil {
		return nil
	}
	out := &VMInstanceParameters{}
	in.DeepCopyInto(out)
	return out
}

func (in *VMInstanceSpec) DeepCopy() *VMInstanceSpec {
	return &VMInstanceSpec{
		Parameters: *in.Parameters.DeepCopy(),
	}
}

func (in *VMInstance) DeepCopy() *VMInstance {
	if in == nil {
		return nil
	}

	out := &VMInstance{}
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = *in.Spec.DeepCopy()

	temp := new(VMInstanceStatus)
	in.Status.ResourceStatus.DeepCopyInto(&temp.ResourceStatus)
	out.Status = *temp

	return out
}

func (in *VMInstanceList) DeepCopy() *VMInstanceList {
	if in == nil {
		return nil
	}
	out := &VMInstanceList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	var newItems []VMInstance
	for _, item := range in.Items {
		newItems = append(newItems, *item.DeepCopy())
	}
	out.Items = newItems
	return out
}

func (in *VMInstance) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func (in *VMInstanceList) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}
