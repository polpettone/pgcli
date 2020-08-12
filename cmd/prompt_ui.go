package cmd

import (
	"github.com/manifoldco/promptui"
)

func showPipelineSelectionPrompt(pipelines []GitlabPipeline) (*GitlabPipeline, error) {

	templates := &promptui.SelectTemplates{
		Active:   "\U00001B61 {{ .Id | red }} {{ .CreatedAt | green }} {{ .UpdatedAt | green }} {{ .Status | green }}",
		Inactive:  "{{ .Id | blue }} {{ .CreatedAt | blue }} {{ .UpdatedAt | blue }} {{ .Status | blue }}",
		Selected:   "\U00001B61 {{ .Id | green }} {{ .CreatedAt | green }} {{ .UpdatedAt | green }} {{ .Status | green }}",
	}

	prompt := promptui.Select{
		Label:     "Which Pipeline",
		Items:     pipelines,
		Templates: templates,
		Size:      20,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return &pipelines[i], nil
}

