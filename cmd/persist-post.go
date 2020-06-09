package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"time"
)

type Post struct {
	UserId string `json:"userId" dynamodbav:"UserId"`
	Message string `json:"message" dynamodbav:"Message"`
	CreatedTimestamp time.Time `json:"createdTimestamp" dynamodbav:"CreatedTimestamp"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" dynamodbav:"UpdatedTimestamp"`
}

func handle(_ context.Context, post Post) (string, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		fmt.Println("Couldn't marshall post:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("posts"),
	}

	avOut, err := svc.PutItem(input)
	if err != nil {
		fmt.Println("Could not PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return avOut.String(), nil
}

func main() {
	lambda.Start(handle)
}