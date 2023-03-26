package project

type CreateProjectRequest struct {
	ProjectName string `json:"project_name"`
	OwnerID     string `json:"owner_id"`
}

type CreateProjectResponse struct {
	ProjectID string `json:"project_id"`
}
