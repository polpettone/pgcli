package cmd

import (
	"bytes"
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

func (m MockGitlabAPIClient) getJobs(pipelineId string) ([]GitlabJob, error) {
	panic("implement me")
}

func (m MockGitlabAPIClient) getPipelines(status string) ([]GitlabPipeline, error) {

	var p0 = GitlabPipeline{
		Id:        0,
		Status:    "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	var p1 = GitlabPipeline{
		Id:        1,
		Status:    "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	return []GitlabPipeline {p0, p1}, nil
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