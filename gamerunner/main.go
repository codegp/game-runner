package main

import (
	"log"

	"github.com/codegp/cloud-persister"
	"github.com/codegp/cloud-persister/models"
	"github.com/codegp/game-runner/gamestate"
)

var projects []*models.Project
var cp *cloudpersister.CloudPersister
var game *models.Game

func init() {
	var err error
	cp, err = cloudpersister.NewCloudPersister()
	if err != nil {
		log.Fatalf("Failed to start cloud persister: %v", err)
	}

	fetchGameEntities()
}

func main() {
	gameInfo, err := gamestate.NewGameInfo(cp, game)
	if err != nil {
		log.Fatal(err.Error())
	}

	gameStateUtils := gamestate.NewGameStateUtils(cp, gameInfo, projects)

	// start the api server
	server := startServer(gameStateUtils)

	gr := newGameRunner(gameStateUtils, newTurnInformerManager(projects))
	wc, err := gr.doGame()
	if err != nil {
		log.Fatal(err.Error())
	}
	// stop the server
	server.Stop()

	if wc.WinningTeam >= 0 {
		log.Printf("Game Over!\nTeam %d wins!\nReason: %s", wc.WinningTeam, wc.Reason)
	} else {
		log.Printf("Game Over!\nDraw!\nReason: %s", wc.Reason)
	}

	reportMatchHistory(gameStateUtils, wc)
}
