package secret

type CreateSecretRequest struct {
	ProjectID   string `json:"projectID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

type CreateSecretResponse struct {
}

type GetSecretRequest struct {
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`
}

type GetSecretResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListSecretRequest struct {
	ProjectID string `json:"projectID"`
}

type ListSecretResponse = []GetSecretResponse

type UpdateSecretRequest struct {
}

type DeleteSecretRequest struct {
}
