package main

import (
	"encoding/json"
	"log"
	"time"
	"github.com/codegp/game-runner/gamestate"
)

func reportMatchHistory(gameStateUtils *gamestate.GameStateUtils, wc *WinCondition) {

	game.WinningTeam = wc.WinningTeam
	game.Reason = wc.Reason
	game.Complete = true
	game.Finished = time.Now()

	err := cp.UpdateGame(game)
	if err != nil {
		log.Fatalf("Unable to update game entity: %v", err)
	}

	content, err := json.Marshal(gameStateUtils.History())
	if err != nil {
		log.Fatalf("Unable to marshal history into json: %v", err)
	}

	err = cp.WriteHistory(game.ID, content)
	if err != nil {
		log.Fatalf("Unable to write history: %v", err)
	}
}
