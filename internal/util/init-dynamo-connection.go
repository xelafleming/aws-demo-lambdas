package util

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func InitDynamoConnection() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	return dynamodb.New(sess)
}
