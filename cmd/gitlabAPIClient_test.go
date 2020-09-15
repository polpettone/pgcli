package cmd

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/models"
	"testing"
	"time"
)


func Test_convertJsonToGitlabJobs(t *testing.T) {


	jsonData := []byte (`
	
		[

			{
				"id" : 123,
				"status": "running",
				"started_at" : "2020-07-20T13:19:16.151Z",
				"finished_at" : "2020-07-20T13:19:16.151Z",
				"some" : "foo",
				"duration" : 32.124

			},
			{
				"id" : 124,
				"status": "skipped",
				"started_at" : "2020-07-20T13:19:16.151Z",
				"some" : "foo"
			}

		]

	`)

	gitlabJobs, err := models.convertJsonToJobs(jsonData)

	fmt.Println(gitlabJobs)
	fmt.Println(err)
}

func Test_getLastFailedPipeline(t *testing.T) {

	var pipelines []models.Pipeline

	p0 := models.Pipeline{
		CreatedAt:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
		Id: 0,
		Status: "failed",
	}

	p1 := models.Pipeline{
		CreatedAt:   time.Date(2021, 12, 17, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
		Id: 1,
		Status: "failed",
	}

	p2 := models.Pipeline{
		CreatedAt:   time.Date(2019, 12, 17, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
		Id: 2,
		Status: "failed",
	}

	pipelines = append(pipelines, p0)
	pipelines = append(pipelines, p1)
	pipelines = append(pipelines, p2)

	lastFailedPipeline := getLastFailedPipeline(pipelines)

	if lastFailedPipeline.Id != 1 {
		t.Errorf("lastFailedPipeline should be 1 not %d", lastFailedPipeline.Id)
	}



}