package cmd

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/models"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)


type GitlabAPIClient struct {
	GitlabAPIToken   string
	GitlabProjectURL string
	ProjectID        string
	Logging          *Logging
}

func NewGitlabAPIClient() *GitlabAPIClient {
	return &GitlabAPIClient{
		GitlabAPIToken:   viper.GetString("api_token"),
		GitlabProjectURL: viper.GetString("url"),
		ProjectID:        viper.GetString("project_id"),
		Logging:          NewLogging(),
	}
}

func (gitlabAPIClient *GitlabAPIClient) getJobs(pipelineId string) ([]models.Job, error) {

	var url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines/" + pipelineId + "/jobs"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jobs, err := models.ConvertJsonToJobs(body)
	if err != nil {
		return nil, err
	}

	return *jobs, nil
}

func (gitlabAPIClient GitlabAPIClient) getPipelines(status string, withUser bool, count int, withCommitTitle bool) ([]*models.Pipeline, error) {

	pipelineCount := strconv.Itoa(count)

	var url string
	if status == "" {
		url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines?order_by=updated_at&per_page=" + pipelineCount
	} else {
		url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines?order_by=updated_at&per_page=" + pipelineCount + "&status=" + status
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pipelines, err := models.ConvertJsonToPipelines(body)
	if err != nil {
		return nil, err
	}

	if withUser {
		var enrichedPipelines []*models.Pipeline
		for _, p := range pipelines {
			enriched, err := gitlabAPIClient.getPipeline(p.Id)
			if err != nil {
				return nil, err
			}
			enrichedPipelines = append(enrichedPipelines, enriched)
		}
		pipelines = enrichedPipelines
	}

	if withCommitTitle {
		for _, p := range pipelines {
			jobs, err := gitlabAPIClient.getJobs(strconv.Itoa(p.Id))
			if err != nil {
				return nil, err
			}
			p.Jobs = jobs
		}
	}

	return pipelines, nil
}

func (gitlabAPIClient GitlabAPIClient) getPipeline(id int) (*models.Pipeline, error) {
	url := gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines" + "/" + strconv.Itoa(id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pipeline, err := models.ConvertJsonToPipeline(body)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (gitlabAPIClient GitlabAPIClient) getLog(jobID string) (string, error) {
	var url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/jobs/" + jobID + "/trace"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (gitlabAPIClient GitlabAPIClient) getLastFailLog() (string, error) {

	pipelines, err := gitlabAPIClient.getPipelines("failed", false, 20, false)
	if err != nil {
		return "", err
	}

	failedPipeline := getLastFailedPipeline(pipelines)
	jobs, err := gitlabAPIClient.getJobs(strconv.Itoa(failedPipeline.Id))

	if err != nil {
		return "", err
	}

	failedJob := getLastFailedJob(jobs)

	log, err := gitlabAPIClient.getLog(strconv.Itoa(failedJob.Id))
	if err != nil {
		return "", err
	}

	fmt.Printf(failedJob.NiceString())

	return log, nil
}

func getLastFailedPipeline(pipelines []*models.Pipeline) *models.Pipeline {
	var failedPipeline *models.Pipeline

	sort.Slice(pipelines[:], func(i, j int) bool {
		return pipelines[i].CreatedAt.Before(pipelines[j].CreatedAt)
	})

	for _, p := range pipelines {
		if p.Status == "failed" {
			failedPipeline = p
		}
	}
	return failedPipeline
}

func getLastFailedJob(jobs []models.Job) models.Job {
	var failedJob models.Job

	sort.Slice(jobs[:], func(i, j int) bool {
		return jobs[i].StartedAt.Before(jobs[j].StartedAt)
	})

	for _, j := range jobs {
		if j.Status == "failed" {
			failedJob = j
		}
	}
	return failedJob
}
