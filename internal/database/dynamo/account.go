package dynamo

import (
	"context"
	"errors"
  "strconv"

	"github.com/google/uuid"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timjdinkins/go-allowance/internal/account"
)

func (db DynamoRepo) GetAllAccounts(ctx context.Context) ([]account.Account, error) {
  accounts := make([]account.Account, 0)
  var token map[string]types.AttributeValue

  for {
    result, err := db.client.Scan(context.TODO(), &dynamodb.ScanInput{
      TableName: awsTableName,
      ExclusiveStartKey: token,
      FilterExpression: aws.String("pk = :pk AND begins_with(sk, :sk)"),
      ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk": &types.AttributeValueMemberS{Value: "family#"+db.familyID},
        ":sk": &types.AttributeValueMemberS{Value: "account#"},
      },
    })
    if err != nil {
      return []account.Account{}, err
    }
    var fetchedAccounts []account.Account
    err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedAccounts)
    if err != nil {
      return []account.Account{}, err
    }

    accounts = append(accounts, fetchedAccounts...)
    token = result.LastEvaluatedKey
    if token == nil {
      break
    }
  }
  return accounts, nil
}

func (db DynamoRepo) GetAccountById(ctx context.Context, id string) (account.Account, error) {
  strPk := "family#"+db.familyID
  strSk := "account#"+id
	acc := account.Account{}

  err := getItem(ctx, db, strPk, strSk, &acc)

  transactions, err := getTransactionsForAccount(ctx, db, strSk)
  if err != nil {
    return account.Account{}, err
  }

  acc.Transactions = append(acc.Transactions, transactions...)
	return acc, nil
}

func getTransactionsForAccount(ctx context.Context, db DynamoRepo, pk string) ([]account.Transaction, error) {
  ts := make([]account.Transaction, 0)
  var token map[string]types.AttributeValue

  for {
    result, err := db.client.Scan(ctx, &dynamodb.ScanInput{
      TableName: awsTableName,
      ExclusiveStartKey: token,
      FilterExpression: aws.String("pk = :pk AND begins_with(sk, :sk)"),
      ExpressionAttributeValues: map[string]types.AttributeValue{
        ":pk": &types.AttributeValueMemberS{Value: pk},
        ":sk": &types.AttributeValueMemberS{Value: "transaction#"},
      },
    })
    if err != nil {
      return []account.Transaction{}, err
    }
    var fetchedTransactions []account.Transaction
    err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedTransactions)
    if err != nil {
      return []account.Transaction{}, err
    }

    ts = append(ts, fetchedTransactions...)
    token = result.LastEvaluatedKey
    if token == nil {
      break
    }
  }
  return ts, nil
}

func (db DynamoRepo) CreateAccount(ctx context.Context, a account.Account) (account.Account, error) {
  a.ID = uuid.New().String()
  strBalance := strconv.Itoa(a.Balance)
  _, err := db.client.PutItem(ctx, &dynamodb.PutItemInput{
    TableName: awsTableName,
    Item: map[string]types.AttributeValue{
      "pk": &types.AttributeValueMemberS{Value: "family#"+db.familyID},
      "sk": &types.AttributeValueMemberS{Value: "account#"+a.ID},
      "id": &types.AttributeValueMemberS{Value: a.ID},
      "firstName": &types.AttributeValueMemberS{Value: a.FirstName},
      "lastName": &types.AttributeValueMemberS{Value: a.LastName},
      "balance": &types.AttributeValueMemberN{Value: strBalance},
    },
  })
  if err != nil {
    return account.Account{}, err
  }
  return a, nil
}

func (db DynamoRepo) CreateTransaction(ctx context.Context, a account.Account, t account.Transaction) (account.Transaction, error) {

  t.ID = uuid.New().String()

  if t.Type == "debit" {
    a.Balance = a.Balance - t.Amount
  } else {
    a.Balance = a.Balance + t.Amount
  }

  if a.Balance < 0 {
    return t, errors.New("Overdraft")
  }

  strAmount := strconv.Itoa(t.Amount)
  strBalance := strconv.Itoa(a.Balance)

  _, err := db.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
    TransactItems: []types.TransactWriteItem{
      {
        Put: &types.Put{
          TableName: awsTableName,
          Item: map[string]types.AttributeValue{
            "pk":    &types.AttributeValueMemberS{Value: "account#"+a.ID},
            "sk":    &types.AttributeValueMemberS{Value: "transaction#"+t.ID},
            "id":    &types.AttributeValueMemberS{Value: t.ID},
            "type":  &types.AttributeValueMemberS{Value: t.Type},
            "amount": &types.AttributeValueMemberN{Value: strAmount},
            "description": &types.AttributeValueMemberS{Value: t.Description},
          },
        },
      },
      {
        Update: &types.Update{
          TableName: awsTableName,
          UpdateExpression: aws.String("SET balance = :balance"),
          ExpressionAttributeValues: map[string]types.AttributeValue{
            ":balance": &types.AttributeValueMemberN{Value: strBalance},
          },
          Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "family#"+db.familyID},
            "sk": &types.AttributeValueMemberS{Value: "account#"+a.ID},
          },
        },
      },
    },
  })

  if err != nil {
    return account.Transaction{}, err
  }

  return t, nil
}
