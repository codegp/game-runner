package gamestate

import (
	"fmt"

	"github.com/codegp/cloud-persister"
	"github.com/codegp/game-object-types/types"
	"github.com/codegp/game-runner/gameinfo"
	gobj "github.com/codegp/game-runner/gameobjects"
)

type GameStateUtils struct {
	gs       *gameState
	gameInfo *gameinfo.GameInfo
}

func NewGameStateUtils(cp *cloudpersister.CloudPersister, gameInfo *gameinfo.GameInfo) *GameStateUtils {
	gs := newGameState()
	gsu := &GameStateUtils{
		gs,
		gameInfo,
	}
	gsu.parseMap(cp)
	return gsu
}

// newInvalidMovef creates an error
func newInvalidMovef(msg string, args ...interface{}) *gobj.InvalidMove {
	err := gobj.NewInvalidMove()
	err.Message = fmt.Sprintf(msg, args...)
	return err
}

/*
*  Getters
 */

func (u *GameStateUtils) GameInfo() *gameinfo.GameInfo {
	return u.gameInfo
}

func (u *GameStateUtils) Round() int {
	return u.gs.round
}

func (u *GameStateUtils) CurrentBot() *gobj.Bot {
	return u.gs.currentBot
}

func (u *GameStateUtils) Bots() []*gobj.Bot {
	bots := make([]*gobj.Bot, len(u.gs.bots))
	i := 0
	for _, b := range u.gs.bots {
		bots[i] = b
		i++
	}
	return bots
}

func (u *GameStateUtils) Bot(id int32) *gobj.Bot {
	return u.gs.bots[id]
}

func (u *GameStateUtils) Item(id int32) *gobj.Item {
	return u.gs.items[id]
}

func (u *GameStateUtils) Map() gobj.Map {
	return u.gs.currentMap
}

func (u *GameStateUtils) History() *gobj.History {
	return u.gs.history
}

func (u *GameStateUtils) BotsToCreate() []*gobj.Bot {
	return u.gs.botsToCreate
}

func (u *GameStateUtils) BotsToDestroy() []int32 {
	return u.gs.botsToDestroy
}

func (u *GameStateUtils) LocationInfoAtLocation(loc *gobj.Location) (*gobj.LocationInfo, error) {
	if !u.OnTheMap(loc) {
		return nil, fmt.Errorf("Location is off the map")
	}
	return u.Map()[loc.X][loc.Y], nil
}

func (u *GameStateUtils) OnTheMap(loc *gobj.Location) bool {
	m := u.Map()
	if loc.X < 0 ||
		loc.Y < 0 ||
		loc.X >= int32(len(m)) ||
		loc.Y >= int32(len(m[0])) {
		return false
	}
	return true
}

/*
* 	GameState manipulation
 */

func (u *GameStateUtils) IncrementRound() {
	u.gs.round++
}

func (u *GameStateUtils) SetCurrentBot(bot *gobj.Bot) {
	u.gs.currentBot = bot
}

func (u *GameStateUtils) AppendMapToHistory() {
	u.gs.history.Maps = append(u.gs.history.Maps, copyMap(u.gs.currentMap))
}

func (u *GameStateUtils) TakeHealthFromBotAtLocation(loc *gobj.Location, damage float64) *gobj.InvalidMove {
	locInfo, err := u.LocationInfoAtLocation(loc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}

	if !u.IsLocationInfoOccupied(locInfo) {
		return newInvalidMovef("Location %v is not occupied", locInfo.Loc)
	}

	locInfo.Bot.Health -= damage
	if locInfo.Bot.Health <= 0 {
		u.gs.botsToDestroy = append(u.gs.botsToDestroy, locInfo.Bot.ID)
		delete(u.gs.bots, locInfo.Bot.ID)
		locInfo.Bot = nil
	}
	return nil
}

func (u *GameStateUtils) InitBot(teamID int32, loc *gobj.Location, botType *types.BotType) (*gobj.Bot, *gobj.InvalidMove) {
	locInfo, err := u.LocationInfoAtLocation(loc)
	if err != nil {
		return nil, newInvalidMovef(err.Error())
	}

	if u.IsLocationInfoOccupied(locInfo) {
		return nil, newInvalidMovef("Location %v is occupied", locInfo.Loc)
	}

	u.gs.maxBotID++
	bot := &gobj.Bot{
		ID:          u.gs.maxBotID,
		Loc:         loc,
		Health:      botType.MaxHealth,
		AttackDelay: 0,
		MoveDelay:   0,
		TeamID:      teamID,
		BotTypeID:   botType.ID,
		Items:       []*gobj.Item{},
	}

	u.gs.bots[u.gs.maxBotID] = bot
	locInfo.Bot = bot
	u.PickUpItemAtLoc(locInfo.Loc)

	u.gs.botsToCreate = append(u.gs.botsToCreate, bot)

	return bot, nil
}

func (u *GameStateUtils) InitItem(loc *gobj.Location, itemType *types.ItemType) (*gobj.Item, *gobj.InvalidMove) {
	locInfo, err := u.LocationInfoAtLocation(loc)
	if err != nil {
		return nil, newInvalidMovef(err.Error())
	}

	if u.LocationInfoContainsItem(locInfo) {
		return nil, newInvalidMovef("Location %v has item already", locInfo.Loc)
	}

	u.gs.maxItemID++
	item := &gobj.Item{
		ItemTypeID: itemType.ID,
		ID:         u.gs.maxBotID,
	}

	u.gs.items[u.gs.maxItemID] = item
	locInfo.Item = item

	return item, nil
}

func (u *GameStateUtils) MoveBotToLocation(bot *gobj.Bot, loc *gobj.Location) *gobj.InvalidMove {
	botLocInfo, err := u.LocationInfoAtLocation(bot.Loc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}

	locInfo, err := u.LocationInfoAtLocation(loc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}

	if u.IsLocationInfoOccupied(locInfo) {
		return newInvalidMovef("Location %v is occupied", locInfo.Loc)
	}

	botLocInfo.Bot = nil
	locInfo.Bot = bot
	bot.Loc = locInfo.Loc
	u.PickUpItemAtLoc(locInfo.Loc)

	return nil
}

func (u *GameStateUtils) MoveItemToLocation(item *gobj.Item, loc *gobj.Location) *gobj.InvalidMove {
	// TODO: find the current location of the item and remove the reference
	// itemLocaInfo, err := u.LocationInfoAtLocation(item.Loc)
	// if err != nil {
	// 	return newInvalidMovef(err.Error())
	// }

	locInfo, err := u.LocationInfoAtLocation(loc)
	if err != nil {
		return newInvalidMovef(err.Error())
	}

	if u.LocationInfoContainsItem(locInfo) {
		return newInvalidMovef("Location %v is occupied", locInfo.Loc)
	}

	// itemLocaInfo.Item = nil
	locInfo.Item = item
	// item.Loc = locInfo.Loc
	u.PickUpItemAtLoc(locInfo.Loc)

	return nil
}

func (u *GameStateUtils) PickUpItemAtLoc(loc *gobj.Location) {
	locInfo, err := u.LocationInfoAtLocation(loc)
	if err == nil && locInfo.Bot != nil && locInfo.Item != nil {
		locInfo.Bot.Items = append(locInfo.Bot.Items, locInfo.Item)
		locInfo.Item = nil
	}
}

func (u *GameStateUtils) ClearPendingBotActions() {
	u.gs.botsToCreate = []*gobj.Bot{}
	u.gs.botsToDestroy = []int32{}
}
