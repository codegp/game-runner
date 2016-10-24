package gamedefutils

import (
	"fmt"

	"github.com/codegp/game-object-types/types"
	gobj "github.com/codegp/game-runner/gameobjects"
	"github.com/codegp/game-runner/gamestate"
)

// GameDefUtils is the object passed to game type definition to manipulate the game state
type GameDefUtils struct {
	gsu *gamestate.GameStateUtils
}

// NewGameDefUtils returns an instance of GameDefUtils
func NewGameDefUtils(gsu *gamestate.GameStateUtils) *GameDefUtils {
	return &GameDefUtils{
		gsu,
	}
}

// Bots returns a list of all active bots
func (u *GameDefUtils) Bots() []gobj.Bot {
	botPointers := u.gsu.Bots()
	bots := make([]gobj.Bot, len(botPointers))
	for i, ptr := range botPointers {
		bots[i] = *ptr
	}
	return bots
}

// Items returns a list of all active items
func (u *GameDefUtils) Items() []gobj.Item {
	itemPointers := u.gsu.Items()
	items := make([]gobj.Item, len(itemPointers))
	for i, ptr := range itemPointers {
		items[i] = *ptr
	}
	return items
}

// InitBot initializes a new bot of type botTypeID for the team teamID at the location loc.
// returns an error if botTypeID is not a valid bot type identifier, if location loc
// does not exist, or if location loc is already occupied by a bot
func (u *GameDefUtils) InitBot(teamID int32, loc gobj.Location, botTypeID int64) (gobj.Bot, error) {
	botType := u.gsu.GameInfo().BotType(botTypeID)
	if botType == nil {
		return gobj.Bot{}, fmt.Errorf("botTypeID %d is not a valid id", botTypeID)
	}
	bot, err := u.gsu.InitBot(teamID, &loc, botType)
	if err != nil {
		return gobj.Bot{}, fmt.Errorf("Failed no init bot, err:\n%v", err.Message)
	}
	return *bot, nil
}

// InitItem initializes a new item of type itemTypeID at the location loc.
// returns an error if itemTypeID is not a valid item type identifier, if location loc
// does not exist, or if location loc is already occupied by an item
func (u *GameDefUtils) InitItem(loc gobj.Location, itemTypeID int64) (gobj.Item, error) {
	itemType := u.gsu.GameInfo().ItemType(itemTypeID)
	if itemType == nil {
		return gobj.Item{}, fmt.Errorf("itemTypeID %d is not a valid id", itemTypeID)
	}
	item, err := u.gsu.InitItem(&loc, itemType)
	if err != nil {
		return gobj.Item{}, fmt.Errorf("Failed no move bot, err:\n%v", err.Message)
	}

	return *item, nil
}

// MoveBotToLocation moves the bot with the id botID to the location loc
// returns an error if there is not bot with the id botID, if location loc
// does not exist, or if location loc is already occupied by a bot
func (u *GameDefUtils) MoveBotToLocation(botID int32, loc gobj.Location) error {
	bot := u.gsu.Bot(botID)
	if bot == nil {
		return fmt.Errorf("botID %d is not valid id", botID)
	}
	err := u.gsu.MoveBotToLocation(bot, &loc)
	if err != nil {
		return fmt.Errorf("Failed no move bot, err:\n%v", err.Message)
	}
	return nil
}

// MoveItemToLocation moves the item with the id itemID to the location loc
// returns an error if there is not item with the id itemID, if location loc
// does not exist, or if location loc is already occupied by a item
func (u *GameDefUtils) MoveItemToLocation(itemID int32, loc gobj.Location) error {
	item := u.gsu.Item(itemID)
	if item == nil {
		return fmt.Errorf("itemID %d is not valid id", itemID)
	}
	err := u.gsu.MoveItemToLocation(item, &loc)
	if err != nil {
		return fmt.Errorf("Failed no move bot, err:\n%v", err.Message)
	}
	return nil
}

// RemoveItem removes an item from the map or from possession of a bot
func (u *GameDefUtils) RemoveItem(itemID int32) error {
	return u.gsu.RemoveItem(itemID)
}

// GiveHealthToBot increases the health of the bot with id botID by health
// or decreases if the value of healths is negative
func (u *GameDefUtils) GiveHealthToBot(botID int32, health float64) error {
	bot := u.gsu.Bot(botID)
	if bot == nil {
		return fmt.Errorf("botID %d is not valid id", botID)
	}

	err := u.gsu.TakeHealthFromBotAtLocation(bot.Loc, health*-1)
	if err != nil {
		return fmt.Errorf("Failed no move bot, err:\n%v", err.Message)
	}

	return nil
}

// GiveSpawnDelayToBot -
func (u *GameDefUtils) GiveSpawnDelayToBot(botID int32, delay float64) error {
	bot := u.gsu.Bot(botID)
	if bot == nil {
		return fmt.Errorf("botID %d is not valid id", botID)
	}

	bot.SpawnDelay += delay
	return nil
}

// GiveMoveDelayToBot -
func (u *GameDefUtils) GiveMoveDelayToBot(botID int32, delay float64) error {
	bot := u.gsu.Bot(botID)
	if bot == nil {
		return fmt.Errorf("botID %d is not valid id", botID)
	}

	bot.MoveDelay += delay
	return nil
}

// GiveAttackDelayToBot -
func (u *GameDefUtils) GiveAttackDelayToBot(botID int32, delay float64) error {
	bot := u.gsu.Bot(botID)
	if bot == nil {
		return fmt.Errorf("botID %d is not valid id", botID)
	}

	bot.AttackDelay += delay
	return nil
}

// BotAtLocation -
func (u *GameDefUtils) BotAtLocation(loc gobj.Location) (gobj.Bot, error) {
	return gobj.Bot{}, nil
}

// ItemAtLocation -
func (u *GameDefUtils) ItemAtLocation(loc gobj.Location) (gobj.Item, error) {
	return gobj.Item{}, nil
}

// TerrainAtLocation -
func (u *GameDefUtils) TerrainAtLocation(loc gobj.Location) (types.TerrainType, error) {
	return types.TerrainType{}, nil
}

// ChangeTerrainAtLocation -
func (u *GameDefUtils) ChangeTerrainAtLocation(terrainType types.TerrainType, loc gobj.Location) error {
	return nil
}

// MapDimensions -
func (u *GameDefUtils) MapDimensions() (int, int, error) {
	return 0, 0, nil
}
