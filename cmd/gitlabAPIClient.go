package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)


type APIClient interface {
	getJobs(pipelineId string) ([]Job, error)
	getPipelines(status string, withUser bool, count int) ([]Pipeline, error)
	getLog(jobID string) (string, error)
	getLastFailLog() (string, error)
}

type GitlabAPIClient struct {
	GitlabAPIToken     string
	GitlabProjectURL    string
	ProjectID          string
}

func NewGitlabAPIClient(apiToken string, projectURL string, projectID string) APIClient {
	return &GitlabAPIClient{
		GitlabAPIToken:   apiToken,
		GitlabProjectURL: projectURL,
		ProjectID:        projectID,
	}
}

func (gitlabAPIClient *GitlabAPIClient) getJobs(pipelineId string) ([]Job, error) {

	var url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines/" + pipelineId + "/jobs"
	req, err := http.NewRequest("GET",  url, nil)

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

	jobs, err := convertJsonToJobs(body)
	if err != nil {
		return nil, err
	}

	return *jobs, nil
}

func (gitlabAPIClient GitlabAPIClient) getPipelines(status string, withUser bool, count int) ([]Pipeline, error) {

	pipelineCount := strconv.Itoa(count)

	var url string
	if status == "" {
		 url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines?per_page="+pipelineCount
	} else {
		 url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines?per_page="+pipelineCount+"&status=" + status
	}

	req, err := http.NewRequest("GET",  url, nil)
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

	pipelines, err := convertJsonToPipelines(body)
	if err != nil {
		return nil, err
	}

	var sortedPipelines []Pipeline

	if withUser {
		var enrichedPipelines []Pipeline
		for _, p := range *pipelines {
			enriched, err := gitlabAPIClient.getPipeline(p.Id)
			if err != nil {
				return nil, err
			}
			enrichedPipelines = append(enrichedPipelines, *enriched)
		}
		sortedPipelines = enrichedPipelines
	} else {
		sortedPipelines = *pipelines
	}

	sort.Slice(sortedPipelines[:], func(i, j int) bool {
		return sortedPipelines[i].CreatedAt.After(sortedPipelines[j].CreatedAt)
	})

	return sortedPipelines, nil
}


func (gitlabAPIClient GitlabAPIClient) getPipeline(id int) (*Pipeline, error) {
	url := gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines" + "/" + strconv.Itoa(id)

	req, err := http.NewRequest("GET",  url, nil)
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

	pipeline, err := convertJsonToPipeline(body)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}


func (gitlabAPIClient GitlabAPIClient) getLog(jobID string) (string, error) {
	var url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/jobs/"+ jobID +"/trace"
	req, err := http.NewRequest("GET",  url, nil)
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

	pipelines, err := gitlabAPIClient.getPipelines("failed", false, 20)
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

	fmt.Printf(failedJob.niceString())

	return log, nil
}

func getLastFailedPipeline(pipelines []Pipeline) Pipeline {
	var failedPipeline Pipeline

	sort.Slice(pipelines[:], func(i, j int) bool {
		return pipelines[i].CreatedAt.Before(pipelines[j].CreatedAt)
	})

	for _, p  := range pipelines {
		if p.Status == "failed"	{
			failedPipeline = p
		}
	}
	return failedPipeline
}

func getLastFailedJob(jobs []Job) Job {
	var failedJob Job

	sort.Slice(jobs[:], func(i, j int) bool {
		return jobs[i].StartedAt.Before(jobs[j].StartedAt)
	})

	for _, j  := range jobs {
		if j.Status == "failed"	{
			failedJob = j
		}
	}
	return failedJob
}
