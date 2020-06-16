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
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := util.InitDynamoConnection()
	username := util.GetUsername(req)

	var post model.Post
	err := json.Unmarshal([]byte(req.Body), &post)
	if err != nil {
		fmt.Println("Couldn't unmarshal request:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't read request")
	}

	_, err = svc.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"MessageId": {
				S: aws.String(post.MessageId),
			},
			"UserId": {
				S: aws.String(username),
			},
		},
		TableName: aws.String("posts"),
	})
	if err != nil {
		fmt.Println("Couldn't delete item from posts:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("couldn't delete post")
	}

	return events.APIGatewayProxyResponse{StatusCode: 204}, nil
}

func main() {
	lambda.Start(handler)
}
