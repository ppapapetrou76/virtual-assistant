package main

import (
	"io/ioutil"
	"log"
	"os"

	labeler "github.com/ppapapetrou76/virtual-assistant/pkg/actions"
	"github.com/ppapapetrou76/virtual-assistant/pkg/config"
	"github.com/ppapapetrou76/virtual-assistant/pkg/github"
)

const errGeneral = "Unable to execute action: %+v"

func main() {
	eventPayload := getEventPayload()
	eventName := os.Getenv(github.EventNameEnvVar)

	repo := github.NewRepo()
	cfgRaw, err := repo.LoadFile(os.Getenv(github.InputConfigPathEnvVar),
		os.Getenv(github.ShaEnvVar))
	checkErr(err)

	cfg, err := config.Load(cfgRaw)
	checkErr(err)

	log.Printf("Re-evaluating labels on %s@%s",
		os.Getenv(github.RepoEnvVar),
		os.Getenv(github.ShaEnvVar))

	log.Printf("Trigger event: %s", os.Getenv(github.EventNameEnvVar))

	err = labeler.New(cfg, repo).HandleEvent(eventName, eventPayload)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf(errGeneral, err)
	}
}

func getEventPayload() *[]byte {
	payloadPath := os.Getenv(github.EventPathEnvVar)
	file, err := os.Open(payloadPath)
	if err != nil {
		log.Fatalf("Failed to open event payload file %s: %s", payloadPath, err)
	}
	eventPayload, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to load event payload from %s: %s", payloadPath, err)
	}
	return &eventPayload
}
