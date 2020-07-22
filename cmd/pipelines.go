package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewPipelinesCmd(apiClient APIClient) *cobra.Command{
	return &cobra.Command{
		Use:   "pipelines",
		Short: "shows the last 20 Pipelines",
		Long:  "shows the last 20 Pipelines",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handlePipelineCommand(cmd, apiClient)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), stdout)
		},
	}
}

func handlePipelineCommand(cobraCommand *cobra.Command, apiClient APIClient) (string, error) {
	status, _  := cobraCommand.Flags().GetString("status")
	pipelines, err := apiClient.getPipelines(status)
	if err != nil {
		return "", err
	}
	value := ""
	for _, p := range pipelines {
		value = value + "\n" + p.niceString()
	}
	return value, nil
}


func init() {
	pipelinesCmd := NewPipelinesCmd(gitlabAPIClient)
	pipelinesCmd.Flags().StringP(
		"status",
			  "s",
			  "",
			  "filter pipelines by status: " +
			  "running, pending, success, failed, canceled, skipped, created, manual",
		)
	rootCmd.AddCommand(pipelinesCmd)

}
