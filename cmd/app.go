package cmd

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	initConfig()
}

type Application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	debugLog      *log.Logger
	AppClient     *GitlabAPIClient
}

func NewApplication() *Application {

	infoLog := log.New(openLogFile("info.log"), "INFO\t", log.Ldate|log.Ltime)
	debugLog := log.New(openLogFile("debug.log"), "DEBUG\t", log.Ldate|log.Ltime)
	errorLog := log.New(openLogFile("error.log"),"ERROR\t", log.Ldate|log.Ltime)

	gitlabAPIClient := NewGitlabAPIClient(
		viper.GetString("api_token"),
		viper.GetString("url"),
		viper.GetString("project_id"))

	app := &Application{
		errorLog: errorLog,
		infoLog: infoLog,
		debugLog: debugLog,
		AppClient: gitlabAPIClient,
	}

	return app
}

func openLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}
