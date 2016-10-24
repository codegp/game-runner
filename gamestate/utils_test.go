package gamestate

import (
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/codegp/test-utils"
	"github.com/codegp/game-object-types/types"
	gobj "github.com/codegp/game-runner/gameobjects"
)

const testID = 1

func testGameStateUtils(t *testing.T) *GameStateUtils {
  gsu := &GameStateUtils{
    newGameState(),
    testGameInfo(),
		&TestRandomizer{},
  }
  parseTestMap(t, gsu)
  gsu.SetCurrentBot(gsu.Bot(1))
  gsu.ClearPendingBotActions()
  return gsu
}

func testGameInfo() *GameInfo {
  return &GameInfo{
		botTypeMap:     map[int64]*types.BotType{testID: testutils.UnitTestBotType()},
		attackTypeMap:  map[int64]*types.AttackType{testID: testutils.UnitTestAttackType()},
		itemTypeMap:    map[int64]*types.ItemType{testID: testutils.UnitTestItemType()},
		terrainTypeMap: map[int64]*types.TerrainType{testID: testutils.UnitTestTerrainType()},
		moveTypeMap:    map[int64]*types.MoveType{testID: testutils.UnitTestMoveType()},
		gameType:       testutils.UnitTestGameType(),
		mapID:          1,
	}
}

func parseTestMap(t *testing.T, gsu *GameStateUtils) {
  initialMap := [][]*gobj.LocationInfo{}
  for i := 0; i < 2; i++ {
    initialMap = append(initialMap, []*gobj.LocationInfo{})

    for j := 0; j < 2; j++ {
			locInfo := &gobj.LocationInfo{
				Loc: &gobj.Location{
					X: int32(i),
					Y: int32(j)},
				Terrain: 1,
				Item:    nil,
				Bot:     nil,
			}
      initialMap[i] = append(initialMap[i], locInfo)
		}
	}

  gsu.gs.currentMap = initialMap
  _, err := gsu.InitBot(1, &gobj.Location{X:0,Y:0}, gsu.GameInfo().BotType(testID))
  assert.Nil(t, err)
  _, err = gsu.InitBot(2, &gobj.Location{X:1,Y:1}, gsu.GameInfo().BotType(testID))
  assert.Nil(t, err)
}
