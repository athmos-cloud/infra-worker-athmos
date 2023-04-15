package application

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/common/dto/project"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
	"reflect"
	"testing"
)

// application := projectService.ProjectService{
// ProjectRepository: mongo.Client,
// }
//
// id1, _ := application.CreateProject(ctx, project.CreateProjectRequest{
// ProjectName: "test1",
// OwnerID:     "toto",
// })
// _, _ = application.CreateProject(ctx, project.CreateProjectRequest{
// ProjectName: "test2",
// OwnerID:     "toto",
// })
// projectByID, err := application.GetProjectByID(ctx, project.GetProjectByIDRequest{
// ProjectID: id1.ProjectID,
// })
// if !err.IsOk() {
// logger.Info.Println("Err: ", err)
// }
// logger.Info.Println(ctx, "Project found with id: ", projectByID)
// err = application.UpdateProjectName(ctx, project.UpdateProjectRequest{
// ProjectID:   id1.ProjectID,
// ProjectName: "test1-updated",
// })
// if !err.IsOk() {
// logger.Error.Println(ctx, "Error: ", err)
// }
// projectAll, err := application.GetProjectByOwnerID(ctx, project.GetProjectByOwnerIDRequest{
// OwnerID: "toto",
// })
//
// err = application.DeleteProject(ctx, project.DeleteRequest{
// ProjectID: id1.ProjectID,
// })
// if !err.IsOk() {
// logger.Error.Println(ctx, "Error: ", err)
// }
// err = application.DeleteProject(ctx, project.DeleteRequest{
// ProjectID: id1.ProjectID,
// })
// if !err.IsOk() {
// logger.Error.Println(ctx, "Error: ", err)
// }
// logger.Info.Println(ctx, "Project found with ownerID: ", projectAll)
func TestProjectService_Create(t *testing.T) {
	type fields struct {
		ProjectRepository repository.IRepository
	}
	type args struct {
		ctx     context.Context
		request project.CreateProjectRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   project.CreateProjectResponse
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &ProjectService{
				ProjectRepository: tt.fields.ProjectRepository,
			}
			got, got1 := ps.CreateProject(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateProject() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateProject() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestProjectService_Delete(t *testing.T) {
	type fields struct {
		ProjectRepository repository.IRepository
	}
	type args struct {
		ctx     context.Context
		request project.DeleteRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &ProjectService{
				ProjectRepository: tt.fields.ProjectRepository,
			}
			if got := ps.DeleteProject(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProjectService_GetByID(t *testing.T) {
	type fields struct {
		ProjectRepository repository.IRepository
	}
	type args struct {
		ctx     context.Context
		request project.GetProjectByIDRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   project.GetProjectByIDResponse
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &ProjectService{
				ProjectRepository: tt.fields.ProjectRepository,
			}
			got, got1 := ps.GetProjectByID(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProjectByID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetProjectByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestProjectService_GetByOwnerID(t *testing.T) {
	type fields struct {
		ProjectRepository repository.IRepository
	}
	type args struct {
		ctx     context.Context
		request project.GetProjectByOwnerIDRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   project.GetProjectByOwnerIDResponse
		want1  errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &ProjectService{
				ProjectRepository: tt.fields.ProjectRepository,
			}
			got, got1 := ps.GetProjectByOwnerID(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProjectByOwnerID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetProjectByOwnerID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestProjectService_Update(t *testing.T) {
	type fields struct {
		ProjectRepository repository.IRepository
	}
	type args struct {
		ctx     context.Context
		request project.UpdateProjectRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errors.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &ProjectService{
				ProjectRepository: tt.fields.ProjectRepository,
			}
			if got := ps.UpdateProjectName(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateProjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}
