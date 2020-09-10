package cmd

import (
	"encoding/json"
	"fmt"
	"time"
)

type GitlabJobs struct {
	Jobs []GitlabJob
}

type GitlabJob struct {
	Id int `json:"id"`
	Status string
	StartedAt time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Duration float64
	Name string
}

type GitlabPipeline struct {
	Id int `json:"id"`
	Status string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Duration time.Duration
}

type Report struct {
	Pipelines []GitlabPipeline
	PipelineSuccessCount int
	PipelineFailedCount int
}


func NewReport(pipelines []GitlabPipeline) *Report {

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
		Pipelines: pipelines,
		PipelineSuccessCount: pipelineSuccessCounter,
		PipelineFailedCount: pipelineFailCounter,
	}
}

func (report *Report) niceString() string {
	pipelineCount := len(report.Pipelines)
	out := fmt.Sprintf("Pipeline Count: %d\n", pipelineCount)
	out += fmt.Sprintf("Failed Pipelines:  %d\n", report.PipelineFailedCount)
	out += fmt.Sprintf("Succeeded Pipelines:  %d\n", report.PipelineSuccessCount)
	return out
}






func NewGitlabPipeline(pipeline GitlabPipeline) *GitlabPipeline {
	duration := pipeline.UpdatedAt.Sub(pipeline.CreatedAt)
	return &GitlabPipeline{
		Id: pipeline.Id,
		Status: pipeline.Status,
		CreatedAt: pipeline.CreatedAt,
		UpdatedAt: pipeline.UpdatedAt,
		Duration: duration,
	}
}

func (job *GitlabJob) niceString() string {

	durationInMinutes := 0.0
	if job.Duration != 0 {
		durationInMinutes = job.Duration / 60
	}

	return fmt.Sprintf("%d \t %s \t %s \t %s \t %f \t %s",
		job.Id, job.Status, job.StartedAt, job.FinishedAt, durationInMinutes, job.Name)
}

func (p *GitlabPipeline) niceString() string {

	duration := p.UpdatedAt.Sub(p.CreatedAt).Minutes()

	return fmt.Sprintf("%d \t %s \t %s \t %s \t %f",
		p.Id, p.Status, p.CreatedAt, p.UpdatedAt, duration)
}

func convertJsonToGitlabJobs(jsonData []byte) (*[]GitlabJob, error) {
	jobs :=  make([]GitlabJob, 0)
	err := json.Unmarshal(jsonData, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func convertJsonToGitlabPipelines(jsonData []byte) (*[]GitlabPipeline, error) {
	pipelines :=  make([]GitlabPipeline, 0)
	err := json.Unmarshal(jsonData, &pipelines)
	if err != nil {
		return nil, err
	}

	pipelinesWithDuration :=  make([]GitlabPipeline, 0)

	for _, p := range(pipelines) {
		withDuration := NewGitlabPipeline(p)
		pipelinesWithDuration = append(pipelinesWithDuration, *withDuration)
	}

	return &pipelinesWithDuration, nil
}

