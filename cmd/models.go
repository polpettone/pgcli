package cmd

import (
	"encoding/json"
	"fmt"
	"time"
)

type Jobs struct {
	Jobs []Job
}

type Job struct {
	Id         int `json:"id"`
	Status     string
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Duration   float64
	Name       string
}

type Pipeline struct {
	Id        int `json:"id"`
	Status    string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Duration  time.Duration
	PipelineUser PipelineUser `json:"user"`
}

type PipelineUser struct {
	Name      string `json:"name"`
	UserName  string `json:"username"`
	ID        int `json:"id"`
	State     string `json:"state"`
	AvatarUrl string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

type Report struct {
	Pipelines            []Pipeline
	PipelineSuccessCount int
	PipelineFailedCount  int
}

func NewReport(pipelines []Pipeline) *Report {

	pipelineSuccessCounter := 0
	pipelineFailCounter := 0

	for _, p := range pipelines {
		if p.Status == "success" {
			pipelineSuccessCounter++
		}
		if p.Status == "failed" {
			pipelineFailCounter++
		}
	}

	return &Report{
		Pipelines:            pipelines,
		PipelineSuccessCount: pipelineSuccessCounter,
		PipelineFailedCount:  pipelineFailCounter,
	}
}

func (report *Report) niceString() string {
	pipelineCount := len(report.Pipelines)
	out := fmt.Sprintf("Pipeline Count: %d\n", pipelineCount)
	out += fmt.Sprintf("Failed Pipelines:  %d\n", report.PipelineFailedCount)
	out += fmt.Sprintf("Succeeded Pipelines:  %d\n", report.PipelineSuccessCount)
	return out
}

func NewPipeline(pipeline Pipeline) *Pipeline {
	duration := pipeline.UpdatedAt.Sub(pipeline.CreatedAt)
	return &Pipeline{
		Id:        pipeline.Id,
		Status:    pipeline.Status,
		CreatedAt: pipeline.CreatedAt,
		UpdatedAt: pipeline.UpdatedAt,
		PipelineUser: pipeline.PipelineUser,
		Duration:  duration,
	}
}

func (job *Job) niceString() string {

	durationInMinutes := 0.0
	if job.Duration != 0 {
		durationInMinutes = job.Duration / 60
	}

	return fmt.Sprintf("%d \t %s \t %s \t %s \t %f \t %s",
		job.Id, job.Status, job.StartedAt, job.FinishedAt, durationInMinutes, job.Name)
}

func (p *Pipeline) niceString() string {
	return fmt.Sprintf("%d \t %s \t %s \t %s \t %s \t %s",
		p.Id, p.Status, p.CreatedAt, p.UpdatedAt, p.Duration, p.PipelineUser.UserName)
}

func convertJsonToJobs(jsonData []byte) (*[]Job, error) {
	jobs := make([]Job, 0)
	err := json.Unmarshal(jsonData, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func convertJsonToPipeline(jsonData []byte) (*Pipeline, error) {
	var pipeline Pipeline
	err := json.Unmarshal(jsonData, &pipeline)
	if err != nil {
		return nil, err
	}
	return NewPipeline(pipeline), nil
}

func convertJsonToPipelines(jsonData []byte) (*[]Pipeline, error) {
	pipelines := make([]Pipeline, 0)
	err := json.Unmarshal(jsonData, &pipelines)
	if err != nil {
		return nil, err
	}

	pipelinesWithDuration := make([]Pipeline, 0)

	for _, p := range pipelines {
		withDuration := NewPipeline(p)
		pipelinesWithDuration = append(pipelinesWithDuration, *withDuration)
	}

	return &pipelinesWithDuration, nil
}
