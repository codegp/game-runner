package gamestate

import (
  "testing"

  "github.com/stretchr/testify/assert"
  	gobj "github.com/codegp/game-runner/gameobjects"
)

/*
  Testing constants are based off the follow assumptions, backed by setup of testGameStateUtils
  CurrentBot is at loc 0,0
  Map is 2x2
  Opponent bot is at loc 1,1
  All bots can move and attack using move/attack key 1 and no others
*/
const (
  validKey = 1
  invalidKey = 2
  offMapDirection = gobj.Direction_NORTH
  occupiedMapDirection = gobj.Direction_SOUTH_EAST
  unoccupiedMapDirection = gobj.Direction_SOUTH
)

func occupiedLoc() *gobj.Location {
  return &gobj.Location{
    X: 1,
    Y: 1,
  }
}

func unoccupiedLoc() *gobj.Location {
  return &gobj.Location{
    X: 0,
    Y: 1,
  }
}

func offMapLoc() *gobj.Location {
  return &gobj.Location{
    X: -1,
    Y: 1,
  }
}

/*
 *   SpawnInDirection / ValidateSpawnAttempt
 */
func TestInvalidSpawnType(t *testing.T) {
  u := testGameStateUtils(t)
  bot, err := u.SpawnInDirection(unoccupiedMapDirection, invalidKey)
  assert.NotNil(t, err)
  assert.Nil(t, bot)
}

func TestHasSpawnDelay(t *testing.T) {
  u := testGameStateUtils(t)
  u.CurrentBot().SpawnDelay = 1
  bot, err := u.SpawnInDirection(unoccupiedMapDirection, validKey)
  assert.NotNil(t, err)
  assert.Nil(t, bot)
}

func TestCanSpawnTypeIsFalse(t *testing.T) {
  u := testGameStateUtils(t)
  u.GameInfo().BotType(testID).CanSpawn = false
  bot, err := u.SpawnInDirection(unoccupiedMapDirection, validKey)
  assert.NotNil(t, err)
  assert.Nil(t, bot)
}

func TestCanBeSpawnedTypeIsFalse(t *testing.T) {
  u := testGameStateUtils(t)
  u.GameInfo().BotType(testID).CanBeSpawned = false
  bot, err := u.SpawnInDirection(unoccupiedMapDirection, validKey)
  assert.NotNil(t, err)
  assert.Nil(t, bot)
}

func TestCannotSpawnOffMap(t *testing.T) {
  u := testGameStateUtils(t)
  bot, err := u.SpawnInDirection(offMapDirection, validKey)
  assert.NotNil(t, err)
  assert.Nil(t, bot)
}

func TestCannotSpawnOnOccupiedLocation(t *testing.T) {
  u := testGameStateUtils(t)
  bot, err := u.SpawnInDirection(occupiedMapDirection, validKey)
  assert.NotNil(t, err)
  assert.Nil(t, bot)
}

func TestValidSpawn(t *testing.T) {
  u := testGameStateUtils(t)
  bot, err := u.SpawnInDirection(unoccupiedMapDirection, validKey)
  assert.NotNil(t, bot)
  assert.Nil(t, err)
  assert.Equal(t, len(u.BotsToCreate()), 1)
  assert.Equal(t, u.BotsToCreate()[0], bot)
  assert.NotNil(t, u.Map()[0][1].Bot)
  assert.Equal(t, u.Map()[0][1].Bot, bot)
}

/*
 *   MoveInDirection / ValidateMoveAttempt
 */
func TestInvalidMoveType(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.MoveInDirection(unoccupiedMapDirection, invalidKey)
  assert.NotNil(t, err)
}

func TestHasMoveDelay(t *testing.T) {
  u := testGameStateUtils(t)
  u.CurrentBot().MoveDelay = 1
  err := u.MoveInDirection(unoccupiedMapDirection, validKey)
  assert.NotNil(t, err)
}

func TestCannotUseMoveType(t *testing.T) {
  u := testGameStateUtils(t)
  u.GameInfo().BotType(testID).MoveTypeIDs = []int64{}
  err := u.MoveInDirection(unoccupiedMapDirection, validKey)
  assert.NotNil(t, err)
}

func TestCannotMoveOffMap(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.MoveInDirection(offMapDirection, validKey)
  assert.NotNil(t, err)
}

func TestCannotMoveOnOccupiedLocation(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.MoveInDirection(occupiedMapDirection, validKey)
  assert.NotNil(t, err)
}

func TestValidMove(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.MoveInDirection(unoccupiedMapDirection, validKey)
  assert.Nil(t, err)
  assert.NotNil(t, u.Map()[0][1].Bot)
  assert.Nil(t, u.Map()[0][0].Bot)
  assert.Equal(t, u.Map()[0][1].Bot.ID, u.CurrentBot().ID)
}

/*
 *   AttackLocation / ValidateAttackAttempt
 */


func TestInvalidAttackType(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.AttackLocation(occupiedLoc(), invalidKey)
  assert.NotNil(t, err)
}

func TestHasAttackDelay(t *testing.T) {
  u := testGameStateUtils(t)
  u.CurrentBot().AttackDelay = 1
  err := u.AttackLocation(occupiedLoc(), validKey)
  assert.NotNil(t, err)
}

func TestCannotUseAttackType(t *testing.T) {
  u := testGameStateUtils(t)
  u.GameInfo().BotType(testID).AttackTypeIDs = []int64{}
  err := u.AttackLocation(occupiedLoc(), validKey)
  assert.NotNil(t, err)
}

func TestCannotAttackOffMap(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.AttackLocation(offMapLoc(), validKey)
  assert.NotNil(t, err)
}

func TestCannotAttackOnUnoccupiedLocation(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.AttackLocation(unoccupiedLoc(), validKey)
  assert.NotNil(t, err)
}

func TestValidAttack(t *testing.T) {
  u := testGameStateUtils(t)
  err := u.AttackLocation(occupiedLoc(), validKey)
  assert.Nil(t, err)
  assert.Equal(t, u.Bot(2).Health, 99.0)
}
