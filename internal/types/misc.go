package types

type HomeResponse struct {
	Version string `json:"version"`
}

type HomeResponseEnvelope struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    HomeResponse `json:"data"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type HealthResponseEnvelope struct {
	Success bool           `json:"success"`
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    HealthResponse `json:"data"`
}
