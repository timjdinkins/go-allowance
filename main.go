package main

import (
	"github.com/timjdinkins/go-allowance/internal/database/memory"
	"github.com/timjdinkins/go-allowance/internal/web"
)

func main() {
  repo := memory.NewWithTestData()
  web.Start(repo)
}
