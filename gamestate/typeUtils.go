package gamestate

import (
	"fmt"
	"math"

	"github.com/codegp/game-object-types/types"
	gobj "github.com/codegp/game-runner/gameobjects"
)

/*
*   BotType Utils
 */

func (u *GameStateUtils) BotHasMoveDelay(bot *gobj.Bot) bool {
	return bot.MoveDelay >= 1
}

func (u *GameStateUtils) BotHasAttackDelay(bot *gobj.Bot) bool {
	return bot.AttackDelay >= 1
}

func (u *GameStateUtils) BotHasSpawnDelay(bot *gobj.Bot) bool {
	return bot.SpawnDelay >= 1
}

func (u *GameStateUtils) BotCanPerformMoveType(bot *gobj.Bot, moveType *types.MoveType) bool {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	for _, k := range botType.MoveTypeIDs {
		if k == moveType.ID {
			return true
		}
	}
	return false
}

func (u *GameStateUtils) BotCanPerformAttackType(bot *gobj.Bot, attackType *types.AttackType) bool {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	for _, k := range botType.AttackTypeIDs {
		if k == attackType.ID {
			return true
		}
	}
	return false
}

func (u *GameStateUtils) BotIsHoldingItemOfType(bot *gobj.Bot, itemType *types.ItemType) bool {
	for _, i := range bot.Items {
		if i.ItemTypeID == itemType.ID {
			return true
		}
	}
	return false
}

func (u *GameStateUtils) BotCanSpawnBotOfType(bot *gobj.Bot, toSpawnType *types.BotType) (bool, string) {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	if !botType.CanSpawn {
		return false, fmt.Sprintf("%s bots can not spawn other bots", botType.Name)
	}

	if !toSpawnType.CanBeSpawned {
		return false, fmt.Sprintf("%s bots can not be spawned", toSpawnType.Name)
	}
	return true, ""
}

func (u *GameStateUtils) SpawnDelayForBotType(bot *gobj.Bot, toSpawnType *types.BotType) float64 {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	return toSpawnType.SpawnDelay * botType.SpawnDelayFactor
}

func (u *GameStateUtils) MoveDelayForMoveType(bot *gobj.Bot, moveType *types.MoveType) float64 {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	locInfo, _ := u.LocationInfoAtLocation(bot.Loc)
	terrainType := u.GameInfo().TerrainType(locInfo.Terrain)
	return moveType.Delay * botType.MoveDelayFactor * terrainType.MoveDelayFactor
}

func (u *GameStateUtils) AttackRangeWithAttackType(bot *gobj.Bot, attackType *types.AttackType) float64 {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	return botType.RangeFactor * 1 // TODO: attackType.Range
}

func (u *GameStateUtils) AttackDamageWithAttackType(bot *gobj.Bot, attackType *types.AttackType) float64 {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	return botType.DamageFactor * attackType.Damage
}

func (u *GameStateUtils) AttackAccuracyWithAttackType(bot *gobj.Bot, attackType *types.AttackType) float64 {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	return botType.AccuracyFactor * attackType.Accuracy
}

func (u *GameStateUtils) AttackDelayWithAttackType(bot *gobj.Bot, attackType *types.AttackType) float64 {
	botType := u.GameInfo().BotType(bot.BotTypeID)
	return botType.AttackDelayFactor * attackType.Delay
}

/*
*   LocationInfo Utils
 */

func (u *GameStateUtils) IsLocationInfoOccupied(locInfo *gobj.LocationInfo) bool {
	return locInfo.Bot != nil
}

func (u *GameStateUtils) LocationInfoContainsItem(locInfo *gobj.LocationInfo) bool {
	return locInfo.Item != nil
}

func (u *GameStateUtils) LocationInfoContainsItemOfType(locInfo *gobj.LocationInfo, itemType *types.ItemType) bool {
	return locInfo.Item != nil && locInfo.Item.ItemTypeID == itemType.ID
}

/*
*   Location Utils
 */

func (u *GameStateUtils) DistanceSquaredTo(loc1 *gobj.Location, loc2 *gobj.Location) float64 {
	return math.Pow(float64(loc1.X-loc2.X), 2) + math.Pow(float64(loc1.Y-loc2.Y), 2)
}

func (u *GameStateUtils) LocationInAttackRange(loc *gobj.Location, bot *gobj.Bot, attackType *types.AttackType) bool {
	return u.DistanceSquaredTo(loc, bot.Loc) <= u.AttackRangeWithAttackType(bot, attackType)
}

func (u *GameStateUtils) AddDirectionToLocation(loc *gobj.Location, dir gobj.Direction) *gobj.Location {
	newX := loc.X
	newY := loc.Y

	switch dir {
	case gobj.Direction_NORTH:
		newY = loc.Y - 1
		break
	case gobj.Direction_NORTH_EAST:
		newY = loc.Y - 1
		newX = loc.X + 1
		break
	case gobj.Direction_EAST:
		newX = loc.X + 1
		break
	case gobj.Direction_SOUTH_EAST:
		newY = loc.Y + 1
		newX = loc.X + 1
		break
	case gobj.Direction_SOUTH:
		newY = loc.Y + 1
		break
	case gobj.Direction_SOUTH_WEST:
		newY = loc.Y + 1
		newX = loc.X - 1
		break
	case gobj.Direction_WEST:
		newX = loc.X - 1
		break
	case gobj.Direction_NORTH_WEST:
		newY = loc.Y - 1
		newX = loc.X - 1
		break
	}

	return &gobj.Location{
		X: newX,
		Y: newY,
	}
}
