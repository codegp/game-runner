package gameinfo

import (
	"fmt"
	"github.com/codegp/cloud-persister"
	"github.com/codegp/cloud-persister/models"
	"github.com/codegp/game-object-types/types"
)

type GameInfo struct {
	botTypeMap     map[int64]*types.BotType
	attackTypeMap  map[int64]*types.AttackType
	itemTypeMap    map[int64]*types.ItemType
	terrainTypeMap map[int64]*types.TerrainType
	moveTypeMap    map[int64]*types.MoveType
	gameType       *models.GameType
	mapID          int64
}

func (gi *GameInfo) BotType(id int64) *types.BotType {
	return gi.botTypeMap[id]
}

func (gi *GameInfo) AttackType(id int64) *types.AttackType {
	return gi.attackTypeMap[id]
}

func (gi *GameInfo) ItemType(id int64) *types.ItemType {
	return gi.itemTypeMap[id]
}

func (gi *GameInfo) TerrainType(id int64) *types.TerrainType {
	return gi.terrainTypeMap[id]
}

func (gi *GameInfo) MoveType(id int64) *types.MoveType {
	return gi.moveTypeMap[id]
}

func (gi *GameInfo) GameType() *models.GameType {
	return gi.gameType
}

func (gi *GameInfo) MapID() int64 {
	return gi.mapID
}

func NewGameInfo(cp *cloudpersister.CloudPersister, game *models.Game) (*GameInfo, error) {
	attackTypeMap := map[int64]*types.AttackType{}
	botTypeMap := map[int64]*types.BotType{}
	itemTypeMap := map[int64]*types.ItemType{}
	terrainTypeMap := map[int64]*types.TerrainType{}
	moveTypeMap := map[int64]*types.MoveType{}

	gameType, err := cp.GetGameType(game.GameTypeID)
	if err != nil {
		return nil, err
	}

	for _, botTypeID := range gameType.BotTypes {
		botType, err := cp.GetBotType(botTypeID)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve bot type: %v", botTypeID)
		}
		botTypeMap[botType.ID] = botType
	}

	for _, itemTypeID := range gameType.ItemTypes {
		itemType, err := cp.GetItemType(itemTypeID)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve item type: %v", itemTypeID)
		}
		itemTypeMap[itemType.ID] = itemType
	}

	for _, terrainTypeID := range gameType.TerrainTypes {
		terrainType, err := cp.GetTerrainType(terrainTypeID)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve terrain type: %v", terrainTypeID)
		}
		terrainTypeMap[terrainType.ID] = terrainType
	}

	for _, botType := range botTypeMap {
		for _, attackTypeID := range botType.AttackTypeIDs {
			if _, exists := attackTypeMap[attackTypeID]; exists {
				continue
			}

			attackType, err := cp.GetAttackType(attackTypeID)
			if err != nil {
				return nil, fmt.Errorf("Failed to retrieve attack type: %v", attackTypeID)
			}
			attackTypeMap[attackType.ID] = attackType
		}

		for _, moveTypeID := range botType.MoveTypeIDs {
			if _, exists := moveTypeMap[moveTypeID]; exists {
				continue
			}

			moveType, err := cp.GetMoveType(moveTypeID)
			if err != nil {
				return nil, fmt.Errorf("Failed to retrieve move type: %v", moveTypeID)
			}
			moveTypeMap[moveType.ID] = moveType
		}
	}

	return &GameInfo{
		botTypeMap:     botTypeMap,
		attackTypeMap:  attackTypeMap,
		itemTypeMap:    itemTypeMap,
		terrainTypeMap: terrainTypeMap,
		moveTypeMap:    moveTypeMap,
		gameType:       gameType,
		mapID:          game.MapID,
	}, nil
}
