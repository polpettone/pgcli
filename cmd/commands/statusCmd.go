package commands

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/polpettone/pgcli/cmd/adapter"
	"github.com/polpettone/pgcli/cmd/models"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func StatusCmd(apiClient *adapter.App) *cobra.Command {
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

func handleStatusCommand(args []string, apiClient *adapter.App) error {
	var pipelineId string
	if len(args) < 1 || args[0] == "" {
		pipelines, _ := apiClient.GetPipelines("", 10)
		pipeline, _ := adapter.ShowPipelineSelectionPrompt(pipelines)
		pipelineId = strconv.Itoa(pipeline.Id)
	} else {
		pipelineId = args[0]
	}

	jobs, err := apiClient.GetJobs(pipelineId)
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
	table.SetBorder(false)
	for _, d := range data {
		table.Append(d)
	}

	pipeline, err := apiClient.GetPipeline(pipelineId)

	if err != nil {
		return err
	}

	pipelines := []*models.Pipeline{pipeline}
	enrichedPipelines, err := apiClient.EnrichPipelinesByUser(pipelines, 1)

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
	InitConfig()
	statusCmd := StatusCmd(adapter.NewApp())
	rootCmd.AddCommand(statusCmd)
}
