package main

//LevelingService provides access and modification to player's experience.
type LevelingService interface {
	addExperience(uuid string, experience int) (bool, error)
	getAccount(uuid string) (LevelAccount, error)
	consumeLevel(uuid string, level int) (bool, error)
}

type levelingService struct{}


func (levelingService) addExperience(uuid string, experience int) (success bool, err error) {
	return false, nil
}

func (levelingService) getAccount(uuid string)  (account LevelAccount, err error){
	return LevelAccount{uuid, 0, 0, make([]int, 0)}, nil
}

func (levelingService) consumeLevel(uuid string, level int) (success bool, err error){
	return false, nil;
}