package account

import "context"

type Repo interface {
  GetAllAccounts(context.Context) ([]Account, error)
  GetAccountById(context.Context, string) (Account, error)
  CreateAccount(context.Context, Account) (Account, error)
  CreateTransaction(context.Context, Account, Transaction) (Transaction, error)
  CreateFamily(context.Context, string, string) (Family, error)
  GetFamily(context.Context, string) (Family, error)
}

type Family struct {
  ID string `json:"id" dynamodbav:"id"`
  Email string `json:"email" dynamodbav:"email"`
  Phone string `json:"phone" dynamodbav:"phone"`
  Accounts []Account
}

type Account struct {
  ID string `json:"id" dynamodbav:"id"`
  FirstName string `json:"firstName" dynamodbav:"firstName"`
  LastName string `json:"lastName" dynamodbav:"lastName"`
  Balance int `json:"balance" dynamodbav:"balance"`
  Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
  ID string `json:"id" dynamodbav:"id"`
  Type string `json:"type" dynamodbav:"type"` // debit or deposit
  Amount int `json:"amount" dynamodbav:"amount"`
  Description string `json:"description" dynamodbav:"description"`
}
