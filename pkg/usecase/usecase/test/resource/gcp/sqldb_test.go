package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/sql/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

type wantSqlDB struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.DatabaseInstanceSpec
}

func initSqlDB(t *testing.T) (context.Context, *testResource.TestResource, usecase.SqlDB) {
	ctx, testNet, nuc := initNetwork(t)
	parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
	req := dto.CreateNetworkRequest{
		ParentIDProvider: &parentID,
		Name:             "test-net",
		Managed:          false,
	}
	ctx.Set(context.RequestKey, req)
	net := &network.Network{}
	err := nuc.Create(ctx, net)
	require.True(t, err.IsOk())
	ctx.Set(testResource.NetworkIDKey, net.IdentifierID)
	uc := usecase.NewSqlDBUseCase(testNet.ProjectRepo, gcp.NewRepository(), nil, nil)

	return ctx, testNet, uc
}

func clearSqlDB(ctx context.Context) {
	clearNetwork(ctx)
	dbs := &v1beta1.DatabaseInstanceList{}

	err := kubernetes.Client().Client.List(ctx, dbs)
	if err != nil {
		return
	}
	for _, db := range dbs.Items {
		err = kubernetes.Client().Client.Delete(ctx, &db)
		if err != nil {
			logger.Warning.Printf("Error deleting db %s: %v", db.Name, err)
			continue
		}
	}
}

func Test_sqlDBUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, duc := initSqlDB(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearSqlDB(ctx)
	}()
	t.Run("Create a valid db should succeed", func(t *testing.T) {
		db := SqlDBFixture(ctx, t, duc)
		kubeResource := &v1beta1.DatabaseInstance{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB}, kubeResource)
		require.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          db.IdentifierID.Provider,
			"identifier.vpc":               db.IdentifierID.VPC,
			"identifier.network":           db.IdentifierID.Network,
			"identifier.sqldb":             db.IdentifierID.SqlDB,
			"name.provider":                "test",
			"name.vpc":                     "test",
			"name.network":                 "test-net",
			"name.sqldb":                   "test-sqldb",
		}

		version := "POSTGRES_12"
		diskType := "pd-ssd"
		diskSize := float64(10)
		resizeLimit := float64(0)
		tier := "db-f1-micro"
		wantSpec := v1beta1.DatabaseInstanceSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Delete",
				ProviderConfigReference: &v1.Reference{
					Name: db.IdentifierID.Provider,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.DatabaseInstanceParameters{
				DatabaseVersion: &version,
				RootPasswordSecretRef: &v1.SecretKeySelector{
					Key: "password",
					SecretReference: v1.SecretReference{
						Name:      db.IdentifierID.SqlDB,
						Namespace: ctx.Value(context.CurrentNamespace).(string),
					},
				},
				Region: &db.Region,
				Settings: []v1beta1.SettingsParameters{
					{
						Tier:                &tier,
						DiskType:            &diskType,
						DiskSize:            &diskSize,
						DiskAutoresize:      &db.Disk.Autoresize,
						DiskAutoresizeLimit: &resizeLimit,
					},
				},
			},
		}
		wantNet := wantSqlDB{
			Name:   db.IdentifierID.SqlDB,
			Labels: wantLabels,
			Spec:   wantSpec,
		}
		gotNet := wantSqlDB{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assert.Equal(t, wantNet, gotNet)
	})

	t.Run("Create a DB with an already existing name should fail", func(t *testing.T) {
		db := SqlDBFixture(ctx, t, duc)
		ctx.Set(context.RequestKey, dto.CreateSqlDBRequest{
			Name: db.IdentifierName.SqlDB,
			ParentID: identifier.Network{
				Provider: db.IdentifierID.Provider,
				VPC:      db.IdentifierID.VPC,
				Network:  db.IdentifierID.Network,
			},
		})
		toCreate := &instance.SqlDB{}
		err := duc.Create(ctx, toCreate)
		require.Equal(t, errors.Conflict.Code, err.Code)
	})

}

func Test_sqlDBUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, duc := initSqlDB(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearSqlDB(ctx)
	}()
	t.Run("Delete an existing db should succeed", func(t *testing.T) {
		db := SqlDBFixture(ctx, t, duc)
		ctx.Set(context.RequestKey, dto.DeleteSqlDBRequest{IdentifierID: db.IdentifierID})
		toDelete := &instance.SqlDB{}
		err := duc.Delete(ctx, toDelete)
		require.Equal(t, errors.NoContent.Code, err.Code)
		kubeResource := &v1beta1.DatabaseInstance{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB}, kubeResource)
		require.True(t, k8serrors.IsNotFound(errk))
	})

	t.Run("Delete a DB which does not exist should return NotFound error", func(t *testing.T) {
		ctx.Set(context.RequestKey, dto.DeleteSqlDBRequest{
			IdentifierID: identifier.SqlDB{
				Provider: "test",
				VPC:      "test",
				Network:  "test-net",
				SqlDB:    "this-db-does-not-exist",
			},
		})
		toDelete := &instance.SqlDB{}
		err := duc.Delete(ctx, toDelete)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_sqlDBUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, duc := initSqlDB(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearSqlDB(ctx)
	}()
	t.Run("Get an existing db should succeed", func(t *testing.T) {
		db := SqlDBFixture(ctx, t, duc)
		ctx.Set(context.RequestKey, dto.GetSqlDBRequest{IdentifierID: db.IdentifierID})
		toGet := &instance.SqlDB{}
		err := duc.Get(ctx, toGet)
		require.Equal(t, errors.OK.Code, err.Code)
	})

	t.Run("Get a DB which does not exist should return NotFound error", func(t *testing.T) {
		ctx.Set(context.RequestKey, dto.GetSqlDBRequest{
			IdentifierID: identifier.SqlDB{
				Provider: "test",
				VPC:      "test",
				Network:  "test-net",
				SqlDB:    "this-db-does-not-exist",
			},
		})
		toGet := &instance.SqlDB{}
		err := duc.Get(ctx, toGet)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_sqlDBUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, duc := initSqlDB(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearSqlDB(ctx)
	}()
	t.Run("Update an existing db should succeed", func(t *testing.T) {
		db := SqlDBFixture(ctx, t, duc)
		machineType := "db-g1-small"
		ctx.Set(context.RequestKey, dto.UpdateSqlDBRequest{
			IdentifierID: db.IdentifierID,
			MachineType:  &machineType,
			Disk: &instance.SqlDbDisk{
				Type:    instance.DiskTypeHDD,
				SizeGib: 5,
			},
			SQLType: &instance.SQLTypeVersion{
				Type:    instance.PostgresSQLType,
				Version: "13",
			},
		})
		toUpdate := &instance.SqlDB{}
		err := duc.Update(ctx, toUpdate)
		require.Equal(t, errors.NoContent.Code, err.Code)
		kubeResource := &v1beta1.DatabaseInstance{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: db.IdentifierID.SqlDB}, kubeResource)
		require.NoError(t, errk)
		require.Equal(t, machineType, *kubeResource.Spec.ForProvider.Settings[0].Tier)
		require.Equal(t, float64(5), *kubeResource.Spec.ForProvider.Settings[0].DiskSize)
		require.Equal(t, "POSTGRES_13", *kubeResource.Spec.ForProvider.DatabaseVersion)
	})

	t.Run("Update a DB which does not exist should return NotFound error", func(t *testing.T) {
		ctx.Set(context.RequestKey, dto.UpdateSqlDBRequest{
			IdentifierID: identifier.SqlDB{
				Provider: "test",
				VPC:      "test",
				Network:  "test-net",
				SqlDB:    "this-db-does-not-exist",
			},
		})
		toUpdate := &instance.SqlDB{}
		err := duc.Update(ctx, toUpdate)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}
