package cmd

import (
	"bytes"
	"github.com/polpettone/pgcli/cmd/models"
	"io/ioutil"
	"testing"
	"time"
)

var gitlabAPIClientMock = MockGitlabAPIClient{
	GitlabAPIToken:   "dummy",
	GitlabProjectURL: "dummy",
	ProjectID:        "dummy",
}

type MockGitlabAPIClient struct {
	GitlabAPIToken     string
	GitlabProjectURL    string
	ProjectID          string
}

func (m MockGitlabAPIClient) getJobs(pipelineId string) ([]models.Job, error) {
	panic("implement me")
}

func (m MockGitlabAPIClient) getPipelines(status string) ([]models.Pipeline, error) {

	var p0 = models.Pipeline{
		Id:        0,
		Status:    "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	var p1 = models.Pipeline{
		Id:        1,
		Status:    "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	return []models.Pipeline{p0, p1}, nil
}

func (m MockGitlabAPIClient) getLog(jobID string) (string, error) {
	panic("implement me")
}

func (m MockGitlabAPIClient) getLastFailLog() (string, error) {
	panic("implement me")
}

func TestExecute(t *testing.T) {

	cmd := NewPipelinesCmd(gitlabAPIClientMock)
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	if string(out) == "" {
		t.Fatalf("kaputt")
	}
}