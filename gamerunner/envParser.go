package main

import (
	"log"
	"os"
	"strconv"

	"github.com/codegp/cloud-persister/models"
)

func fetchGameEntities() {
	var err error
	gameIDStr := os.Getenv("GAME_ID")
	if gameIDStr == "" {
		log.Fatalf("No gameID provided")
	}
	ID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		log.Fatalf("failed to get game: Invalid game id, %v, must be integer", gameIDStr)
	}
	game, err = cp.GetGame(ID)
	if err != nil {
		log.Fatalf("Error getting game from datastore, %v", err)
	}

	var proj *models.Project
	for _, projID := range game.ProjectIDs {
		proj, err = cp.GetProject(projID)
		if err != nil {
			log.Fatalf("failed to get projects: %v", err)
		}

		projects = append(projects, proj)
	}
}

func getShouldReport() bool {
	reportResults := os.Getenv("REPORT_RESULTS")
	if reportResults == "" {
		return true
	}
	shouldReport, err := strconv.ParseBool(reportResults)
	if err != nil {
		return false
	}
	return shouldReport
}

func getIP() string {
	ip := os.Getenv("POD_IP")
	if ip == "" {
		log.Fatal("Pod id env var not set")
	}
	return ip
}
