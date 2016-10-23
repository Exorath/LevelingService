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

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
)

type dynamoService struct {
	Db        dynamodb.DynamoDB
	TableName string
}

const UUID_KEY string = "uuid"
const XP_KEY string = "xp"
const LEVEL_KEY string = "lvl"
const CONSUMABLE_KEY string = "cons"

func (svc dynamoService) addExperience(uuid string, experience int) (success bool, err error) {
	input := getAddExperienceParams(svc.TableName, uuid, experience)
	out, err := svc.Db.UpdateItem(&input)
	if (err != nil) {
		return false, err
	}
	_ = out;
	return true, nil
}

func (svc dynamoService) getAccount(uuid string) (account *LevelAccount, err error) {
	input := getGetAccountParams(svc.TableName, uuid);
	out, err := svc.Db.GetItem(&input)
	if (err != nil || out.Item == nil) {
		return nil, err
	}
	acc, err := getAccountFromAttributes(uuid, out.Item)
	if (err != nil) {
		return nil, err
	}
	return &acc, nil
}

func (svc dynamoService) consumeLevel(uuid string, level int) (success bool, err error) {
	//TODO: IMPLEMENTATION
	return false, nil
}

func getAddExperienceParams(tableName string, uuid string, experience int) dynamodb.UpdateItemInput {
	return dynamodb.UpdateItemInput{
		Key: getUuidUpdateKey(uuid),
		TableName: aws.String(tableName),
		AttributeUpdates: getAddExperienceAttributeUpdates(experience),
	}
}

func getUuidUpdateKey(uuid string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		UUID_KEY: {
			S: aws.String(uuid),
		},
	}
}

func getAddExperienceAttributeUpdates(experience int) map[string]*dynamodb.AttributeValueUpdate {
	return map[string]*dynamodb.AttributeValueUpdate{
		XP_KEY: {
			Action: aws.String(dynamodb.AttributeActionAdd),
			Value: &dynamodb.AttributeValue{
				N: aws.String(strconv.Itoa(experience)),
			},
		},
	}
}

func getGetAccountParams(tableName string, uuid string) dynamodb.GetItemInput {
	return dynamodb.GetItemInput{
		Key: getUuidUpdateKey(uuid),
		TableName: aws.String(tableName),
	}
}

func getAccountFromAttributes(uuid string, values map[string]*dynamodb.AttributeValue) (LevelAccount, error) {
	acc := LevelAccount{Experience: 0, Level: 0, UnconsumedLevels: make([]int, 0), Uuid: uuid,}

	if tLvl := getLevelFromValues(values); tLvl != nil {
		acc.Level = *tLvl
	}
	if tXp := getXpFromValues(values); tXp != nil {
		acc.Experience = *tXp
	}
	if tConsumable := getConsumableFromValues(values); tConsumable != nil {
		acc.UnconsumedLevels = *tConsumable;
	}
	return acc, nil
}

func getLevelFromValues(values map[string]*dynamodb.AttributeValue) *int {
	if (values[LEVEL_KEY] == nil) {
		return nil
	}
	lvlNum, err := strconv.Atoi(*values[LEVEL_KEY].N)
	if (err != nil) {
		return nil
	}
	return &lvlNum
}
func getXpFromValues(values map[string]*dynamodb.AttributeValue) *int {
	if (values[XP_KEY] == nil) {
		return nil
	}
	xpNum, err := strconv.Atoi(*values[XP_KEY].N)
	if (err != nil) {
		return nil
	}
	return &xpNum
}
func getConsumableFromValues(values map[string]*dynamodb.AttributeValue) *[]int {
	if (values[CONSUMABLE_KEY] == nil) {
		return nil
	}
	consumableAttributes := values[CONSUMABLE_KEY].L
	var length = len(consumableAttributes)
	consumable := make([]int, length)
	for i := 0; i < length; i++ {
		consumableNum, err := strconv.Atoi(*consumableAttributes[i].N)
		if (err != nil) {
			return nil
		}
		consumable[i] = consumableNum
	}
	return &consumable
}
