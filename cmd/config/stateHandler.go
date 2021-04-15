package config

import (
	"encoding/json"
	"io/ioutil"
)

type State struct {
	CurrentProject  string `json:"current_project"`
	SSH_url_to_repo string `json:"ssh_url_to_repo"`
}

func convertJsonToState(jsonData []byte) (*State, error) {
	var state State
	err := json.Unmarshal(jsonData, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func convertStateToJson(state State) ([]byte, error) {
	json, err := json.Marshal(state)
	if err != nil {
		return []byte{}, err
	}
	return json, nil
}

func ReadState(file string) (*State, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	state, err := convertJsonToState(content)
	return state, err
}

func WriteState(state State, file string) error {
	json, err := convertStateToJson(state)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, json, 0644)
	if err != nil {
		return err
	}

	return nil
}
