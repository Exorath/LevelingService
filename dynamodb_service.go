package main

import (
 	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
)

type dynamoService struct{
	Db dynamodb.DynamoDB
	TableName string
}
const UUID_KEY string = "uuid"
const XP_KEY string = "xp"
const LEVEL_KEY string = "lvl"
const CONSUMABLE_KEY string = "cons"

func (svc dynamoService) addExperience(uuid string, experience int) (success bool, err error) {
	input := getAddExperienceParams(svc.TableName, uuid, experience)
	out, err := svc.Db.UpdateItem(&input)
	if(err != nil){
		return false, err
	}
	_ = out;
	return true, nil
}


func (svc dynamoService) getAccount(uuid string)  (account *LevelAccount, err error){
	input := getGetAccountParams(svc.TableName, uuid);
	out, err := svc.Db.GetItem(&input)
	if(err != nil || out.Item == nil){
		return nil, err
	}
	acc, err := getAccountFromAttributes(uuid, out.Item)
	if(err != nil || acc == nil){
		return nil, err
	}
	return acc, nil
}

func (svc dynamoService) consumeLevel(uuid string, level int) (success bool, err error){
	//TODO: IMPLEMENTATION
	return false, nil
}

func getAddExperienceParams(tableName string, uuid string, experience int) dynamodb.UpdateItemInput {
	var xpString = strconv.Itoa(experience)
	return dynamodb.UpdateItemInput{
		Key: getUuidUpdateKey(uuid),
		TableName: aws.String(tableName),
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			XP_KEY: {
				Action: aws.String(dynamodb.AttributeActionAdd),
				Value: &dynamodb.AttributeValue{
					N: aws.String(xpString),
				},
			},
		},
	}
}
func getGetAccountParams(tableName string, uuid string) dynamodb.GetItemInput{
	return dynamodb.GetItemInput{
		Key: getUuidUpdateKey(uuid),
		TableName: aws.String(tableName),
	}
}

func getAccountFromAttributes(uuid string, values map[string]*dynamodb.AttributeValue) (acc *LevelAccount, err error) {
	level := 0
	xp := 0
	var consumable []int

	if(values[LEVEL_KEY] != nil){
		lvlNum, err := strconv.Atoi(*values[LEVEL_KEY].N)
		if(err != nil){
			return nil, err
		}
		level = lvlNum
	}

	if(values[XP_KEY] != nil){
		xpNum, err := strconv.Atoi(*values[XP_KEY].N)
		if(err != nil){
			return nil, err
		}
		xp = xpNum
	}
	if(values[CONSUMABLE_KEY] != nil){
		consumableAttributes := values[CONSUMABLE_KEY].L
		var length = len(consumableAttributes)
		consumable := make([]int, length)
		for i := 0; i < length; i++ {
			consumableNum, err := strconv.Atoi(*consumableAttributes[i].N)
			if(err != nil){
				return nil, err
			}
			consumable[i] = consumableNum
		}
	}
	return &LevelAccount{
		Experience: xp,
		Level: level,
		UnconsumedLevels: consumable,
		Uuid: uuid,
	}, nil
}

func getUuidUpdateKey(uuid string) map[string]*dynamodb.AttributeValue{
	return map[string]*dynamodb.AttributeValue{
		UUID_KEY: {
			S: aws.String(uuid),
		},
	}
}

