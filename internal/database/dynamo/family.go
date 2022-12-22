package dynamo

import (
	"context"

	"github.com/google/uuid"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timjdinkins/go-allowance/internal/account"
)

func (db DynamoRepo) CreateFamily(ctx context.Context, email string, phone string) (account.Family, error) {
  fam := account.Family{
    ID: uuid.New().String(),
    Email: email,
    Phone: phone,
  }
  _, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
    TableName: awsTableName,
    Item: map[string]types.AttributeValue{
      "pk": &types.AttributeValueMemberS{Value: "family#"+fam.ID},
      "sk": &types.AttributeValueMemberS{Value: "family#"+fam.ID},
      "id": &types.AttributeValueMemberS{Value: fam.ID},
      "email": &types.AttributeValueMemberS{Value: fam.Email},
      "phone": &types.AttributeValueMemberS{Value: fam.Phone},
    },
  })
  if err != nil {
    return account.Family{}, err
  }
  return fam, nil
}

func (db DynamoRepo) GetFamily(ctx context.Context) (account.Family, error) {
  strKey := "family#"+db.familyID
  var fam account.Family
  if err := getItem[account.Family](ctx, db, strKey, strKey, &fam); err != nil {
    return account.Family{}, err
  }

  return fam, nil
}
