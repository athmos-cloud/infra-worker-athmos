package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-gcp/apis/sql/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func (gcp *gcpRepository) FindSqlDB(ctx context.Context, opt option.Option) (*instance.SqlDB, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpDB := &v1beta1.DatabaseInstance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, gcpDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("SQL database %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get SQL database %s", req.Name))
	}
	mod, err := gcp.toModelSqlDB(ctx, gcpDB)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllSqlDBs(ctx context.Context, opt option.Option) (*instance.SqlDBCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpDBList := &v1beta1.DatabaseInstanceList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpDBList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list vm"))
	}
	modDBs, err := gcp.toModelSqlDBCollection(ctx, gcpDBList)
	if !err.IsOk() {
		return nil, err
	}

	return modDBs, errors.OK
}

func (gcp *gcpRepository) FindAllRecursiveSqlDBs(ctx context.Context, opt option.Option, ch *resourceRepo.SqlDBChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpDBList := &v1beta1.DatabaseInstanceList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpDBList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage("unable to list vm")
		return
	}
	if dbs, err := gcp.toModelSqlDBCollection(ctx, gcpDBList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- dbs
	}
}

func (gcp *gcpRepository) CreateSqlDB(ctx context.Context, db *instance.SqlDB) errors.Error {
	if exists, err := gcp.SqlDBExists(ctx, db); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("db %s already exists in network %s", db.IdentifierName.SqlDB, db.IdentifierID.Network))
	}
	if err := gcp._createSqlPasswordSecret(ctx, db); !err.IsOk() {
		return err
	}
	gcpVM, err := gcp.toGCPSqlDB(ctx, db)
	if !err.IsOk() {
		return err
	}
	if errCreate := kubernetes.Client().Client.Create(ctx, gcpVM); errCreate != nil {
		if k8serrors.IsAlreadyExists(errCreate) {
			return errors.Conflict.WithMessage(fmt.Sprintf("db %s already exists", db.IdentifierName.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create db %s", db.IdentifierName.SqlDB))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateSqlDB(ctx context.Context, db *instance.SqlDB) errors.Error {
	existingDB := &v1beta1.DatabaseInstance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB}, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", db.IdentifierID.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", db.IdentifierID.SqlDB))
	}
	if err := gcp._updateSqlPasswordSecret(ctx, db); !err.IsOk() {
		return err
	}
	gcpVM, err := gcp.toGCPSqlDB(ctx, db)
	if !err.IsOk() {
		return err
	}
	existingDB.Spec = gcpVM.Spec
	existingDB.Labels = gcpVM.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("db %s not found", db.IdentifierName.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update db %s", db.IdentifierName.SqlDB))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteSqlDB(ctx context.Context, db *instance.SqlDB) errors.Error {
	existingDB := &v1beta1.DatabaseInstance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB}, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("db %s not found", db.IdentifierName.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get db %s", db.IdentifierName.SqlDB))
	}
	if err := kubernetes.Client().Client.Delete(ctx, existingDB); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("db %s not found", db.IdentifierName.SqlDB))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete db %s", db.IdentifierName.SqlDB))
	}
	if err := gcp._deleteSqlPasswordSecret(ctx, db); !err.IsOk() {
		return err
	}
	return errors.NoContent
}

func (gcp *gcpRepository) SqlDBExists(ctx context.Context, db *instance.SqlDB) (bool, errors.Error) {
	gcpDBs := &v1beta1.DatabaseInstanceList{}
	parentID := identifier.Network{
		Provider: db.IdentifierID.Provider,
		VPC:      db.IdentifierID.VPC,
		Network:  db.IdentifierID.Network,
	}
	searchLabels := lo.Assign(parentID.ToIDLabels(), map[string]string{identifier.SqlDBNameKey: db.IdentifierName.SqlDB})
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpDBs, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list db")
	}
	return len(gcpDBs.Items) > 0, errors.OK
}

func (gcp *gcpRepository) toModelSqlDB(ctx context.Context, db *v1beta1.DatabaseInstance) (*instance.SqlDB, errors.Error) {
	id := identifier.SqlDB{}
	name := identifier.SqlDB{}
	if err := id.IDFromLabels(db.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(db.Labels); !err.IsOk() {
		return nil, err
	}
	version, err := getDBTypeVersion(*db.Spec.ForProvider.DatabaseVersion)
	logger.Info.Println("version", version)
	if !err.IsOk() {
		return nil, err
	}
	logger.Info.Println("autoresize", *db.Spec.ForProvider.Settings[0].DiskAutoresizeLimit)
	autoResizeLimit := int(*db.Spec.ForProvider.Settings[0].DiskAutoresizeLimit)
	modelDB := &instance.SqlDB{
		Metadata: metadata.Metadata{
			Status:  metadata.StatusFromKubernetesStatus(db.Status.Conditions),
			Managed: db.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
		},
		IdentifierID:   id,
		IdentifierName: name,
		SQLTypeVersion: *version,
		MachineType:    *db.Spec.ForProvider.Settings[0].Tier,
		Region:         *db.Spec.ForProvider.Region,
		Disk: instance.SqlDbDisk{
			Type:               fromGCPDiskType(*db.Spec.ForProvider.Settings[0].DiskType),
			SizeGib:            int(*db.Spec.ForProvider.Settings[0].DiskSize),
			AutoresizeLimitGib: &autoResizeLimit,
			Autoresize:         *db.Spec.ForProvider.Settings[0].DiskAutoresize,
		},
	}
	if errPwd := gcp._getSqlPasswordSecret(ctx, modelDB); !errPwd.IsOk() {
		return nil, errPwd
	}
	return modelDB, errors.OK
}

func (gcp *gcpRepository) toGCPSqlDB(ctx context.Context, db *instance.SqlDB) (*v1beta1.DatabaseInstance, errors.Error) {
	version, err := _getDBVersion(db.SQLTypeVersion)
	if !err.IsOk() {
		return nil, err
	}
	ns, ok := ctx.Value(context.CurrentNamespace).(string)
	if !ok {
		return nil, errors.InternalError.WithMessage("unable to get current namespace")
	}
	diskType := toGCPDiskType(db.Disk.Type)
	diskSize := float64(db.Disk.SizeGib)
	resizeLimit := float64(*db.Disk.AutoresizeLimitGib)
	labelsDB := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), db.IdentifierID.ToIDLabels(), db.IdentifierName.ToNameLabels())
	return &v1beta1.DatabaseInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:        db.IdentifierID.SqlDB,
			Labels:      labelsDB,
			Annotations: crossplane.GetAnnotations(db.Metadata.Managed, db.IdentifierName.SqlDB),
		},
		Spec: v1beta1.DatabaseInstanceSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(db.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: db.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.DatabaseInstanceParameters{
				DatabaseVersion: version,
				RootPasswordSecretRef: &v1.SecretKeySelector{
					Key: "password",
					SecretReference: v1.SecretReference{
						Name:      db.IdentifierID.SqlDB,
						Namespace: ns,
					},
				},
				Region: &db.Region,
				Settings: []v1beta1.SettingsParameters{
					{
						Tier:                &db.MachineType,
						DiskType:            &diskType,
						DiskSize:            &diskSize,
						DiskAutoresize:      &db.Disk.Autoresize,
						DiskAutoresizeLimit: &resizeLimit,
					},
				},
			},
		},
	}, errors.OK
}

func (gcp *gcpRepository) toModelSqlDBCollection(ctx context.Context, instanceList *v1beta1.DatabaseInstanceList) (*instance.SqlDBCollection, errors.Error) {
	items := instance.SqlDBCollection{}
	for _, item := range instanceList.Items {
		sqlDB, err := gcp.toModelSqlDB(ctx, &item)

		if !err.IsOk() {
			return nil, err
		}
		items[sqlDB.IdentifierName.SqlDB] = *sqlDB
	}
	return &items, errors.OK
}

var (
	mySQLAvailableVersion     = []string{"5.6", "5.7", "8.0"}
	postgresAvailableVersion  = []string{"9.6", "10", "11", "12", "13", "14"}
	sqlServerAvailableVersion = []string{"2017_STANDARD", "2017_STANDARD", "2017_ENTERPRISE", "2017_EXPRESS_2017", "2017_WEB", "2019_STANDARD", "2019_ENTERPRISE", "2019_EXPRESS", "2019_WEB"}
)

func _getDBVersion(dbVersion instance.SQLTypeVersion) (*string, errors.Error) {
	isAvailable := func(availableVersions []string, version string) bool {
		for _, v := range availableVersions {
			if v == version {
				return true
			}
		}
		return false
	}
	sqlTypeVersion := func(sqlType string, version string) *string {
		res := fmt.Sprintf("%s_%s", sqlType, version)
		return &res
	}
	switch dbVersion.Type {
	case instance.MySqlSQLType:
		if !isAvailable(mySQLAvailableVersion, dbVersion.Version) {
			return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid mysql version %s", dbVersion.Version))
		}
		return sqlTypeVersion("MYSQL", dbVersion.Version), errors.OK
	case instance.PostgresSQLType:
		if !isAvailable(postgresAvailableVersion, dbVersion.Version) {
			return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid postgres version %s", dbVersion.Version))
		}
		return sqlTypeVersion("POSTGRES", dbVersion.Version), errors.OK
	case instance.SQLServerType:
		if !isAvailable(sqlServerAvailableVersion, dbVersion.Version) {
			return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid sql server version %s", dbVersion.Version))
		}
		return sqlTypeVersion("SQLSERVER", dbVersion.Version), errors.OK
	}
	return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid sql type %s", dbVersion.Type))
}

func getDBTypeVersion(dbVersion string) (*instance.SQLTypeVersion, errors.Error) {
	split := strings.Split(dbVersion, "_")
	if len(split) < 2 {
		return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid db version %s", dbVersion))
	}
	sqlVersion := strings.Join(split[1:], "_")
	var sqlType instance.SQLType
	switch split[0] {
	case "MYSQL":
		sqlType = instance.MySqlSQLType
	case "POSTGRES":
		sqlType = instance.PostgresSQLType
	case "SQLSERVER":
		sqlType = instance.SQLServerType
	default:
		return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid db type %s", split[0]))
	}
	return &instance.SQLTypeVersion{Version: sqlVersion, Type: sqlType}, errors.OK
}
