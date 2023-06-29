package resourceUc

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type SqlDB interface {
	Get(context.Context, *instance.SqlDB) errors.Error
	Create(context.Context, *instance.SqlDB) errors.Error
	Update(context.Context, *instance.SqlDB) errors.Error
	Delete(context.Context, *instance.SqlDB) errors.Error
}

type sqlDBUseCase struct {
	projectRepo repository.Project
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewSqlDBUseCase(projectRepo repository.Project, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) SqlDB {
	return &sqlDBUseCase{projectRepo: projectRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (suc *sqlDBUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return suc.gcpRepo
	case types.ProviderAWS:
		return suc.awsRepo
		//case types.ProviderAZURE:
		//	return suc.azureRepo
	}
	return nil
}

func (suc *sqlDBUseCase) Get(ctx context.Context, db *instance.SqlDB) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s db not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetResourceRequest)

	err := _setSqlNamespace(ctx, suc)
	if !err.IsOk() {
		return err
	}

	foundDB, err := repo.FindSqlDB(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.Identifier},
	})
	if !err.IsOk() {
		return err
	}
	*db = *foundDB

	return errors.OK
}

func (suc *sqlDBUseCase) Create(ctx context.Context, db *instance.SqlDB) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s db not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateSqlDBRequest)
	defaults.SetDefaults(&req)

	err := _setSqlNamespace(ctx, suc)
	if !err.IsOk() {
		return err
	}

	project, errRepo := suc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errRepo.IsOk() {
		return errRepo
	}
	ctx.Set(context.CurrentNamespace, project.Namespace)
	network, errNet := repo.FindNetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.ParentID.Network},
	})
	if !errNet.IsOk() {
		return errNet
	}

	toCreateDB := &instance.SqlDB{
		Metadata: metadata.Metadata{
			Managed: req.Managed,
			Tags:    req.Tags,
		},
		IdentifierID: identifier.SqlDB{
			Provider: network.IdentifierID.Provider,
			VPC:      network.IdentifierID.VPC,
			Network:  network.IdentifierID.Network,
			SqlDB:    usecase.IdFromName(req.Name),
		},
		IdentifierName: identifier.SqlDB{
			Provider: network.IdentifierName.Provider,
			VPC:      network.IdentifierName.VPC,
			Network:  network.IdentifierName.Network,
			SqlDB:    req.Name,
		},
		MachineType: req.MachineType,
		SQLTypeVersion: instance.SQLTypeVersion{
			Version: req.SQLVersion,
			Type:    req.SQLType,
		},
		Region: req.Region,
		Auth: instance.SqlDBAuth{
			RootPassword: req.RootPassword,
		},
		Disk: req.Disk,
	}
	if err := repo.CreateSqlDB(ctx, toCreateDB); !err.IsOk() {
		return err
	}
	*db = *toCreateDB

	return errors.Created
}

func (suc *sqlDBUseCase) Update(ctx context.Context, db *instance.SqlDB) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s db not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateSqlDBRequest)
	defaults.SetDefaults(&req)

	err := _setSqlNamespace(ctx, suc)
	if !err.IsOk() {
		return err
	}

	foundDB, err := repo.FindSqlDB(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.SqlDB}})
	if !err.IsOk() {
		return err
	}

	*db = *foundDB
	if req.Name != nil {
		db.IdentifierName.SqlDB = *req.Name
	}
	if req.Tags != nil {
		db.Metadata.Tags = *req.Tags
	}
	if req.Managed != nil {
		db.Metadata.Managed = *req.Managed
	}
	if req.SQLType != nil {
		db.SQLTypeVersion = *req.SQLType
	}
	if req.MachineType != nil {
		db.MachineType = *req.MachineType
	}
	if req.Region != nil {
		db.Region = *req.Region
	}
	if req.RootPassword != nil {
		db.Auth.RootPassword = *req.RootPassword
	}
	if req.Disk != nil {
		db.Disk = *req.Disk
	}
	if errUpdate := repo.UpdateSqlDB(ctx, db); !errUpdate.IsOk() {
		return errUpdate
	}

	return errors.NoContent
}

func (suc *sqlDBUseCase) Delete(ctx context.Context, db *instance.SqlDB) errors.Error {
	err := _setSqlNamespace(ctx, suc)
	if !err.IsOk() {
		return err
	}

	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteSqlDBRequest)
	defaults.SetDefaults(&req)

	foundDB, err := repo.FindSqlDB(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID},
	})
	if !err.IsOk() {
		return err
	}
	*db = *foundDB

	if delErr := repo.DeleteSqlDB(ctx, db); !delErr.IsOk() {
		return delErr
	}

	return errors.NoContent
}

func _setSqlNamespace(ctx context.Context, suc *sqlDBUseCase) errors.Error {
	project, err := suc.projectRepo.Find(ctx, option.Option{
		Value: repository.FindProjectByIDRequest{
			ID: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !err.IsOk() {
		return err
	}

	ctx.Set(context.CurrentNamespace, project.Namespace)
	return errors.OK
}
