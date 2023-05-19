package dto

type GetSecretResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListSecretResponse []ListSecretResponseItem

type ListSecretResponseItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateSecretRequest struct {
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       []byte `json:"value"`
}

type CreateSecretResponse struct {
	ID             string `json:"id"`
	RedirectionURL string `json:"redirection_url"`
}

type UpdateSecretRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       []byte `json:"value"`
}

type DeleteSecretRequest struct {
	Name string `json:"name"`
}
