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
	"github.com/google/uuid"
	"os"
	"time"
)

type Post struct {
	UserId           string    `dynamodbav:"UserId"`
	PostId           string    `dynamodbav:"MessageId"`
	Message          string    `dynamodbav:"Message"`
	CreatedTimestamp time.Time `dynamodbav:"CreatedTimestamp"`
	UpdatedTimestamp time.Time `dynamodbav:"UpdatedTimestamp"`
}

type Message struct {
	Message string `json:"message"`
}

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	username := req.RequestContext.Authorizer["claims"].(map[string]interface{})["username"].(string)

	var message Message
	err := json.Unmarshal([]byte(req.Body), &message)
	checkErr(err, "Couldn't unmarshall JSON")

	post := Post{
		UserId:           username,
		PostId:           uuid.New().String(),
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
