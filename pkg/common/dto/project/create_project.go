package project

type CreateProjectRequest struct {
	ProjectName string `json:"projectName"`
	OwnerID     string `json:"ownerID"`
}

type CreateProjectResponse struct {
	ProjectID string `json:"projectID"`
}
