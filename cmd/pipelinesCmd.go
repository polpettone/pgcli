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
	withUser, _ := cobraCommand.Flags().GetBool("user")
	count, _ := cobraCommand.Flags().GetInt("count")
	pipelines, err := apiClient.getPipelines(status, withUser, count)
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
	pipelinesCmd := NewPipelinesCmd(gitlabAPIClient)

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
		20,
		"Count of Pipelines returned. 100 is max" ,
	)

	rootCmd.AddCommand(pipelinesCmd)

}
