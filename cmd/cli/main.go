package main

import (
	"context"
	"fmt"

	"github.com/timjdinkins/go-allowance/internal/database/dynamo"
)

func main() {
  repo, err := dynamo.New()
  if err != nil {
    panic(err)
  }

  accs, err := repo.GetAll(context.TODO())
  if err != nil {
    panic(err)
  }
  fmt.Printf("Accounts: %+v\n", accs)
}
