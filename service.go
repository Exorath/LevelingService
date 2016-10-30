/*
 *    Copyright 2016 Exorath
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */
package main

//LevelingService provides access and modification to player's experience.
type LevelingService interface {
	addExperience(uuid string, experience int) (bool, error)
	getAccount(uuid string) (*LevelAccount, error)
	consumeLevel(uuid string, level int) (bool, error)
}

type levelingService struct{}


func (levelingService) addExperience(uuid string, experience int) (success bool, err error) {
	//TODO: IMPLEMENTATION
	return false, nil
}

func (levelingService) getAccount(uuid string)  (account *LevelAccount, err error){
	//TODO: IMPLEMENTATION
	return &LevelAccount{uuid, 0, 0, make([]int, 0)}, nil
}

func (levelingService) consumeLevel(uuid string, level int) (success bool, err error){
	//TODO: IMPLEMENTATION
	return false, nil;
}