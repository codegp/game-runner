exception InvalidMove {
    1: string message;
}

enum Direction {
  NORTH,
  NORTH_EAST,
  EAST,
  SOUTH_EAST,
  SOUTH,
  SOUTH_WEST,
  WEST,
  NORTH_WEST,
  NONE
}


/*
  A simple struct to represent the x and y coordinates used to index into the Map
*/
struct Location  {
  1: required i32 x;
  2: required i32 y;
}


struct Item {
  1: required i64 itemTypeID;
  2: required i32 ID;
}


/*
` @prop ID - a unique id for the bot
  @prop loc - the bots current location on the map
  @prop health - the bots remaining health
  @prop attackDelay - the bots current attack delay. if > 1 the bot can't attack in the current turn
  @prop moveDelay - the bots current move delay. if > 1 the bot can't move in the current turn
  @prop teamID - a unique id for the team the bot is a part of. -1 is no team
  @prop type - the bot type
  @prop items - list of items the bot is currently holding
*/
struct Bot {
  1: required i32 ID;
  2: required Location loc;
  3: required double health;
  4: required double attackDelay;
  5: required double moveDelay;
  6: required double spawnDelay;
  7: required i64 teamID;
  8: required i64 botTypeID;
  9: required list<Item> items;
}



/*
  A LocationInfo represents the current state of one location in the game world.
  A location must have a terrain type and may contain one Bot and one Item at
  any given time.
  @prop TerrainType - the terrain on the location, eg water, lava, sannd
  @prop Bot - the bot this is currently occupying the location
  @prop Item - an item residing on the location. When a bot moves onto the
    location the Item will be moved to the bots item list
*/
struct LocationInfo {
  1: required i64 terrain;
  2: required Bot bot;
  3: required Item item;
  4: required Location loc;
}

/*
  The state of the game world at any given time is represented by a 2D array
  of LocationInfos. The first array represents the x coor and the second
  represents the y, eg Map[2][10] corresponds to api.Location(2, 10)
*/
typedef list<list<LocationInfo>> Map

/*
  History keeps a list of actions for each round that have passed for the sake of replay
  @prop actionID - 0: move, 1: attack, 2: itemgrab
  @prop target - the location move to, attacked, or grabbed from
*/
struct Action {
  1: required i64 ID;
  2: required Location source;
  3: required Location target;
  4: required i64 actionTypeID;
}

/*
  History a log of the actions and game state at each round
  @prop maps - list of maps, one for each round
  @prop actions - list of list of actions, one list of actions for each round
*/
struct History {
  1: required list<Map> maps;
  2: required list<list<Action>> actions;
}
