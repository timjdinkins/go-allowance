package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timjdinkins/go-allowance/internal/account"
)


type DynamoRepo struct {
  client *dynamodb.Client
  familyID string
}

type Getable interface {
  account.Family | account.Account
}

const TableName = "Allowance"
var awsTableName *string = aws.String(TableName)
var client *dynamodb.Client

func init() {
  cfg, err := config.LoadDefaultConfig(context.TODO())
  if err != nil {
    return
  }
  client = dynamodb.NewFromConfig(cfg)
}

func New(familID string) (*DynamoRepo, error) {
  return &DynamoRepo{
    client: client,
    familyID: familID,
  }, nil
}

func getItem[T Getable](ctx context.Context, db DynamoRepo, strPk string, strSk string, t *T) error {
  pk, err := attributevalue.Marshal(strPk)
  sk, err1 := attributevalue.Marshal(strSk)
  if err != nil || err1 != nil {
    return nil
  }

  input := &dynamodb.GetItemInput{
    TableName: awsTableName,
    Key: map[string]types.AttributeValue{
      "pk": pk,
      "sk": sk,
    },
  }

	result, err := db.client.GetItem(ctx, input)
	if err != nil || result.Item == nil {
		return err
	}

	err = attributevalue.UnmarshalMap(result.Item, &t)
	if err != nil {
		return err
	}
  
  return nil
}
