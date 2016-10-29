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
	"golang.org/x/net/context"
	"net/http"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"log"
	"errors"
	"strconv"
)

func main() {
	region:= os.Getenv("AWS_REGION")
	baseExperience, err := strconv.Atoi(os.Getenv("LS_BASE_XP"))
	if err != nil{
		log.Fatal(errors.New("Environment variable LS_BASE_XP was not configured correctly."))
		return;
	}
	incrementExperience, err := strconv.Atoi(os.Getenv("LS_INC_XP"))
	if err != nil{
		log.Fatal(errors.New("Environment variable LS_INC_XP was not configured correctly."))
		return;
	}
	if region == ""{
		log.Fatal(errors.New("Did not find AWS_REGION environment variable."))
		return;
	}
	db := *dynamodb.New(session.New(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewEnvCredentials(),

	}))
	ctx := context.Background()
	svc := dynamoService{Db: db, TableName: "test", LevelFunction: LevelFunction{BaseExperience: baseExperience, IncrementExperience: incrementExperience}}

	http.ListenAndServe(":8080", MakeHTTPHandler(ctx, svc))
}