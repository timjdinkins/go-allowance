package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timjdinkins/go-allowance/internal/database/dynamo"
	"github.com/timjdinkins/go-allowance/internal/transport"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  familyId := request.PathParameters["id"]
  repo, err := dynamo.New(familyId)
  if err != nil {
    return transport.GwResponse(err.Error(), 500), err
  }

  id := request.PathParameters["id"]
  if id == "" {
    return transport.GwResponse("No id given", 404), nil
  }

  acc, err := repo.GetAccountById(context.TODO(), id)
  if err != nil {
    return transport.GwResponse(err.Error(), 404), nil
  }

  var accJson []byte
  accJson, responseJsonErr := json.Marshal(acc)
  if responseJsonErr != nil {
    return transport.GwResponse(responseJsonErr.Error(), 500), responseJsonErr
  }

  return transport.GwResponse(string(accJson), 200), nil
}

func main() {
  lambda.Start(handler)
}
