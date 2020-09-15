package models

import (
	"fmt"
)

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

func (report *Report) NiceString() string {
	pipelineCount := len(report.Pipelines)
	out := fmt.Sprintf("Pipeline Count: %d\n", pipelineCount)
	out += fmt.Sprintf("Failed Pipelines:  %d\n", report.PipelineFailedCount)
	out += fmt.Sprintf("Succeeded Pipelines:  %d\n", report.PipelineSuccessCount)
	return out
}

