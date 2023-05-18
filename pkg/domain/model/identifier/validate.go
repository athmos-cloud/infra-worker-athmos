package identifier

//func IDParentMatchesWithResource(idParent ID, resourceType types.ResourceType) bool {
//	switch reflect.TypeOf(idParent) {
//	case reflect.TypeOf(Empty{}):
//		return resourceType == types.Provider
//	case reflect.TypeOf(Provider{}):
//		return resourceType == types.VPC || resourceType == types.Network
//	case reflect.TypeOf(VPC{}):
//		return resourceType == types.Network
//	case reflect.TypeOf(Network{}):
//		return resourceType == types.Subnetwork || resourceType == types.Firewall
//	case reflect.TypeOf(Subnetwork{}):
//		return resourceType == types.VM
//	default:
//		return false
//	}
//}
