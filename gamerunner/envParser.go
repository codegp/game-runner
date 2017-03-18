package main

import (
	"log"
	"os"
	"strconv"
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

	projects, err = cp.GetMultiProject(game.ProjectIDs)
	if err != nil {
		log.Fatalf("Error getting projects from datastore, %v", err)
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
