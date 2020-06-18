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
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := util.InitDynamoConnection()
	username := util.GetUsername(req)

	result, err := svc.Query(&dynamodb.QueryInput{
		ConsistentRead:         aws.Bool(true),
		TableName:              aws.String("posts"),
		KeyConditionExpression: aws.String("#userid = :userid"),
		ExpressionAttributeNames: map[string]*string{
			"#userid": aws.String("UserId"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userid": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to get posts from dynamodb:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't get posts for user")
	}

	var posts []model.Post
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		fmt.Println("Couldn't unmarshal items from dynamodb:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't return posts for user")
	}

	jsonOut, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("Couldn't marshal posts to JSON:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't format posts to return")
	}
	return events.APIGatewayProxyResponse{StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "POST,PUT,OPTIONS,DELETE",
			"Access-Control-Max-Age":       "86400",
		},
		Body: string(jsonOut)}, nil
}

func main() {
	lambda.Start(handler)
}
