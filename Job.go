package hypernova_go

type ReactProps interface {}

type Job struct {
	Name string `json:"name"`
	Data ReactProps `json:"data"`
}

type Jobs map[string]Job
