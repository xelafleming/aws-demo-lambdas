package util

import "github.com/aws/aws-lambda-go/events"

func GetUsername(req events.APIGatewayProxyRequest) string {
	return req.RequestContext.Authorizer["claims"].(map[string]interface{})["username"].(string)
}
