package transport

import (
	"github.com/aws/aws-lambda-go/events"
)

func GwResponse(body string, code int) events.APIGatewayProxyResponse {
  return events.APIGatewayProxyResponse{
    Body:       body,
    StatusCode: code,
  }
}
