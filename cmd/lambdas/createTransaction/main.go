package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timjdinkins/go-allowance/internal/account"
	"github.com/timjdinkins/go-allowance/internal/database/dynamo"
	"github.com/timjdinkins/go-allowance/internal/transport"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  familyId := request.PathParameters["familyId"]
  if familyId == "" {
    return transport.GwResponse("Invalid request: Unknown account.", 400), nil
  }
  accId := request.PathParameters["id"]
  if accId == "" {
    return transport.GwResponse("Account ID is required", 400), nil
  }
  repo, err := dynamo.New(familyId)
  ctx := context.Background()

  acc, err := repo.GetAccountById(ctx, accId)
  if err != nil {
    return transport.GwResponse(err.Error(), 404), nil
  }

  var trans account.Transaction
  json.Unmarshal([]byte(request.Body), &trans)
  if trans.Amount == 0 {
    return transport.GwResponse("Invalid transaction, no amount given", 400), nil
  }

  trans, err = repo.CreateTransaction(ctx, acc, trans)
  if err != nil {
    return transport.GwResponse(err.Error(), 500), nil
  }

  transJson, err := json.Marshal(trans)
  if err != nil {
    return transport.GwResponse(err.Error(), 500), nil
  }

  return transport.GwResponse(string(transJson), 200), nil
}

func main() {
  lambda.Start(handler)
}
