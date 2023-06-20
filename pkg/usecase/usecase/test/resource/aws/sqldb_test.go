package aws

import (
<<<<<<< HEAD
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws/xrds"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
	"time"
)

type wantSqlDB struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Spec      xrds.SQLDatabaseSpec
}

func Test_sqlDBUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	duc := usecase.NewSqlDBUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer suiteTeardown(ctx, t, mongoC)

	t.Run("Create a valid db should succeed", func(t *testing.T) {
		defer ClearSqlFixtures(ctx)

		db := SqlDBFixture(ctx, t, duc)
		namespace := ctx.Value(test.TestNamespaceContextKey).(string)
		fmt.Println(fmt.Sprintf("Current namespace is %s", namespace))
		kubeResource := &xrds.SQLDatabase{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{
			Name:      db.IdentifierID.SqlDB,
			Namespace: namespace,
		}, kubeResource)
		require.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          db.IdentifierID.Provider,
			"identifier.vpc":               db.IdentifierID.VPC,
			"identifier.network":           db.IdentifierID.Network,
			"identifier.sqldb":             db.IdentifierID.SqlDB,
			"name.provider":                "fixture-provider",
			"name.vpc":                     "test",
			"name.network":                 "fixture-network",
			"name.sqldb":                   "fixture-db",
		}

		version := "12"
		diskSize := float64(10)
		resizeLimit := float64(0)
		machineType := "db.m7g"
		subnetGroup := fmt.Sprintf("%s-subnet-group", db.IdentifierID.SqlDB)
		subnet1 := fmt.Sprintf("%s-subnet1", db.IdentifierID.SqlDB)
		subnet2 := fmt.Sprintf("%s-subnet2", db.IdentifierID.SqlDB)
		sqlType := "postgres"
		passwordRef := fmt.Sprintf("%s-password", db.IdentifierID.SqlDB)
		wantSpec := xrds.SQLDatabaseSpec{
			Parameters: xrds.SQLDatabaseParameters{
				MachineType:       &machineType,
				NetworkRef:        &db.IdentifierID.Network,
				PasswordNamespace: &namespace,
				PasswordRef:       &passwordRef,
				ProviderRef:       &db.IdentifierID.Provider,
				Region:            &db.Region,
				ResourceName:      &db.IdentifierID.SqlDB,
				SqlType:           &sqlType,
				SqlVersion:        &version,
				StorageGB:         &diskSize,
				StorageGBLimit:    &resizeLimit,
				Subnet1:           &subnet1,
				Subnet2:           &subnet2,
				SubnetGroupName:   &subnetGroup,
			},
		}
		wantNet := wantSqlDB{
			Name:      db.IdentifierID.SqlDB,
			Namespace: namespace,
			Labels:    wantLabels,
			Spec:      wantSpec,
		}
		gotNet := wantSqlDB{
			Name:      kubeResource.Name,
			Namespace: kubeResource.Namespace,
			Labels:    kubeResource.Labels,
			Spec:      kubeResource.Spec,
		}
		assert.Equal(t, wantNet, gotNet)
	})

	t.Run("Creating a DB should create a corresponding password secret", func(t *testing.T) {
		defer ClearSqlFixtures(ctx)

		db := SqlDBFixture(ctx, t, duc)
		passwordSecretName := fmt.Sprintf("%s-password", db.IdentifierID.SqlDB)
		kubeResource := &v1.Secret{}
		err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{
			Name:      passwordSecretName,
			Namespace: ctx.Value(test.TestNamespaceContextKey).(string),
		}, kubeResource)

		assert.NoError(t, err)
		assert.Equal(t, db.Auth.RootPassword, string(kubeResource.Data["Password"][:]))
	})

	t.Run("Create a DB with an already existing name should fail", func(t *testing.T) {
		defer ClearSqlFixtures(ctx)

		SqlDBFixture(ctx, t, duc)
		/*		ctx.Set(context.RequestKey, dto.CreateSqlDBRequest{
				Name: db.IdentifierName.SqlDB,
				ParentID: identifier.Network{
					Provider: db.IdentifierID.Provider,
					VPC:      db.IdentifierID.VPC,
					Network:  db.IdentifierID.Network,
				},
			})*/
		err := duc.Create(ctx, &instance.SqlDB{})
		assert.Equal(t, errors.Conflict.Code, err.Code)
	})

	t.Run("Creating a DB with HDD storage should fail", func(t *testing.T) {
		defer ClearSqlFixtures(ctx)

		region := "eu-west-1"
		req := dto.CreateSqlDBRequest{
			ParentID:    ctx.Value(testResource.NetworkIDKey).(identifier.Network),
			Name:        "fixture-db",
			Region:      region,
			MachineType: "db.m7g",
			Disk: instance.SqlDbDisk{
				Type:    instance.DiskTypeHDD,
				SizeGib: 10,
			},
			SQLType:      instance.PostgresSQLType,
			SQLVersion:   "12",
			Managed:      true,
			RootPassword: "proEsgi7656$!",
		}
		ctx.Set(context.RequestKey, req)

		db := &instance.SqlDB{}
		err := duc.Create(ctx, db)
		assert.Equal(t, errors.BadRequest.Code, err.Code)
	})
}

func Test_sqlDBUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	duc := usecase.NewSqlDBUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer suiteTeardown(ctx, t, mongoC)

	t.Run("Delete an existing db should succeed", func(t *testing.T) {
		defer ClearSqlFixtures(ctx)

		db := SqlDBFixture(ctx, t, duc)
		ctx.Set(context.RequestKey, dto.DeleteSqlDBRequest{IdentifierID: db.IdentifierID})
		err := duc.Delete(ctx, db)
		assert.Equal(t, errors.NoContent.Code, err.Code)
	})

	t.Run("Delete a DB which does not exist should return NotFound error", func(t *testing.T) {
		defer ClearSqlFixtures(ctx)

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
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_sqlDBUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	duc := usecase.NewSqlDBUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)
	db := SqlDBFixture(ctx, t, duc)

	defer suiteTeardown(ctx, t, mongoC)

	t.Run("Get an existing db should succeed", func(t *testing.T) {
		ctx.Set(context.RequestKey, dto.GetResourceRequest{Identifier: db.IdentifierID.SqlDB})
		toGet := &instance.SqlDB{}
		err := duc.Get(ctx, toGet)
		assert.Equal(t, errors.OK.Code, err.Code)
	})

	t.Run("Get a DB which does not exist should return NotFound error", func(t *testing.T) {
		ctx.Set(context.RequestKey, dto.GetResourceRequest{
			Identifier: identifier.SqlDB{
				Provider: "test",
				VPC:      "test",
				Network:  "test-net",
				SqlDB:    "this-db-does-not-exist",
			}.SqlDB,
		})
		toGet := &instance.SqlDB{}
		err := duc.Get(ctx, toGet)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_sqlDBUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()

	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	duc := usecase.NewSqlDBUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)
	db := SqlDBFixture(ctx, t, duc)

	defer suiteTeardown(ctx, t, mongoC)

	t.Run("Update an existing db should succeed", func(t *testing.T) {
		machineType := "db.m6i"
		ctx.Set(context.RequestKey, dto.UpdateSqlDBRequest{
			IdentifierID: db.IdentifierID,
			MachineType:  &machineType,
			Disk: &instance.SqlDbDisk{
				Type:    instance.DiskTypeSSD,
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

		time.Sleep(10 * 1000 * time.Millisecond)

		kubeResource := &xrds.SQLDatabase{}
		errk := kubernetes.Client().Client.Get(
			ctx,
			types.NamespacedName{
				Name:      db.IdentifierID.SqlDB,
				Namespace: ctx.Value(test.TestNamespaceContextKey).(string),
			},
			kubeResource,
		)
		require.NoError(t, errk)
		require.Equal(t, machineType, *kubeResource.Spec.Parameters.MachineType)
		require.Equal(t, float64(5), *kubeResource.Spec.Parameters.StorageGB)
		require.Equal(t, "13", *kubeResource.Spec.Parameters.SqlVersion)
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
=======
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_sqlDBUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Create a valid db should succeed", func(t *testing.T) {
		repo := aws.NewRepository()

	})

	t.Run("Create a DB with an already existing name should fail", func(t *testing.T) {
>>>>>>> feat: xrds definition
	})
}
