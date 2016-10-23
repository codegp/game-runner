package main

import (
	"github.com/codegp/game-runner/gamedefutils"
	"github.com/codegp/game-runner/gamestate"
)

// GameType is an interface defining the games parameters and rules
type GameType interface {

	/*
	   WinCondition lets you define the parameters to end the game
	   A GameState is passed which gives you the means to decide if the game should
	   be over based on round number, bot count, or if other objectives have been met
	   @param GameState - the state of the game at the time we are checking for a win condition
	   @return boolean - true when the game is over
	   @return int - id of team that won
	*/
	WinCondition() (bool, *WinCondition)

	/*
	   Make any mutations to the state of the game world you would like to take place
	   between every round. This could be adding new bots, dealing additional damage,
	   dealing additional items, etc.
	   The round property of the gamestate will represent the round that just finished
	   @param GameState the gamestate at the end of the round
	*/
	BetweenRound()
}

// WinCondition represents if any team has won, and if so which team and why
type WinCondition struct {
	WinningTeam int64
	Reason      string
}

// GameDefinition is the gametypes implementation of GameType interface
type GameDefinition struct {
	gameDefUtils *gamedefutils.GameDefUtils
}

func newGameDefinition(gameStateUtils *gamestate.GameStateUtils) *GameDefinition {
	return &GameDefinition{
		gameDefUtils: gamedefutils.NewGameDefUtils(gameStateUtils),
	}
}
