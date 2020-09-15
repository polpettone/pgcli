package models

import (
	"encoding/json"
	"fmt"
	"time"
)

func NewPipeline(pipeline Pipeline) *Pipeline {
	duration := pipeline.UpdatedAt.Sub(pipeline.CreatedAt)
	return &Pipeline{
		Id:           pipeline.Id,
		Status:       pipeline.Status,
		CreatedAt:    pipeline.CreatedAt,
		UpdatedAt:    pipeline.UpdatedAt,
		PipelineUser: pipeline.PipelineUser,
		Duration:     duration,
	}
}

func (p *Pipeline) NiceString() string {
	return fmt.Sprintf("%d \t %s \t %s \t %s \t %s \t %s",
		p.Id, p.Status, p.CreatedAt, p.UpdatedAt, p.Duration, p.PipelineUser.UserName)
}

func ConvertJsonToPipeline(jsonData []byte) (*Pipeline, error) {
	var pipeline Pipeline
	err := json.Unmarshal(jsonData, &pipeline)
	if err != nil {
		return nil, err
	}
	return NewPipeline(pipeline), nil
}

func ConvertJsonToPipelines(jsonData []byte) (*[]Pipeline, error) {
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

type Pipeline struct {
	Id           int `json:"id"`
	Status       string
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Duration     time.Duration
	PipelineUser PipelineUser `json:"user"`
}

type PipelineUser struct {
	Name      string `json:"name"`
	UserName  string `json:"username"`
	ID        int    `json:"id"`
	State     string `json:"state"`
	AvatarUrl string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}
