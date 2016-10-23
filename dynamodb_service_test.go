package main

import (
	"testing"
	"github.com/satori/go.uuid"
	//"github.com/aws/aws-sdk-go/service/dynamodb"
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