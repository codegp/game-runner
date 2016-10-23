package main

import (
	"fmt"

	gobj "github.com/codegp/game-runner/gameobjects"
	"github.com/codegp/game-runner/gamestate"
)

// apiHandler implements the api interface defined by thrift
type apiHandler struct {
	gameStateUtils *gamestate.GameStateUtils
}

// NewapiHandler creates a new instance of the apihandler
func newAPIHandler(gameStateUtils *gamestate.GameStateUtils) *apiHandler {
	return &apiHandler{
		gameStateUtils: gameStateUtils,
	}
}

// Me returns the current Bot
func (a *apiHandler) Me() (*gobj.Bot, error) {
	bot := a.gameStateUtils.CurrentBot()
	return bot, nil
}

func (a *apiHandler) CanSpawn(dir gobj.Direction, botID int64) (bool, error) {
	return a.gameStateUtils.ValidateSpawnAttempt(dir, botID) == nil, nil
}

// Spawn cretes a new bot and adds it to the board
func (a *apiHandler) Spawn(dir gobj.Direction, botID int64) (*gobj.Bot, error) {
	return a.gameStateUtils.SpawnInDirection(dir, botID)
}

// CanMove returns true if the bot does not have move delay and the
// location in the input direction can be moved into
func (a *apiHandler) CanMove(dir gobj.Direction, moveID int64) (bool, error) {
	return a.gameStateUtils.ValidateMoveAttempt(dir, moveID) == nil, nil
}

// Move checks if the bot can move in the direction, if not throws an exception
// Otherwise move updtates the bots location and moveDelay and also updates the map
func (a *apiHandler) Move(dir gobj.Direction, moveID int64) error {
	return a.gameStateUtils.MoveInDirection(dir, moveID)
}

// CanAttack returns true if the bot does not have attack delay, the input location
// is on the map, and the location is within the bots attack range
func (a *apiHandler) CanAttack(loc *gobj.Location, attackID int64) (bool, error) {
	return a.gameStateUtils.ValidateAttackAttempt(loc, attackID) == nil, nil
}

// Attack checks if the bot can attack, if not raises an exception
// Otherwise deals damage and delay to bots and checks if attacked robot is dead
func (a *apiHandler) Attack(loc *gobj.Location, attackID int64) error {
	return a.gameStateUtils.AttackLocation(loc, attackID)
}

// BotAtLocation returns a Bot if there is one on the input location
// If the input location is not on the map raises and exception
// Otherwise returns nil
func (a *apiHandler) BotAtLocation(loc *gobj.Location) (*gobj.Bot, error) {
	locInfo, err := a.gameStateUtils.LocationInfoAtLocation(loc)
	if err != nil {
		return nil, newInvalidMoveError("Invalid location")
	}
	return locInfo.Bot, nil
}

// newInvalidMoveError creates an error
func newInvalidMoveError(msg string, args ...interface{}) *gobj.InvalidMove {
	err := gobj.NewInvalidMove()
	err.Message = fmt.Sprintf(msg, args...)
	return err
}
