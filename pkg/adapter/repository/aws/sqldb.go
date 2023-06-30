package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws/xrds"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (aws *awsRepository) FindSqlDB(ctx context.Context, opt option.Option) (*instance.SqlDB, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsRDSInstance := &xrds.SQLDatabase{}
	namespace := ctx.Value(context.CurrentNamespace).(string)
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: namespace}, awsRDSInstance); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("SQL database %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get SQL database %s", req.Name))
	}

	passwordSecretName := fmt.Sprintf("%s-password", req.Name)
	passwordSecret, err := aws.getPasswordSecret(ctx, &passwordSecretName, &namespace)
	if !err.IsOk() {
		return nil, err
	}

	mod, err := aws.toSqlDBModel(awsRDSInstance, passwordSecret)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindAllSqlDBs(ctx context.Context, opt option.Option) (*instance.SqlDBCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsRDSInstances := &xrds.SQLDatabaseList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
		Namespace:     ctx.Value(context.CurrentNamespace).(string),
	}
	if err := kubernetes.Client().Client.List(ctx, awsRDSInstances, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to list dbs"))
	}
	modDBs, err := aws.toSqlDBModelCollection(ctx, awsRDSInstances)
	if !err.IsOk() {
		return nil, err
	}

	return modDBs, errors.OK
}

func (aws *awsRepository) FindAllRecursiveSqlDBs(ctx context.Context, opt option.Option, ch *resourceRepo.SqlDBChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsRDSInstances := &xrds.SQLDatabaseList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsRDSInstances, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage("Unable to list dbs")
		return
	}
	if dbs, err := aws.toSqlDBModelCollection(ctx, awsRDSInstances); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- dbs
	}
}

func (aws *awsRepository) CreateSqlDB(ctx context.Context, db *instance.SqlDB) errors.Error {
	if exists, err := aws.SqlDBExists(ctx, db); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("db %s already exists in network %s", db.IdentifierID.SqlDB, db.IdentifierID.Network))
	}

	passwordSecret, err := aws.toPasswordSecret(ctx, db)
	if !err.IsOk() {
		return err
	}
	if kubeErr := kubernetes.Client().Client.Create(ctx, passwordSecret); kubeErr != nil {
		if k8serrors.IsAlreadyExists(kubeErr) {
			return errors.Conflict.WithMessage(fmt.Sprintf("Password secret for %s already exists", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to create password secret for %s", db.IdentifierID.SqlDB))
	}

	rdsInstance, err := aws.toAWSRDSInstance(ctx, db)
	if !err.IsOk() {
		return err
	}
	if kubeErr := kubernetes.Client().Client.Create(ctx, rdsInstance); kubeErr != nil {
		if k8serrors.IsAlreadyExists(kubeErr) {
			return errors.Conflict.WithMessage(fmt.Sprintf("Db %s already exists", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to create db %s", db.IdentifierID.SqlDB))
	}
	return errors.Created
}

func (aws *awsRepository) UpdateSqlDB(ctx context.Context, db *instance.SqlDB) errors.Error {
	namespace := ctx.Value(context.CurrentNamespace).(string)
	existingDB := &xrds.SQLDatabase{}
	if kubeErr := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB, Namespace: namespace}, existingDB); kubeErr != nil {
		if k8serrors.IsNotFound(kubeErr) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Db %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to get db %s", db.IdentifierID.SqlDB))
	}

	passwordSecretName := fmt.Sprintf("%s-password", db.IdentifierID.SqlDB)
	existingPasswordSecret, err := aws.getPasswordSecret(ctx, &passwordSecretName, &namespace)
	if !err.IsOk() {
		return err
	}
	newPasswordSecret, err := aws.toPasswordSecret(ctx, db)
	if !err.IsOk() {
		return err
	}

	existingPasswordSecret.Data = newPasswordSecret.Data
	existingPasswordSecret.Labels = newPasswordSecret.Labels
	if kubeErr := kubernetes.Client().Client.Update(ctx, existingPasswordSecret); kubeErr != nil {
		if k8serrors.IsNotFound(kubeErr) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Password secret for %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to update password secret for %s", db.IdentifierID.SqlDB))
	}

	newRDSInstance, err := aws.toAWSRDSInstance(ctx, db)
	if !err.IsOk() {
		return err
	}

	existingDB.Name = newRDSInstance.Name
	existingDB.Namespace = newRDSInstance.Namespace
	existingDB.Spec = newRDSInstance.Spec
	existingDB.Labels = newRDSInstance.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Db %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update Db %s", db.IdentifierID.SqlDB))
	}
	return errors.NoContent
}

func (aws *awsRepository) DeleteSqlDB(ctx context.Context, db *instance.SqlDB) errors.Error {
	namespace := ctx.Value(context.CurrentNamespace).(string)
	existingDB := &xrds.SQLDatabase{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB, Namespace: namespace}, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Db %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to get db %s", db.IdentifierID.SqlDB))
	}

	passwordSecretName := fmt.Sprintf("%s-password", db.IdentifierID.SqlDB)
	existingPasswordSecret, err := aws.getPasswordSecret(ctx, &passwordSecretName, &namespace)
	if !err.IsOk() {
		return err
	}

	if err := kubernetes.Client().Client.Delete(ctx, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("db %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete db %s", db.IdentifierID.SqlDB))
	}

	if err := kubernetes.Client().Client.Delete(ctx, existingPasswordSecret); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("Password secret for %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to delete password secret for %s", db.IdentifierID.SqlDB))
	}
	return errors.NoContent
}

func (aws *awsRepository) SqlDBExists(ctx context.Context, db *instance.SqlDB) (bool, errors.Error) {
	awsDBs := &xrds.SQLDatabaseList{}
	parentID := identifier.Network{
		Provider: db.IdentifierID.Provider,
		VPC:      db.IdentifierID.VPC,
		Network:  db.IdentifierID.Network,
	}
	searchLabels := lo.Assign(parentID.ToIDLabels(), map[string]string{identifier.SqlDBNameKey: db.IdentifierName.SqlDB})
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsDBs, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list db")
	}
	return len(awsDBs.Items) > 0, errors.OK
}

func (aws *awsRepository) toSqlDBModel(db *xrds.SQLDatabase, secret *v1.Secret) (*instance.SqlDB, errors.Error) {
	sqlAuthDB, err := aws.toSqlDBAuth(secret)
	if !err.IsOk() {
		return nil, err
	}

	id := identifier.SqlDB{}
	name := identifier.SqlDB{}
	if err := id.IDFromLabels(db.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(db.Labels); !err.IsOk() {
		return nil, err
	}
	version, err := aws.getSqlTypeVersion(db)
	if !err.IsOk() {
		return nil, err
	}
	sizeGib := int(*db.Spec.Parameters.StorageGB)
	maxSizeGib := int(*db.Spec.Parameters.StorageGBLimit)

	return &instance.SqlDB{
		Auth: *sqlAuthDB,
		Disk: instance.SqlDbDisk{
			Type:               instance.DiskTypeSSD,
			SizeGib:            sizeGib,
			AutoresizeLimitGib: &maxSizeGib,
			Autoresize:         0 < maxSizeGib,
		},
		MachineType:    *db.Spec.Parameters.MachineType,
		IdentifierID:   id,
		IdentifierName: name,
		Region:         *db.Spec.Parameters.Region,
		SQLTypeVersion: *version,
	}, errors.OK
}

func (aws *awsRepository) toSqlDBModelCollection(ctx context.Context, dbs *xrds.SQLDatabaseList) (*instance.SqlDBCollection, errors.Error) {
	namespace := ctx.Value(context.CurrentNamespace).(string)
	sqlDBModels := instance.SqlDBCollection{}
	for _, item := range dbs.Items {
		passwordSecretName := fmt.Sprintf("%s-password", item.Name)
		passwordSecret, err := aws.getPasswordSecret(ctx, &passwordSecretName, &namespace)
		if !err.IsOk() {
			return nil, err
		}
		sqlDBModel, err := aws.toSqlDBModel(&item, passwordSecret)
		if !err.IsOk() {
			return nil, err
		}
		sqlDBModels[sqlDBModel.IdentifierName.SqlDB] = *sqlDBModel
	}
	return &sqlDBModels, errors.OK
}

func (aws *awsRepository) toSqlDBAuth(secret *v1.Secret) (*instance.SqlDBAuth, errors.Error) {
	strPassword := fmt.Sprintf("%s", secret.Data["Password"])
	if strPassword == "" {
		return nil, errors.InternalError.WithMessage(fmt.Sprintf("Could not extract password from secret %s", secret.Name))
	}
	return &instance.SqlDBAuth{
		RootPassword: strPassword,
	}, errors.OK
}

func (aws *awsRepository) getSqlTypeVersion(db *xrds.SQLDatabase) (*instance.SQLTypeVersion, errors.Error) {
	sqlVersion := *db.Spec.Parameters.SqlVersion
	var sqlType instance.SQLType
	switch *db.Spec.Parameters.SqlType {
	case "mysql":
		sqlType = instance.MySqlSQLType
		break
	case "postgres":
		sqlType = instance.PostgresSQLType
		break
	case "sqlserver-se":
		sqlType = instance.SQLServerType
		break
	}

	if sqlType == "" {
		return nil, errors.InternalError.WithMessage(fmt.Sprintf("Could not compute database type and version for %s", db.Name))
	}

	return &instance.SQLTypeVersion{
		Type:    sqlType,
		Version: sqlVersion,
	}, errors.OK
}

func (aws *awsRepository) toAWSRDSInstance(ctx context.Context, db *instance.SqlDB) (*xrds.SQLDatabase, errors.Error) {
	labels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), db.IdentifierID.ToIDLabels(), db.IdentifierName.ToNameLabels())
	namespace := ctx.Value(context.CurrentNamespace).(string)
	passwordRef := fmt.Sprintf("%s-password", db.IdentifierID.SqlDB)
	storageSize := float64(db.Disk.SizeGib)
	var maxStorageSize float64
	if db.Disk.Autoresize {
		maxStorageSize = float64(*db.Disk.AutoresizeLimitGib)
	}
	subnetGroupName := fmt.Sprintf("%s-subnet-group", db.IdentifierID.SqlDB)
	subnet1 := fmt.Sprintf("%s-subnet1", db.IdentifierID.SqlDB)
	subnet2 := fmt.Sprintf("%s-subnet2", db.IdentifierID.SqlDB)

	return &xrds.SQLDatabase{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: crossplane.GetAnnotations(db.Metadata.Managed, db.IdentifierName.SqlDB),
			Labels:      labels,
			Name:        db.IdentifierID.SqlDB,
			Namespace:   namespace,
		},
		Spec: xrds.SQLDatabaseSpec{
			Parameters: xrds.SQLDatabaseParameters{
				MachineType:       &db.MachineType,
				NetworkRef:        &db.IdentifierID.Network,
				PasswordNamespace: &namespace,
				PasswordRef:       &passwordRef,
				ProviderRef:       &db.IdentifierID.Provider,
				Region:            &db.Region,
				ResourceName:      &db.IdentifierID.SqlDB,
				SqlType:           aws.toAWSEngine(db.SQLTypeVersion.Type),
				SqlVersion:        &db.SQLTypeVersion.Version,
				StorageGB:         &storageSize,
				StorageGBLimit:    &maxStorageSize,
				SubnetGroupName:   &subnetGroupName,
				Subnet1:           &subnet1,
				Subnet2:           &subnet2,
			},
		},
	}, errors.OK
}

func (aws *awsRepository) toPasswordSecret(ctx context.Context, db *instance.SqlDB) (*v1.Secret, errors.Error) {
	if db.Auth.RootPassword == "" {
		return nil, errors.BadRequest.WithMessage(fmt.Sprintf(
			"Must provide a master user password for db %s",
			db.IdentifierID.SqlDB,
		))
	}

	name := fmt.Sprintf("%s-password", db.IdentifierID.SqlDB)
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ctx.Value(context.CurrentNamespace).(string),
		},
		Data: map[string][]byte{
			"Password": []byte(db.Auth.RootPassword),
		},
	}, errors.OK
}

func (aws *awsRepository) toAWSEngine(sqlType instance.SQLType) *string {
	var engine string
	switch sqlType {
	case instance.MySqlSQLType:
		engine = "mysql"
		break
	case instance.PostgresSQLType:
		engine = "postgres"
		break
	case instance.SQLServerType:
		engine = "sqlserver-se"
		break
	}
	return &engine
}

func (aws *awsRepository) getPasswordSecret(ctx context.Context, name *string, namespace *string) (*v1.Secret, errors.Error) {
	existingPasswordSecret := &v1.Secret{}
	if err := kubernetes.Client().Client.Get(
		ctx,
		types.NamespacedName{Name: *name, Namespace: *namespace}, existingPasswordSecret); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("Password secret %s not found", *name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to get password secret %s", *name))
	}
	return existingPasswordSecret, errors.OK
}
