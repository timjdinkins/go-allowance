package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/timjdinkins/go-allowance/internal/account"
)

func Start(db account.Repo) {
  r := mux.NewRouter()

  r.HandleFunc("/accounts", getAll(db)).Methods("GET")
  r.HandleFunc("/accounts/{id}", getAccountById(db)).Methods("GET")
  r.HandleFunc("/accounts", createAccount(db)).Methods("POST")
  r.HandleFunc("/accounts/{id}", createTransaction(db)).Methods("POST")

  http.ListenAndServe(":8080", r)
}

func getAll(db account.Repo) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    accs, err := db.GetAllAccounts(ctx)
    if err != nil {
      respondWithError(
        w,
        http.StatusInternalServerError,
        fmt.Sprintf("Error fetching accounts: %s", err.Error()),
      )
    } else {
      respondWithJSON(w, http.StatusOK, accs)
    }

  }
}

func getAccountById(db account.Repo) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := mux.Vars(r)["id"]
    acc, err := db.GetAccountById(ctx, id)
    if err != nil {
      respondWithError(w, http.StatusBadRequest, "Record not found with id: " + id)
    } else {
      respondWithJSON(w, http.StatusOK, acc)
    }
  }
}

func createAccount(db account.Repo) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    var acc account.Account
    err := json.NewDecoder(r.Body).Decode(&acc)
    if err != nil {
      respondWithError(w, http.StatusInternalServerError, "Could not parse JSON body")
      return
    }
    acc, err = db.CreateAccount(ctx, acc)
    if err != nil {
      respondWithError(w, http.StatusInternalServerError, "Could not create account")
      return
    }
    respondWithJSON(w, http.StatusOK, acc)
  }
}

func createTransaction(db account.Repo) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := mux.Vars(r)["id"]
    acc, err := db.GetAccountById(ctx, id)
    if err != nil {
      respondWithError(w, http.StatusBadRequest, "Account not found with id: " + id)
      return
    }
    var t account.Transaction
    err = json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
      respondWithError(w, http.StatusInternalServerError, "Could not parse JSON body" + err.Error())
      return
    }
    t, err = db.CreateTransaction(ctx, acc, t)
    if err != nil {
      respondWithError(w, http.StatusInternalServerError, "Could not create transaction")
      return
    }
    respondWithJSON(w, http.StatusOK, t)
  }
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
