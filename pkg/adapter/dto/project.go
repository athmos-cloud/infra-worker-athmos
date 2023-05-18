package dto

type GetProjectResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type ListProjectResponse []ListProjectResponseItem

type ListProjectResponseItem struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Owner          string `json:"owner"`
	RedirectionURL string `json:"redirection_url"`
}

type CreateProjectRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type CreateProjectResponse struct {
	ID             string `json:"id"`
	RedirectionURL string `json:"redirection_url"`
}

type UpdateProjectRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
