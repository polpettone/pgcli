package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Jobs struct {
	Jobs []Job
}

type Commit struct {
	Title string `json:"title"`
	Message string `json:"message"`
}

type Job struct {
	Id         int `json:"id"`
	Status     string
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Duration   float64
	Name       string
	Commit	   Commit `json:"commit"`
	Stage	   string `json:"stage"`
}


func (job Job) NiceString() string {
	durationInMinutes := 0.0
	if job.Duration != 0 {
		durationInMinutes = job.Duration / 60
	}
	return fmt.Sprintf("%d \t %s \t %s \t %s \t %f \t %s",
		job.Id, job.Status, job.StartedAt, job.FinishedAt, durationInMinutes, job.Name)
}


func (job Job) View() string {
	output := fmt.Sprintf("%s \t %s \t %s", job.Name, job.Status, job.Commit.Title)
	return output
}


func ConvertJsonToJobs(jsonData []byte) (*[]Job, error) {
	jobs := make([]Job, 0)
	err := json.Unmarshal(jsonData, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}
