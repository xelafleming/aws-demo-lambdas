package model

import "time"

type Post struct {
	UserId           string    `dynamodbav:"UserId"`
	MessageId        string    `dynamodbav:"MessageId"`
	Message          string    `dynamodbav:"Message"`
	CreatedTimestamp time.Time `dynamodbav:"CreatedTimestamp"`
	UpdatedTimestamp time.Time `dynamodbav:"UpdatedTimestamp"`
}
