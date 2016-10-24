package gamestate

import (
	gobj "github.com/codegp/game-runner/gameobjects"
)

/*
*		Actions and Action validators
*   These functions are invoked by the api handler
 */

func (u *GameStateUtils) ValidateSpawnAttempt(dir gobj.Direction, toSpawnTypeID int64) *gobj.InvalidMove {
	toSpawnType := u.GameInfo().BotType(toSpawnTypeID)
	if toSpawnType == nil {
		return newInvalidMovef("Invalid toSpawnTypeID provided, %d", toSpawnTypeID)
	}

	currentBot := u.CurrentBot()
	if u.BotHasSpawnDelay(currentBot) {
		return newInvalidMovef("Bot has spawn delay: %v", currentBot.SpawnDelay)
	}
	if canSpawn, reason := u.BotCanSpawnBotOfType(currentBot, toSpawnType); !canSpawn {
		return newInvalidMovef("%s", reason)
	}
	spawnLoc := u.AddDirectionToLocation(currentBot.Loc, dir)
	spawnLocInfo, err := u.LocationInfoAtLocation(spawnLoc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}
	if u.IsLocationInfoOccupied(spawnLocInfo) {
		return newInvalidMovef("Location is occupied")
	}
	return nil
}

func (u *GameStateUtils) SpawnInDirection(dir gobj.Direction, toSpawnTypeID int64) (*gobj.Bot, *gobj.InvalidMove) {
	toSpawnType := u.GameInfo().BotType(toSpawnTypeID)
	if toSpawnType == nil {
		return nil, newInvalidMovef("Invalid toSpawnTypeID provided, %d", toSpawnTypeID)
	}

	currentBot := u.CurrentBot()
	defer func() {
		currentBot.SpawnDelay += u.SpawnDelayForBotType(currentBot, toSpawnType)
	}()

	if err := u.ValidateSpawnAttempt(dir, toSpawnTypeID); err != nil {
		return nil, err
	}
	spawnLoc := u.AddDirectionToLocation(currentBot.Loc, dir)
	return u.InitBot(currentBot.TeamID, spawnLoc, toSpawnType)
}

func (u *GameStateUtils) ValidateMoveAttempt(dir gobj.Direction, moveTypeID int64) *gobj.InvalidMove {
	moveType := u.GameInfo().MoveType(moveTypeID)
	if moveType == nil {
		return newInvalidMovef("Invalid moveTypeID provided, %d", moveTypeID)
	}

	currentBot := u.CurrentBot()
	if u.BotHasMoveDelay(currentBot) {
		return newInvalidMovef("Bot has move delay: %v", currentBot.MoveDelay)
	}
	if canUseMoveType := u.BotCanPerformMoveType(currentBot, moveType); !canUseMoveType {
		return newInvalidMovef("Bot cannot use move type %s", moveType.Name)
	}
	moveLoc := u.AddDirectionToLocation(currentBot.Loc, dir)
	moveLocInfo, err := u.LocationInfoAtLocation(moveLoc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}
	if u.IsLocationInfoOccupied(moveLocInfo) {
		return newInvalidMovef("Location is occupied")
	}
	return nil
}

func (u *GameStateUtils) MoveInDirection(dir gobj.Direction, moveTypeID int64) *gobj.InvalidMove {
	moveType := u.GameInfo().MoveType(moveTypeID)
	if moveType == nil {
		return newInvalidMovef("Invalid moveTypeID provided, %d", moveTypeID)
	}
	currentBot := u.CurrentBot()
	defer func() {
		currentBot.MoveDelay += u.MoveDelayForMoveType(currentBot, moveType)
	}()

	if err := u.ValidateMoveAttempt(dir, moveTypeID); err != nil {
		return err
	}
	dest := u.AddDirectionToLocation(currentBot.Loc, dir)
	return u.MoveBotToLocation(currentBot, dest)
}

func (u *GameStateUtils) ValidateAttackAttempt(loc *gobj.Location, attackTypeID int64) *gobj.InvalidMove {
	attackType := u.GameInfo().AttackType(attackTypeID)
	if attackType == nil {
		return newInvalidMovef("Invalid attackType provided, %d", attackTypeID)
	}

	currentBot := u.CurrentBot()
	if u.BotHasAttackDelay(currentBot) {
		return newInvalidMovef("Bot has attack delay: %v", currentBot.AttackDelay)
	}
	if canUseAttackType := u.BotCanPerformAttackType(currentBot, attackType); !canUseAttackType {
		return newInvalidMovef("Bot cannot use attack type %s", attackType.Name)
	}
	if u.LocationInAttackRange(loc, currentBot, attackType) {
		return newInvalidMovef("Bot cannot use attack type %s, loc is out of range", attackType.Name)
	}
	attackLocInfo, err := u.LocationInfoAtLocation(loc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}
	if !u.IsLocationInfoOccupied(attackLocInfo) {
		return newInvalidMovef("Location is not occupied")
	}
	return nil
}

func (u *GameStateUtils) AttackLocation(loc *gobj.Location, attackTypeID int64) *gobj.InvalidMove {
	attackType := u.GameInfo().AttackType(attackTypeID)
	if attackType == nil {
		return newInvalidMovef("Invalid attackType provided, %d", attackTypeID)
	}

	currentBot := u.CurrentBot()
	defer func() {
		currentBot.AttackDelay += u.AttackDelayWithAttackType(currentBot, attackType)
	}()

	if err := u.ValidateAttackAttempt(loc, attackTypeID); err != nil {
		return err
	}
	if u.randomizer.randomPercent() > int(u.AttackAccuracyWithAttackType(currentBot, attackType)*100) {
		return nil
	}
	return u.TakeHealthFromBotAtLocation(loc, u.AttackDamageWithAttackType(currentBot, attackType))
}
