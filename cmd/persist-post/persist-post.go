package main

import (
	"aws-demo-lambdas/internal/model"
	"aws-demo-lambdas/internal/util"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"os"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := util.InitDynamoConnection()
	username := util.GetUsername(req)

	var message Message
	err := json.Unmarshal([]byte(req.Body), &message)
	checkErr(err, "Couldn't unmarshall JSON")

	post := model.Post{
		UserId:           username,
		MessageId:        uuid.New().String(),
		Message:          message.Message,
		CreatedTimestamp: time.Now(),
		UpdatedTimestamp: time.Now(),
	}
	av, err := dynamodbattribute.MarshalMap(post)
	checkErr(err, "Couldn't marshall post")

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("posts"),
	}

	_, err = svc.PutItem(input)
	checkErr(err, "Could not PutItem")

	jsonOut, err := json.Marshal(post)
	checkErr(err, "Error marshalling output to json")
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(jsonOut)}, nil
}

func checkErr(err error, msg string) {
	if err != nil {
		fmt.Println(msg + ": ")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	lambda.Start(handle)
}
