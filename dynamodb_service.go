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
	"sort"
)

type dynamoService struct {
	Db            dynamodb.DynamoDB
	TableName     string
	LevelFunction LevelFunction
}

const UUID_KEY string = "uuid"
const XP_KEY string = "xp"
const LEVEL_KEY string = "lvl"
const CONSUMABLE_KEY string = "cons"

func (svc dynamoService) addExperience(uuid string, experience int) (success bool, err error) {
	input := getAddExperienceParams(svc.TableName, uuid, experience)
	out, err := svc.Db.UpdateItem(&input)
	if err != nil {
		return false, err
	}
	_ = out;
	svc.startLevelupSequence(uuid, *out)
	return true, nil
}
func (svc dynamoService) getAccount(uuid string) (account *LevelAccount, err error) {
	input := getGetAccountParams(svc.TableName, uuid);
	out, err := svc.Db.GetItem(&input)
	if err != nil {
		return nil, err
	}
	acc, err := getAccountFromAttributes(uuid, out.Item)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (svc dynamoService) consumeLevel(uuid string, level int) (success bool, err error) {
	input := getConsumeLevelParams(svc.TableName, uuid, level)
	out, err := svc.Db.UpdateItem(&input)
	if err != nil {
		return false, err
	}
	var updated bool = out.Attributes[CONSUMABLE_KEY] != nil
	return updated, nil
}

func (svc dynamoService) startLevelupSequence(uuid string, output dynamodb.UpdateItemOutput) {
	levelAccount, err := getAccountFromAttributes(uuid, output.Attributes)
	if err != nil {
		return;
	}
	go svc.levelUp(levelAccount)
}
func (svc dynamoService) levelUp(account LevelAccount) {
	if svc.LevelFunction.canLevel(account.Level, account.Experience) {
		toDeduct := svc.LevelFunction.GetRequiredXp(account.Level)
		in := getLevelUpParams(svc.TableName, account.Uuid, account.Level, toDeduct)
		_, err := svc.Db.UpdateItem(&in)
		if err != nil {
			//Udate failed, probably due to updateconditions

			return
		}
		account.Experience -= svc.LevelFunction.GetRequiredXp(account.Level)
		account.Level += 1
		svc.levelUp(account)
	}
}

func getLevelUpParams(tableName string, uuid string, currentLevel int, deductExperience int) dynamodb.UpdateItemInput {
	updateKey := getUuidUpdateKey(uuid);
	return dynamodb.UpdateItemInput{
		Key: updateKey,
		TableName: aws.String(tableName),
		ExpressionAttributeValues: getLevelUpExpressionAttributeValues(currentLevel, deductExperience),
		ConditionExpression: getLevelUpConditionExpression(currentLevel, deductExperience),
		UpdateExpression: getLevelUpUpdateExpression(currentLevel, deductExperience),
	}
}
func getLevelUpConditionExpression(currentLevel int, deductExperience int) *string {
	var lvlCond string
	if currentLevel == 0 {
		lvlCond = "attribute_not_exists(" + LEVEL_KEY + ")"
	} else {
		lvlCond = LEVEL_KEY + "=:lvlReq"
	}
	return aws.String(lvlCond + " AND " + XP_KEY + ">=:expReq")
}
func getLevelUpUpdateExpression(currentLevel int, deductExperience int) *string {

	return aws.String("SET " + LEVEL_KEY + "=:newLvl ADD " + XP_KEY + " :expAdd," + CONSUMABLE_KEY + " :consSet")
}
func getLevelUpExpressionAttributeValues(lvl int, xp int) map[string]*dynamodb.AttributeValue {
	lvlString := strconv.Itoa(lvl)
	newLvlString := strconv.Itoa(lvl + 1)
	xpString := strconv.Itoa(xp)
	xpAddString := strconv.Itoa(-xp)
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{

		":newLvl" : {
			N: aws.String(newLvlString),
		},
		":consSet" : {
			NS: []*string{aws.String(newLvlString), },
		},
		":expReq" : {
			N: aws.String(xpString),
		},
		":expAdd" : {
			N: aws.String(xpAddString),
		},
	}
	if (lvl > 0) {
		expressionAttributeValues[":lvlReq"] = &dynamodb.AttributeValue{
			N: aws.String(lvlString),
		}
	}
	return expressionAttributeValues
}

//Tested!
func getConsumeLevelParams(tableName string, uuid string, level int) dynamodb.UpdateItemInput {
	updateKey := getUuidUpdateKey(uuid);
	return dynamodb.UpdateItemInput{
		Key: updateKey,
		ConditionExpression: aws.String("contains(" + CONSUMABLE_KEY + ",:lvl)"),
		ExpressionAttributeValues: getConsumableLevelExpressionAttributeValues(level),
		UpdateExpression: getConsumableLevelUpdateExpression(),
		TableName: aws.String(tableName),
	}
}
func getConsumableLevelExpressionAttributeValues(lvlToRemove int) (map[string]*dynamodb.AttributeValue) {
	lvlToRemoveString := strconv.Itoa(lvlToRemove)
	return map[string]*dynamodb.AttributeValue{
		":consSet" : {
			NS: []*string{aws.String(lvlToRemoveString), },
		},
		":lvl" : {
			N: aws.String(lvlToRemoveString),
		},
	}
}
func getConsumableLevelUpdateExpression() *string {
	return aws.String("DELETE " + CONSUMABLE_KEY + " :consSet")
}

//Tested!
func getAddExperienceParams(tableName string, uuid string, experience int) dynamodb.UpdateItemInput {
	return dynamodb.UpdateItemInput{
		Key: getUuidUpdateKey(uuid),
		TableName: aws.String(tableName),
		AttributeUpdates: getAddExperienceAttributeUpdates(experience),
		ReturnValues: aws.String(dynamodb.ReturnValueAllNew),
	}
}

//Tested!
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

//Tested!
func getGetAccountParams(tableName string, uuid string) dynamodb.GetItemInput {
	return dynamodb.GetItemInput{
		Key: getUuidUpdateKey(uuid),
		TableName: aws.String(tableName),
	}
}

//Tested!
func getAccountFromAttributes(uuid string, values map[string]*dynamodb.AttributeValue) (LevelAccount, error) {
	acc := LevelAccount{Experience: 0, Level: 0, UnconsumedLevels: make([]int, 0), Uuid: uuid, }

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

//Tested!
func getLevelFromValues(values map[string]*dynamodb.AttributeValue) *int {
	if values[LEVEL_KEY] == nil {
		return nil
	}
	lvlNum, err := strconv.Atoi(*values[LEVEL_KEY].N)
	if err != nil {
		return nil
	}
	return &lvlNum
}

//Tested!
func getXpFromValues(values map[string]*dynamodb.AttributeValue) *int {
	if values[XP_KEY] == nil {
		return nil
	}
	xpNum, err := strconv.Atoi(*values[XP_KEY].N)
	if err != nil {
		return nil
	}
	return &xpNum
}

//Tested!
func getConsumableFromValues(values map[string]*dynamodb.AttributeValue) *[]int {
	if values[CONSUMABLE_KEY] == nil {
		return nil
	}
	consumableAttributes := values[CONSUMABLE_KEY].NS
	var length = len(consumableAttributes)
	consumable := make([]int, length)
	i := 0;
	for _, k := range consumableAttributes {
		consumableNum, err := strconv.Atoi(*k)
		if err != nil {
			return nil
		}
		consumable[i] = consumableNum
		i++
	}
	sort.Ints(consumable)
	return &consumable
}

//Tested!
func getUuidUpdateKey(uuid string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		UUID_KEY: {
			S: aws.String(uuid),
		},
	}
}
