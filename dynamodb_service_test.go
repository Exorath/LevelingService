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
	"testing"
	"github.com/satori/go.uuid"
	"strconv"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)
//**  Interface tests  **//
//Requires dynamodb connection...

//**  Unit tests  **//
//Update to new getconsumelevel system
func TestGetConsumeLevelParams(t *testing.T) {
	var tableName string = "TestTable"
	var uuid string = uuid.NewV4().String()
	var toRemove int = 10;
	itemInput := getConsumeLevelParams(tableName, uuid, toRemove)

	//Key
	if tUuid := *itemInput.Key[UUID_KEY].S; tUuid != uuid {
		t.Error("Expected", uuid, "got", tUuid)
	}
	if tLength := len(itemInput.Key); 1 != tLength {
		t.Error("Expected", 1, "got", tLength)
	}
	//TableName
	if tTableName := *itemInput.TableName; tableName != tTableName {
		t.Error("Expected", tableName, "got", tTableName)
	}
	//ConditionExpression
	conditionExpression := "contains(" + CONSUMABLE_KEY + ",:lvl)"
	if tConditionExpresison := *itemInput.ConditionExpression; conditionExpression != tConditionExpresison {
		t.Error("Expected", conditionExpression, "got", tConditionExpresison)
	}
	//ExpressionAttributeValues
	if tLength := len(itemInput.ExpressionAttributeValues); 2 != tLength {
		t.Error("Expected", 2, "got", tLength)
	}
	if tConsSetLen := len(itemInput.ExpressionAttributeValues[":consSet"].NS); 1 != tConsSetLen {
		t.Error("Expected", 1, "got", tConsSetLen)
	}
	if tConsLvl := *itemInput.ExpressionAttributeValues[":consSet"].NS[0]; "10" != tConsLvl {
		t.Error("Expected", "10", "got", tConsLvl)
	}
	if tLvl := *itemInput.ExpressionAttributeValues[":lvl"].N; "10" != tLvl {
		t.Error("Expected", "10", "got", tLvl)
	}
	//UpdateExpression
	updateExpression := "DELETE " + CONSUMABLE_KEY + " :consSet"
	if tUpdateExpression := *itemInput.UpdateExpression; updateExpression != tUpdateExpression {
		t.Error("Expected", updateExpression, "got", tUpdateExpression)
	}
}

func TestGetLevelUpParams(t *testing.T) {
	var tableName string = "TestTable"
	var uuid string = uuid.NewV4().String()
	var currentLevel int = 3;
	var deductXp int = 100;
	itemInput := getLevelUpParams(tableName, uuid, currentLevel, deductXp)
	//Key
	if tUuid := *itemInput.Key[UUID_KEY].S; tUuid != uuid {
		t.Error("Expected", uuid, "got", tUuid)
	}
	if tLength := len(itemInput.Key); 1 != tLength {
		t.Error("Expected", 1, "got", tLength)
	}
	//TableName
	if tTableName := *itemInput.TableName; tableName != tTableName {
		t.Error("Expected", tableName, "got", tTableName)
	}
	//ExpressionAttributeValues
	if tLength := len(itemInput.ExpressionAttributeValues); 5 != tLength {
		t.Error("Expected", 5, "got", tLength)
	}
	if tNewLvl := *itemInput.ExpressionAttributeValues[":newLvl"].N; "4" != tNewLvl {
		t.Error("Expected", "4", "got", tNewLvl)
	}
	if tLvlReq := *itemInput.ExpressionAttributeValues[":lvlReq"].N; "3" != tLvlReq {
		t.Error("Expected", "3", "got", tLvlReq)
	}

	if tConsSetLen := len(itemInput.ExpressionAttributeValues[":consSet"].NS); 1 != tConsSetLen {
		t.Error("Expected", 1, "got", tConsSetLen)
	}
	if tConsLvl := *itemInput.ExpressionAttributeValues[":consSet"].NS[0]; "4" != tConsLvl {
		t.Error("Expected", "4", "got", tConsLvl)
	}
	if tExpReq := *itemInput.ExpressionAttributeValues[":expReq"].N; "100" != tExpReq {
		t.Error("Expected", "100", "got", tExpReq)
	}
	if tExpAdd := *itemInput.ExpressionAttributeValues[":expAdd"].N; "-100" != tExpAdd {
		t.Error("Expected", "-100", "got", tExpAdd)
	}
	//UpdateExpression
	updateExpression := "SET " + LEVEL_KEY + "=:newLvl ADD " + XP_KEY + " :expAdd," + CONSUMABLE_KEY + " :consSet"
	if tUpdateExpression := *itemInput.UpdateExpression; updateExpression != tUpdateExpression {
		t.Error("Expected", updateExpression, "got", tUpdateExpression)
	}
	//ConditionExpression
	conditionExpression := LEVEL_KEY + "=:lvlReq AND " + XP_KEY + ">=:expReq"
	if tConditionExpression := *itemInput.ConditionExpression; conditionExpression != tConditionExpression {
		t.Error("Expected", conditionExpression, "got", tConditionExpression)
	}
	itemInput2 := getLevelUpParams(tableName, uuid, 0, deductXp)
	conditionExpression2 := "attribute_not_exists(" + LEVEL_KEY + ") AND " + XP_KEY + ">=:expReq"
	if tConditionExpression := *itemInput2.ConditionExpression; conditionExpression2 != tConditionExpression {
		t.Error("Expected", conditionExpression2, "got", tConditionExpression)
	}
}
func TestGetGetAccountParams(t *testing.T) {
	var tableName string = "TestTable"
	var uuid string = uuid.NewV4().String()
	itemInput := getGetAccountParams(tableName, uuid)
	if tUuid := *itemInput.Key[UUID_KEY].S; tUuid != uuid {
		t.Error("Expected", uuid, "got", tUuid)
	}
	if tLength := len(itemInput.Key); 1 != tLength {
		t.Error("Expected", 1, "got", tLength)
	}
	if tTableName := *itemInput.TableName; tableName != tTableName {
		t.Error("Expected", tableName, "got", tTableName)
	}
}
func TestGetAddExperienceParams(t *testing.T) {
	var tableName string = "TestTable"
	var uuid string = uuid.NewV4().String()
	var xp int = 1023
	itemInput := getAddExperienceParams(tableName, uuid, xp)
	if tUuid := *itemInput.Key[UUID_KEY].S; tUuid != uuid {
		t.Error("Expected", uuid, "got", tUuid)
	}
	if tLength := len(itemInput.Key); 1 != tLength {
		t.Error("Expected", 1, "got", tLength)
	}
	if tLength := len(itemInput.AttributeUpdates); 1 != tLength {
		t.Error("Expected", 1, "got", tLength)
	}
	var xpString string = strconv.Itoa(xp);
	if tXpString := *itemInput.AttributeUpdates[XP_KEY].Value.N; tXpString != xpString {
		t.Error("Expected", xpString, "got", tXpString)
	}
	if tTableName := *itemInput.TableName; tableName != tTableName {
		t.Error("Expected", tableName, "got", tTableName)
	}

	if tLength := len(itemInput.Key); tLength != 1 {
		t.Error("Expected", 1, "got", tLength)
	}
	if tReturnValues := *itemInput.ReturnValues; tReturnValues != dynamodb.ReturnValueAllNew {
		t.Error("Expected", dynamodb.ReturnValueAllNew, "got", tReturnValues)
	}

}
func TestGetGetAccountFromAttributes(t *testing.T) {
	var uuid string = uuid.NewV4().String()
	var lvl int = 10
	var xp int = 125
	attributes := map[string]*dynamodb.AttributeValue{
		LEVEL_KEY : {N: aws.String(strconv.Itoa(lvl))},
		XP_KEY : {N: aws.String(strconv.Itoa(xp))},
		CONSUMABLE_KEY : {
			NS: []*string{
				aws.String("5"),
				aws.String("7"),
				aws.String("6"),
			},
		},
	}
	account, err := getAccountFromAttributes(uuid, attributes)
	if (err != nil) {
		t.Error(err)
	}
	if account.Level != lvl {
		t.Error("Expected", lvl, "got", account.Level)
	}
	if account.Experience != xp {
		t.Error("Expected", xp, "got", account.Experience)
	}
	if account.Uuid != uuid {
		t.Error("Expected", uuid, "got", account.Uuid)
	}
	if consumableLength := len(account.UnconsumedLevels); consumableLength != 3 {
		t.Error("Expected", 3, "got", consumableLength)
	}
	if tLevel := account.UnconsumedLevels[0]; tLevel != 5 {
		t.Error("Expected", 5, "got", tLevel)
	}
	if tLevel := account.UnconsumedLevels[1]; tLevel != 6 {
		t.Error("Expected", 6, "got", tLevel)
	}
	if tLevel := account.UnconsumedLevels[2]; tLevel != 7 {
		t.Error("Expected", 7, "got", tLevel)
	}
}

