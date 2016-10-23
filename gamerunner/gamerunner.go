package main

import (
	"log"
	"github.com/codegp/game-runner/gamestate"
)

type gameRunner struct {
	tiManager      *turnInformerManager
	gameStateUtils *gamestate.GameStateUtils
	gameDef        *GameDefinition
}

func newGameRunner(gameStateUtils *gamestate.GameStateUtils, tiManager *turnInformerManager) *gameRunner {
	gameDef := newGameDefinition(gameStateUtils)
	return &gameRunner{
		tiManager:      tiManager,
		gameStateUtils: gameStateUtils,
		gameDef:        gameDef,
	}
}

func (gr *gameRunner) checkForGameOver() (bool, *WinCondition) {
	if over, wc := gr.gameDef.WinCondition(); over {
		return over, wc
	}
	if gr.gameStateUtils.Round() >= 5 {
		return true, &WinCondition{
			WinningTeam: -1,
			Reason:      "Hit round limit",
		}
	}

	return false, nil
}

func (gr *gameRunner) doGame() (*WinCondition, error) {
	log.Println("Starting Game")
	gr.tiManager.createNewClients(gr.gameStateUtils.BotsToCreate())
	gr.gameStateUtils.ClearPendingBotActions()

	for {
		if err := gr.doRound(); err != nil {
			return nil, err
		}
		if gameOver, wc := gr.checkForGameOver(); gameOver {
			gr.tiManager.destroyAllClients()
			return wc, nil
		}
		// execute game type specific logic BetweenRound
		gr.gameDef.BetweenRound()
		gr.gameStateUtils.IncrementRound()
	}
}

func (gr *gameRunner) doRound() error {
	for _, bot := range gr.gameStateUtils.Bots() {
		gr.gameStateUtils.SetCurrentBot(bot)
		if err := gr.doTurn(); err != nil {
			return err
		}
	}

	gr.tiManager.createNewClients(gr.gameStateUtils.BotsToCreate())
	gr.tiManager.destroyClients(gr.gameStateUtils.BotsToDestroy())
	gr.gameStateUtils.ClearPendingBotActions()

	// printMap(CopyMap(gameState.CurrentMap))
	// perserve the state of the map after the round completes
	gr.gameStateUtils.AppendMapToHistory()
	log.Printf("Round %d complete\n", gr.gameStateUtils.Round())
	return nil
}

func (gr *gameRunner) doTurn() error {
	bot := gr.gameStateUtils.CurrentBot()
	err := gr.tiManager.startTurn(bot.ID)
	if err != nil {
		return err
	}

	if bot.MoveDelay > 0 {
		bot.MoveDelay = bot.MoveDelay - 1
		if bot.MoveDelay < 0 {
			bot.MoveDelay = 0
		}
	}
	if bot.AttackDelay > 0 {
		bot.AttackDelay = bot.AttackDelay - 1
		if bot.AttackDelay < 0 {
			bot.AttackDelay = 0
		}
	}
	if bot.SpawnDelay > 0 {
		bot.SpawnDelay = bot.SpawnDelay - 1
		if bot.SpawnDelay < 0 {
			bot.SpawnDelay = 0
		}
	}

	return nil
}

func (gr *gameRunner) printMap() {
	m := gr.gameStateUtils.Map()
	log.Printf("MAP: \n\n")
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			atLoc := m[i][j]
			if atLoc.Bot != nil {
				log.Printf("| %v |", atLoc.Bot.ID)
				continue
			}
			log.Printf("|   |")
		}
		log.Println()
	}
	log.Printf("\n\n")
}
