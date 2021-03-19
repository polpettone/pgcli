package cmd

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/models"
	"github.com/spf13/cobra"
)

func NewPipelinesCmd(apiClient *GitlabAPIClient) *cobra.Command{
	return &cobra.Command{
		Use:   "pipelines",
		Short: "shows the last 5 Pipelines",
		Long:  "shows the last 5 Pipelines",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handlePipelineCommand(cmd, apiClient)
			apiClient.Logging.InfoLog.Printf("Pipeline Command Called")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}

func handlePipelineCommand(cobraCommand *cobra.Command, apiClient *GitlabAPIClient) (string, error) {
	status, _  := cobraCommand.Flags().GetString("status")
	count, _ := cobraCommand.Flags().GetInt("count")
	withUser, _ := cobraCommand.Flags().GetBool("user")
	withCommitTitle, _ := cobraCommand.Flags().GetBool("commit-title")

	pipelines, err := apiClient.getPipelines(status, count)

	var enrichedPipelines []*models.Pipeline

	if withUser {
		enrichedPipelines, err = apiClient.enrichPipelinesByUser(pipelines, 10)
		pipelines = enrichedPipelines
	}

	if withCommitTitle {
		enrichedPipelines, err = apiClient.enrichPipelinesByJobs(pipelines, 10)
		pipelines = enrichedPipelines
	}

	if err != nil {
		return "", err
	}
	value := ""
	for _, p := range pipelines {
		value = value + "\n" + p.NiceString()
	}
	return value, nil
}


func init() {
	pipelinesCmd := NewPipelinesCmd(NewGitlabAPIClient())

	pipelinesCmd.Flags().StringP(
		"status",
			  "s",
			  "",
			  "filter pipelines by status: " +
			  "running, pending, success, failed, canceled, skipped, created, manual",
		)

	pipelinesCmd.Flags().BoolP(
		"user",
		"u",
		false,
		"shows user which triggered the pipeline." +
			"Takes longer due more api calls (each per pipeline)",
	)

	pipelinesCmd.Flags().IntP(
		"count",
		"c",
		5,
		"Count of Pipelines returned. 100 is max" ,
	)

	pipelinesCmd.Flags().BoolP(
		"commit-title",
		"t",
		false,
		"Show commit title for each pipeline. Take longer due more api calls" ,
	)

	rootCmd.AddCommand(pipelinesCmd)
}
