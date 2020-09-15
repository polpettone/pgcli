package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func NewLogsCmd(apiClient APIClient) *cobra.Command {
	return &cobra.Command{
		Use:
		"logs -> interactive mode| logs <jobID> -> logs of job | logs -l -> logs of last failed job",
		Short: "when no job id " +
			"is given, interactive mode started to choose a pipeline, then a job to see logs or use flag -l to see log of the last failed job",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleLogsCommand(cmd, args, apiClient)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(stdout)
		},
	}
}

func handleLogsCommand(cobraCommand *cobra.Command, args []string, apiClient APIClient) (string, error) {
	lastFailed, _ := cobraCommand.Flags().GetBool("lastFailed")

	var jobID string

	if (len(args) < 1 || args[0] == "") && lastFailed == false {

		pipelines, _ := apiClient.getPipelines("", false, 20)
		pipeline, _ := showPipelineSelectionPrompt(pipelines)
		pipelineId := strconv.Itoa(pipeline.Id)

		jobs, _ := apiClient.getJobs(pipelineId)
		job, _ := showJobSelectionPrompt(jobs)
		jobID = strconv.Itoa(job.Id)

		return apiClient.getLog(jobID)
	} else {

		if lastFailed {
			return apiClient.getLastFailLog()
		} else {
			jobID = args[0]
			return apiClient.getLog(jobID)
		}
	}
}

func init() {
	logsCmd := NewLogsCmd(gitlabAPIClient)
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolP(
		"lastFailed",
		"l",
		false,
		"Shows the logs from the last failed job")
}
