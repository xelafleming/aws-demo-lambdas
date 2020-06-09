package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"time"
)

type Post struct {
	UserId           string    `json:"userId" dynamodbav:"UserId"`
	Message          string    `json:"message" dynamodbav:"Message"`
	CreatedTimestamp time.Time `json:"createdTimestamp" dynamodbav:"CreatedTimestamp"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" dynamodbav:"UpdatedTimestamp"`
}

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	var post Post
	err := json.Unmarshal([]byte(req.Body), &post)
	checkErr(err, "Couldn't unmarshall JSON")

	av, err := dynamodbattribute.MarshalMap(post)
	checkErr(err, "Couldn't marshall post")

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("posts"),
	}

	avOut, err := svc.PutItem(input)
	checkErr(err, "Could not PutItem")

	jsonOut, err := json.Marshal(avOut)
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
