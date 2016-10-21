# LevelingService
Exorath's leveling Microservice written in Go.

##Endpoints

###/player [POST]: 
####Adds experience to a user's leveling-service account (And levels the player up if neccessary).
**Arguments**:
- uuid (string): the uuid of the player
- xp (integer): The amount of experience to grant the player

**Response**: JSON formatted updated document, according to GET /player with the player uuid

**Behind the scenes (Design doc)**: When this endpoint is called, an increment is send to the DB. While the player has more experience then required for a level up, he will keep leveling up (as soon as there's not enough experience, this stops). The level-up query is molecular: It requires demands for enough experience to deduct, and the level must be not have increased, if these requirements are not met, the level-up loop terminates.

###/player [GET]: 
####gets current experience, current level and unconsumed levels
**Body**: {"uuid": "71a2d31e-646a-4bf6-9883-448eaf81effb"}
- uuid (string): the uuid of the player

**Response**: {"xp": 1624, "lvl": 20, "consumable":[19,20]}
- xp: the current experience of the player
- lvl: the current level of the player
- consumable: the levels that can be processed for a reward (Maybe a more advanced consumable will be a cool feature in the future, allowing for multiple consumers).

###/player/consume [POST]: 
####Consumes an unconsumed level. Should be called when a reward was granted for the level up (For the reward to be a real transaction, post the consumed level somewhere in the reward document, and check if it already exists to not consume twice before this method can be executed)

**Body**: {"uuid": "71a2d31e-646a-4bf6-9883-448eaf81effb", "lvl": 19}
- uuid (string): the uuid of the player
- lvl (number): the level to consume

**Response**: {"success": true}
- success: whether or not the consumption succeeded
