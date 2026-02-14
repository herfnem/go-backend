package types

type ErrorResponseEnvelope struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type EmptyResponseEnvelope struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
