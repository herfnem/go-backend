package types

type PostResponse struct {
	Slug string `json:"slug"`
}

type PostResponseEnvelope struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    PostResponse `json:"data"`
}

type PostCreateResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type PostCreateResponseEnvelope struct {
	Success bool               `json:"success"`
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    PostCreateResponse `json:"data"`
}
