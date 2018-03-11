package hypernova_go

type JobError struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Stack   []string `json:"stack"`
}

type JobResponse struct {
	Error      JobError `json:"error"`
	Html       string   `json:"html"`
	Duration   float32  `json:"duration"`
	StatusCode int      `json:"statusCode"`
	Success    bool     `json:"success"`
}

type HypernovaResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Results map[string]JobResponse `json:"results"`
}
