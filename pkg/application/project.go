package application

import (
	"context"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/project"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	projectRepository "github.com/PaulBarrie/infra-worker/pkg/repository/project"
)

type ProjectService struct {
	ProjectRepository *projectRepository.Repository
}

func (ps *ProjectService) CreateProject(ctx context.Context, request dto.CreateProjectRequest) (dto.CreateProjectResponse, errors.Error) {
	resp, err := ps.ProjectRepository.Create(ctx, option.Option{
		Value: request,
	})
	if !err.IsOk() {
		return dto.CreateProjectResponse{}, err
	}
	return resp.(dto.CreateProjectResponse), errors.OK
}

func (ps *ProjectService) UpdateProjectName(ctx context.Context, request dto.UpdateProjectRequest) errors.Error {
	err := ps.ProjectRepository.Update(ctx, option.Option{
		Value: request,
	})
	if !err.IsOk() {
		return err
	}
	return errors.OK
}

func (ps *ProjectService) GetProjectByID(ctx context.Context, request dto.GetProjectByIDRequest) (dto.GetProjectByIDResponse, errors.Error) {
	resp, err := ps.ProjectRepository.Create(ctx, option.Option{
		Value: request,
	})
	if !err.IsOk() {
		return dto.GetProjectByIDResponse{}, err
	}
	return resp.(dto.GetProjectByIDResponse), errors.OK
}

func (ps *ProjectService) GetProjectByOwnerID(ctx context.Context, request dto.GetProjectByOwnerIDRequest) (dto.GetProjectByOwnerIDResponse, errors.Error) {
	resp, err := ps.ProjectRepository.List(ctx, option.Option{
		Value: request,
	})
	if !err.IsOk() {
		return dto.GetProjectByOwnerIDResponse{}, err
	}
	return resp.(dto.GetProjectByOwnerIDResponse), errors.OK
}

func (ps *ProjectService) DeleteProject(ctx context.Context, request dto.DeleteRequest) errors.Error {
	return ps.ProjectRepository.Delete(ctx, option.Option{
		Value: request,
	})
}
