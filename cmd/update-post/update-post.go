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
	"time"
)

type Key struct {
	UserId           string    `dynamodbav:"UserId"`
	MessageId        string    `dynamodbav:"MessageId"`
	CreatedTimestamp time.Time `dynamodbav:"CreatedTimestamp"`
}

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := util.InitDynamoConnection()
	username := util.GetUsername(req)

	var post model.Post
	err := json.Unmarshal([]byte(req.Body), &post)
	if err != nil {
		fmt.Println("Couldn't unmarshal request body: ")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't process post from request")
	}

	if post.UserId != username {
		return events.APIGatewayProxyResponse{StatusCode: 403}, errors.New("calling user is trying to edit a post that does not belong to them")
	}

	key, err := dynamodbattribute.MarshalMap(Key{
		UserId:           post.UserId,
		MessageId:        post.MessageId,
		CreatedTimestamp: post.CreatedTimestamp,
	})
	if err != nil {
		fmt.Println("Got error marshaling key:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't process update key")
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":msg": {
				S: aws.String(post.Message),
			},
			":updts": {
				S: aws.String(post.UpdatedTimestamp.String()),
			},
		},
		UpdateExpression: aws.String("set Message = :msg, UpdatedTimestamp = :updts"),
		TableName:        aws.String("posts"),
		Key:              key,
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		fmt.Println("Couldn't update post:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't update post")
	}

	jsonOut, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Couldn't marshal post for output:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "{\"error\":\"Couldn't return post updated. Post was updated.\"}",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonOut),
	}, nil
}

func main() {
	lambda.Start(handle)
}
