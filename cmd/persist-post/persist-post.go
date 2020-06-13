package main

import (
	"aws-demo-lambdas/internal/model"
	"aws-demo-lambdas/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
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
	if err != nil {
		fmt.Println("Couldn't unmarshall message from request:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, errors.New("couldn't process message from request")
	}

	post := model.Post{
		UserId:           username,
		MessageId:        uuid.New().String(),
		Message:          message.Message,
		CreatedTimestamp: time.Now(),
		UpdatedTimestamp: time.Now(),
	}
	av, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		fmt.Println("Couldn't create dynamodbav from post object")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't create post")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("posts"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Couldn't insert item into database:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't store item")
	}

	jsonOut, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Couldn't marshall created item to JSON")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "{\"error\":\"Couldn't return created post. Post was created.\"}",
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(jsonOut)}, nil
}

func main() {
	lambda.Start(handle)
}
