package gamestate

import (
	"encoding/json"
	"log"

	"github.com/codegp/cloud-persister"
	gobj "github.com/codegp/game-runner/gameobjects"
)

type dummyMap [][]*dumbLocInfo

type dumbLocInfo struct {
	Bot     *dumbBot
	Terrain int64
	Item    int64
}

type dumbBot struct {
	TypeID int64
	Team   int32
}

func (u *GameStateUtils) parseMap(cp *cloudpersister.CloudPersister) {
	gameInfo := u.GameInfo()
	var dumbMap dummyMap

	mapEntity, err := cp.GetMap(gameInfo.MapID())
	if err != nil {
		log.Fatalf("could not parse mapID, %v", err)
	}

	raw, err := cp.ReadMap(mapEntity.ID)
	if err != nil {
		log.Fatalf(err.Error())
	}
	json.Unmarshal(raw, &dumbMap)

	initialMap := [][]*gobj.LocationInfo{}

	for i := range dumbMap {
		initialMap = append(initialMap, []*gobj.LocationInfo{})
		for j, dumbLocInfo := range dumbMap[i] {
			locInfo := &gobj.LocationInfo{
				Loc: &gobj.Location{
					X: int32(i),
					Y: int32(j)},
				Terrain: dumbLocInfo.Terrain,
				Item:    nil,
				Bot:     nil,
			}

			initialMap[i] = append(initialMap[i], locInfo)
		}
	}
	u.gs.currentMap = initialMap
	for i := range dumbMap {
		for j, dumbLoc := range dumbMap[i] {
			loc := &gobj.Location{
				X: int32(i),
				Y: int32(j),
			}

			item := dumbLoc.Item
			if item != 0 {
				if itemType := gameInfo.ItemType(item); itemType == nil {
					log.Printf("Invalid itemType %d\n", item)
				} else if _, err := u.InitItem(loc, itemType); err != nil {
					log.Println(err)
				}
			}

			bot := dumbLoc.Bot
			if bot != nil {
				if botType := gameInfo.BotType(bot.TypeID); botType == nil {
					log.Printf("Invalid botType %d\n", bot.TypeID)
				} else if _, err := u.InitBot(bot.Team-1, loc, botType); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
