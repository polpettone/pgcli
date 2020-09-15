package cmd

import (
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
}

func NewApplication() *Application {

	infoLog := log.New(openLogFile("info.log"), "INFO\t", log.Ldate|log.Ltime)
	debugLog := log.New(openLogFile("debug.log"), "DEBUG\t", log.Ldate|log.Ltime)
	errorLog := log.New(openLogFile("error.log"),"ERROR\t", log.Ldate|log.Ltime)

	app := &Application{
		errorLog: errorLog,
		infoLog: infoLog,
		debugLog: debugLog,
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
