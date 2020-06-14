package model

import "time"

type Post struct {
	UserId           string    `json:"userId" dynamodbav:"UserId"`
	MessageId        string    `json:"messageId" dynamodbav:"MessageId"`
	Message          string    `json:"message" dynamodbav:"Message"`
	CreatedTimestamp time.Time `json:"createdTimestamp" dynamodbav:"CreatedTimestamp"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" dynamodbav:"UpdatedTimestamp"`
}
