package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type SqlDB struct {
	SqlDB    string `json:"sql_db"`
	Provider string `json:"provider"`
	VPC      string `json:"vpc"`
	Network  string `json:"network"`
}

func (id *SqlDB) NameFromLabels(labels map[string]string) errors.Error {
	db, ok := labels[SqlDBNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing db identifier")
	}
	provider, ok := labels[ProviderNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpc := labels[VpcNameKey]
	network, ok := labels[NetworkNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}

	*id = SqlDB{
		SqlDB:    db,
		Provider: provider,
		VPC:      vpc,
		Network:  network,
	}
	return errors.OK
}

func (id *SqlDB) Equals(other ID) bool {
	otherVMID, ok := other.(*SqlDB)
	if !ok {
		return false
	}
	return id.SqlDB == otherVMID.SqlDB &&
		id.Provider == otherVMID.Provider &&
		id.VPC == otherVMID.VPC &&
		id.Network == otherVMID.Network
}

func (id *SqlDB) ToIDLabels() map[string]string {
	return map[string]string{
		SqlDBIdentifierKey:    id.SqlDB,
		ProviderIdentifierKey: id.Provider,
		VpcIdentifierKey:      id.VPC,
		NetworkIdentifierKey:  id.Network,
	}
}

func (id *SqlDB) IDFromLabels(labels map[string]string) errors.Error {
	dbID, ok := labels[SqlDBIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing db identifier")
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID := labels[VpcIdentifierKey]
	networkID, ok := labels[NetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}
	*id = SqlDB{
		SqlDB:    dbID,
		Provider: providerID,
		VPC:      vpcID,
		Network:  networkID,
	}
	return errors.OK
}

func (id *SqlDB) ToNameLabels() map[string]string {
	return map[string]string{
		SqlDBNameKey:    id.SqlDB,
		ProviderNameKey: id.Provider,
		VpcNameKey:      id.VPC,
		NetworkNameKey:  id.Network,
	}
}
