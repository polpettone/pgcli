package adapter

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/config"
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
	Logging          *config.Logging
}

func NewGitlabAPIClient() *GitlabAPIClient {

	loggingEnabled := viper.GetBool("logging_enabled")
	logging :=         config.NewLogging(loggingEnabled)
	state, err := config.ReadState("/home/esteban/.config/pgcli/state.json")

	var projectID string
	if err != nil {
		logging.ErrorLog.Printf("Could not read state, using default project ID from config. %v", err)
		projectID = viper.GetString("project_id")
	} else {
		projectID = state.CurrentProject
	}

	return &GitlabAPIClient{
		GitlabAPIToken:   viper.GetString("api_token"),
		GitlabProjectURL: viper.GetString("url"),
		ProjectID:        projectID,
		Logging:          logging,
	}
}

func (gitlabAPIClient *GitlabAPIClient) GetProjects() ([]models.Project, error) {

	var url = "https://gitlab.com/api/v4/groups/6192951/projects/?membership=true&simple=true&per_page=30&include_subgroups=true"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}

	gitlabAPIClient.Logging.DebugLog.Printf("%v", req)
	resp, err := client.Do(req)
	gitlabAPIClient.Logging.DebugLog.Printf("%v", resp)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	projects, err := models.ConvertJsonToProjects(body)
	if err != nil {
		return nil, err
	}

	return *projects, err
}

func (gitlabAPIClient *GitlabAPIClient) GetJobs(pipelineId string) ([]models.Job, error) {

	var url = gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines/" + pipelineId + "/jobs"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}

	gitlabAPIClient.Logging.DebugLog.Printf("%v", req)
	resp, err := client.Do(req)
	gitlabAPIClient.Logging.DebugLog.Printf("%v", resp)

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

func (gitlabAPIClient GitlabAPIClient) GetPipelines(status string, count int) ([]*models.Pipeline, error) {
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

	gitlabAPIClient.Logging.DebugLog.Printf("%v", req)

	resp, err := client.Do(req)

	gitlabAPIClient.Logging.DebugLog.Printf("%v", resp)

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
	return pipelines, nil
}

type enrichedPipelineResult struct {
	index    int
	pipeline *models.Pipeline
	err      error
}

//TODO: needs refactoring
func (gitlabAPIClient GitlabAPIClient) EnrichPipelinesByUser(pipelines []*models.Pipeline, concurrencyLimit int) ([]*models.Pipeline, error) {

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	enrichedPiplineChan := make(chan *enrichedPipelineResult)

	defer func() {
		close(semaphoreChan)
		close(enrichedPiplineChan)
	}()

	for i, pipeline := range pipelines {
		go func(i int, pipeline *models.Pipeline) {
			semaphoreChan <- struct{}{}
			enrichedPipeline, err := gitlabAPIClient.GetPipeline(strconv.Itoa(pipeline.Id))
			enrichedPipelineResult := &enrichedPipelineResult{i, enrichedPipeline, err}
			enrichedPiplineChan <- enrichedPipelineResult
			<-semaphoreChan
		}(i, pipeline)
	}

	var enrichedPipelineResults []enrichedPipelineResult
	for {
		enrichedPipeline := <-enrichedPiplineChan
		enrichedPipelineResults = append(enrichedPipelineResults, *enrichedPipeline)
		if len(enrichedPipelineResults) == len(pipelines) {
			break
		}
	}
	var enrichedPipelines []*models.Pipeline

	sort.Slice(enrichedPipelineResults, func(i, j int) bool {
		return enrichedPipelineResults[i].index < enrichedPipelineResults[j].index
	})

	for _, e := range enrichedPipelineResults {
		enrichedPipelines = append(enrichedPipelines, e.pipeline)
	}

	return enrichedPipelines, nil
}

//TODO: needs refactoring
func (gitlabAPIClient GitlabAPIClient) EnrichPipelinesByJobs(pipelines []*models.Pipeline, concurrencyLimit int) ([]*models.Pipeline, error) {

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	enrichedPiplineChan := make(chan *enrichedPipelineResult)

	defer func() {
		close(semaphoreChan)
		close(enrichedPiplineChan)
	}()

	for i, pipeline := range pipelines {

		go func(i int, pipeline *models.Pipeline) {
			semaphoreChan <- struct{}{}
			jobs, err := gitlabAPIClient.GetJobs(strconv.Itoa(pipeline.Id))
			pipeline.Jobs = jobs
			enrichedPipelineResult := &enrichedPipelineResult{i, pipeline, err}
			enrichedPiplineChan <- enrichedPipelineResult
			<-semaphoreChan
		}(i, pipeline)
	}

	var enrichedPipelineResults []enrichedPipelineResult
	for {
		enrichedPipeline := <-enrichedPiplineChan
		enrichedPipelineResults = append(enrichedPipelineResults, *enrichedPipeline)
		if len(enrichedPipelineResults) == len(pipelines) {
			break
		}
	}
	var enrichedPipelines []*models.Pipeline

	sort.Slice(enrichedPipelineResults, func(i, j int) bool {
		return enrichedPipelineResults[i].index < enrichedPipelineResults[j].index
	})

	for _, e := range enrichedPipelineResults {
		enrichedPipelines = append(enrichedPipelines, e.pipeline)
	}

	return enrichedPipelines, nil
}

func (gitlabAPIClient GitlabAPIClient) GetPipeline(id string) (*models.Pipeline, error) {
	url := gitlabAPIClient.GitlabProjectURL + "/" + gitlabAPIClient.ProjectID + "/pipelines" + "/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabAPIClient.GitlabAPIToken)
	client := &http.Client{Timeout: time.Second * 5}

	gitlabAPIClient.Logging.DebugLog.Printf("%v", req)
	resp, err := client.Do(req)
	gitlabAPIClient.Logging.DebugLog.Printf("%v", resp)

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

func (gitlabAPIClient GitlabAPIClient) GetLog(jobID string) (string, error) {
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

func (gitlabAPIClient GitlabAPIClient) GetLastFailLog() (string, error) {

	pipelines, err := gitlabAPIClient.GetPipelines("failed", 20)
	if err != nil {
		return "", err
	}

	failedPipeline := getLastFailedPipeline(pipelines)
	jobs, err := gitlabAPIClient.GetJobs(strconv.Itoa(failedPipeline.Id))

	if err != nil {
		return "", err
	}

	failedJob := getLastFailedJob(jobs)

	log, err := gitlabAPIClient.GetLog(strconv.Itoa(failedJob.Id))
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
