package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/polpettone/pgcli/cmd/models"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func StatusCmd(apiClient *GitlabAPIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "status <pipelineID>",
		Short: "shows the status of the pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := handleStatusCommand(args, apiClient)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func handleStatusCommand(args []string, apiClient *GitlabAPIClient) error {
	var pipelineId string
	if len(args) < 1 || args[0] == "" {
		pipelines, _ := apiClient.getPipelines("", 10)
		pipeline, _ := showPipelineSelectionPrompt(pipelines)
		pipelineId = strconv.Itoa(pipeline.Id)
	} else {
		pipelineId = args[0]
	}

	jobs, err := apiClient.getJobs(pipelineId)
	if err != nil {
		return err
	}

	value := ""
	for _, job := range jobs {
		value = value + "\n" + job.View()
	}

	var data [][]string
	for _, job := range jobs {
		data = append(data, []string{
			job.Name,
			job.Status,
			job.StartedAt.Format("02.01.2006 15:04:05"),
			job.FinishedAt.Format("02.01.2006 15:04:05")})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Started", "Finished"})
	for _, d := range data {
		table.Append(d)
	}

	pipeline, err := apiClient.getPipeline(pipelineId)

	if err != nil {
		return err
	}

	pipelines := []*models.Pipeline{pipeline}
	enrichedPipelines, err := apiClient.enrichPipelinesByUser(pipelines, 1)

	if err != nil {
		return err
	}
	enrichedPipeline := enrichedPipelines[0]

	pipelineSummary := fmt.Sprintf("%s \n"+
		"%s \t %s \n"+
		"%s \t %s \n"+
		"%s \t %s \n"+
		"%s \t %s \n"+
		"%s \t %s \n"+
		"%s \t %f \n",
		"Pipeline Summary",
		"Status:", enrichedPipeline.Status,
		"Commit:", jobs[0].Commit.Title,
		"Committer:", enrichedPipeline.PipelineUser.Name,
		"Created At:", enrichedPipeline.CreatedAt.String(),
		"Updated At:", enrichedPipeline.UpdatedAt.String(),
		"Duration Total:", enrichedPipeline.Duration.Minutes())

	fmt.Println(pipelineSummary)

	table.Render()

	return nil
}

func init() {
	initConfig()
	statusCmd := StatusCmd(NewGitlabAPIClient())
	rootCmd.AddCommand(statusCmd)
}
