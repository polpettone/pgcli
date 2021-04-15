package commands

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/adapter"
	"github.com/spf13/cobra"
	"strconv"
)

func NewJobsCmd(apiClient *adapter.App) *cobra.Command{
	return &cobra.Command{
		Use:   "jobs <pipelineID>",
		Short: "list the jobs of a specific pipeline",
		Long:  "list the jobs of a specific pipeline, with start time, end time, duration and state",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout, err := handleJobsCommand(args, apiClient)
			if err != nil {
				return err
			}
			fmt.Println(stdout)
			return nil
		},

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return getPipelineSuggestions(apiClient), cobra.ShellCompDirectiveNoFileComp
		},
	}
}

func getPipelineSuggestions(apiClient *adapter.App) []string {
	pipelines, _ := apiClient.GetPipelines("",  20)
	var pipelineIds []string
	for _, p := range pipelines[:5] {
		pipelineIds = append(pipelineIds, strconv.Itoa(p.Id))
	}
	return pipelineIds
}

func handleJobsCommand(args []string, apiClient *adapter.App) (string, error) {

	var pipelineId string

	if len(args) < 1 || args[0] == "" {
		pipelines, _ := apiClient.GetPipelines("",  20)
		pipeline , _ :=  adapter.ShowPipelineSelectionPrompt(pipelines)
		pipelineId = strconv.Itoa(pipeline.Id)
	} else {
		pipelineId = args[0]
	}


	jobs, err := apiClient.GetJobs(pipelineId)
	if err != nil {
		return "", err
	}

	value := ""
	for _, job := range jobs {
		value = value + "\n" + job.NiceString()
	}

	return value, nil
}

func init() {
	InitConfig()
	jobsCmd := NewJobsCmd(adapter.NewApp())
	rootCmd.AddCommand(jobsCmd)
}
