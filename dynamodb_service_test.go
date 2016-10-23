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
)

func TestGetGetAccountParams(t *testing.T) {
	var tableName string = "TestTable"
	var uuid string = uuid.NewV4().String()
	itemInput := getGetAccountParams(tableName, uuid)
	if tUuid := *itemInput.Key[UUID_KEY].S; tUuid != uuid {
		t.Error("Expected", uuid, "got", tUuid)
	}
	if tTableName := *itemInput.TableName; tableName != tTableName {
		t.Error("Expected", tableName, "got", tTableName)
	}

	if tLength := len(itemInput.Key); tLength != 1 {
		t.Error("Expected", 1, "got", tLength)
	}
}

func TestGetUuidUpdateKey(t *testing.T) {
	var uuid string = uuid.NewV4().String()
	updateKey := getUuidUpdateKey(uuid)
	if tUuid := *updateKey[UUID_KEY].S; uuid != tUuid {
		t.Error("Expected", uuid, "got", tUuid)
	}
	if tLength := len(updateKey); tLength != 1 {
		t.Error("Expected", 1, "got", tLength)
	}
}

func TestGetAddExperienceAttributeUpdates(t *testing.T) {
	var experience int = 125
	updateKey := getAddExperienceAttributeUpdates(experience)
	tExperience, err := strconv.Atoi(*updateKey[XP_KEY].Value.N)
	if err != nil {
		t.Error(err)
	}
	if experience != tExperience {
		t.Error("Expected", experience, "got", tExperience)
	}
	action := dynamodb.AttributeActionAdd
	if tAction := *updateKey[XP_KEY].Action; action != tAction {
		t.Error("Expected", action, "got", tAction)
	}
	if tLength := len(updateKey); tLength != 1 {
		t.Error("Expected", 1, "got", tLength)
	}
}