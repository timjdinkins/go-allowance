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
  familyId := request.PathParameters["familyId"]
  repo, err := dynamo.New(familyId)
  if err != nil {
    return transport.GwResponse(err.Error(), 500), err
  }

  accs, err := repo.GetAllAccounts(context.TODO())
  if err != nil {
    return transport.GwResponse(err.Error(), 500), err
  }

  var accJson []byte
  accJson, responseJsonErr := json.Marshal(accs)
  if responseJsonErr != nil {
    return transport.GwResponse(responseJsonErr.Error(), 500), err
  }

  return transport.GwResponse(string(accJson), 200), nil
}

func main() {
  lambda.Start(handler)
}
