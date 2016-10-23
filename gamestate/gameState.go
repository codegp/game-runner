package gamestate

import (
	gobj "github.com/codegp/game-runner/gameobjects"
)

type gameState struct {
	round         int
	currentBot    *gobj.Bot
	bots          map[int32]*gobj.Bot
	items         map[int32]*gobj.Item
	currentMap    gobj.Map
	maxBotID      int32
	maxItemID     int32
	history       *gobj.History
	botsToCreate  []*gobj.Bot
	botsToDestroy []int32
}

func newGameState() *gameState {
	return &gameState{
		round:         0,
		currentBot:    nil,
		bots:          map[int32]*gobj.Bot{},
		items:         map[int32]*gobj.Item{},
		currentMap:    nil,
		maxBotID:      0,
		maxItemID:     0,
		history:       newHistory(),
		botsToCreate:  []*gobj.Bot{},
		botsToDestroy: []int32{},
	}
}

func newHistory() *gobj.History {
	return &gobj.History{
		Maps: [][][]*gobj.LocationInfo{}, // this is dumb. thrift makes [][][]LocationInfo instead of []map
	}
}

func newLocationInfo() *gobj.LocationInfo {
	return &gobj.LocationInfo{
		Bot:     nil,
		Terrain: 0,
	}
}

func copyMap(m gobj.Map) gobj.Map {
	newMap := [][]*gobj.LocationInfo{}

	for i := 0; i < len(m); i++ {
		newMap = append(newMap, []*gobj.LocationInfo{})
		for j := 0; j < len(m[0]); j++ {
			newMap[i] = append(newMap[i], copyLocationInfo(m[i][j]))
		}
	}

	return newMap
}

func copyLocationInfo(loc *gobj.LocationInfo) *gobj.LocationInfo {
	newLoc := newLocationInfo()
	*newLoc = *loc
	return newLoc
}
