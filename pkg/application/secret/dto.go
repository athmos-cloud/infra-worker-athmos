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
}

type GetSecretResponse struct {
}

type ListSecretRequest struct {
}

type ListSecretResponse struct {
}

type UpdateSecretRequest struct {
}

type DeleteSecretRequest struct {
}
