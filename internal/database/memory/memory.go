package memory

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/timjdinkins/go-allowance/internal/account"
)

type MemoryRepo struct {
  Family account.Family
  Accounts []account.Account
}

func NewWithTestData() *MemoryRepo {
  repo := MemoryRepo{}
  repo.loadTestData()
  return &repo
}

func (db *MemoryRepo) CreateFamily(ctx context.Context, email string, phone string) (account.Family, error) {
  return account.Family{
    ID: getUUID(),
    Email: email,
    Phone: phone,
  }, nil
}

func (db *MemoryRepo) GetFamily(ctx context.Context, id string) (account.Family, error) {
  return db.Family, nil
}

func (db *MemoryRepo) GetAllAccounts(ctx context.Context) ([]account.Account, error) {
  return db.Accounts, nil
}

func (db *MemoryRepo) GetAccountById(ctx context.Context, id string) (account.Account, error) {
  var acc account.Account
  for _, a := range db.Accounts {
    if a.ID == id {
      return a, nil
    }
  }
  return acc, errors.New("No account found with that ID")
}

func (db *MemoryRepo) CreateAccount(ctx context.Context, a account.Account) (account.Account, error) {
  a.ID = getUUID()
  db.Accounts = append(db.Accounts, a)
  return a, nil
}

func (db *MemoryRepo) CreateTransaction(ctx context.Context, a account.Account, t account.Transaction) (account.Transaction, error) {
  t.ID = getUUID()
  var acc *account.Account
  for i, aa := range db.Accounts {
    if a.ID == aa.ID {
      acc = &db.Accounts[i]
    }
  }
  acc.Transactions = append(a.Transactions, t)
  if t.Type == "deposit" {
    acc.Balance = acc.Balance + t.Amount
  } else {
    acc.Balance = acc.Balance - t.Amount
  }
  return t, nil
}

func (db *MemoryRepo) loadTestData() {
  t := account.Transaction{
    ID: getUUID(),
    Type: "deposit",
    Amount: 100,
    Description: "Initial Deposit",
  }
  db.Family = account.Family{
    ID: getUUID(),
    Email: "tim@dinkins.me",
    Phone: "702-806-3442",
  }
  acc := account.Account{
    ID: getUUID(),
    FirstName: "Fin",
    LastName: "Dinkins",
    Balance: 100,
    Transactions: []account.Transaction{t},
  }
  db.Accounts = append(db.Accounts, acc)
}

func getUUID() string {
  return uuid.New().String()
}
