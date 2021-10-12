# Simple Bank
## Simple bank service

The service that we’re going to build is a simple bank. It will provide APIs for the frontend to do following things:

1. Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
2. Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
3. Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.

## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup infrastructure

- Create the bank-network

    ``` bash
    make network
    ```

- Start postgres container:

    ```bash
    make postgres
    ```

- Create simple_bank database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### How to generate code

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    migrate create -ext sql -dir db/migration -seq <migration_name>
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run test:

    ```bash
    make test
    ```

## Deploy to kubernetes cluster

- [Install nginx ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/#aws):

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/aws/deploy.yaml
    ```

- [Install cert-manager](https://cert-manager.io/docs/installation/kubernetes/):

    ```bash
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
    ```

---
## Progress
### 1. Setup local environment

### 2. Design [dbdiagram](./db/dbdiagram) with https://dbdiagram.io/
- Foreign Key: `ref: > A.id`, 
- Timestamp Type: `timestamptz`
- Generate sql [000001_init_schema.up.sql](./db/migtation/000001_init_schema.up.sql)

### 3. Setup Postgres with Docker and DB Migration
```
make network
make postgres
make createdb
make migrateup
make dockerexecpostgres
\c simple_bank \dt
```
- now we should be able to see tables created by migration script
- we can also connect DB with [TablePlus](https://tableplus.com/)

### 4. Generate CRUD Golang code from SQL
- Write CRUD SQL query in [db/query](./db/query)
- generate golang code with `make sqlc`
- init go module `go mod init github.com/hhow09/simple_bank`

### 5. Write Golang unit tests for database CRUD with random data
- Write tests
    - [main_test.go](./db/sqlc/main_test.go): to make db connection
    - use `testQueries` to access functions in `[query].sql.go`
    - write following tests
    - [account_test.go](./db/sqlc/account_test.go)
    - [entry_test.go](./db/sqlc/entry_test.go)
    - [transfer_test.go](./db/sqlc/transfer_test.go)
- `make test`
- go [context](https://pkg.go.dev/context): carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.

### 6. implement database transaction in Golang
- Create [store.go](./db/sqlc/store.go)
    - `Store`: provides all funcs to execute queries and transactions
    - `execTx`: define a private transaction function: begin -> Commit or Rollback
    - `TransferTx`: define a public transfer transaction function
        1. create transfer record
	    2. create Entry of from_account
	    3. create Entry of to_account
        4. update accounts' balance
- Write [store_test.go](./db/sqlc/store_test.go)
    - create 5 goroutine to test transaction
    - get the err and result with go [channel](https://tour.golang.org/concurrency/2)
### 7. Handle Transaction Lock
- Now transfer transaction will not pass the test since
    - `GetAccount` SQL is `SELECT` and does not block each other
    - it will result in all concurrent `GetAccount` just return initial value
- Create a SQL that `SELECT FOR UPDATE`
    ```
    -- name: GetAccountForUpdate :one
    SELECT * FROM accounts
    WHERE id = $1 LIMIT 1;
    FOR UPDATE
    ```
- Now we will encounter `pq: deadlock detected`
    - simulate the QUERY BEING EXECUTED TO IDENTIFY LOCK
    - [PSQL - Lock Monitoring](https://wiki.postgresql.org/wiki/Lock_Monitoring)
    - deadlock are created between 2 transactions of `INSERT INTO transfers` and `SELECT * FROM accounts FOR UPDATE` 
- Handle Deadlock
    - since `transfers` Table has foreign key `from_account_id` and `to_account_id` referencing `accounts` Table
    - `INSERT INTO transfers` will acquire a `ExclusiveLock` on `accounts` Table to ensure that `ID` of accounts are not consistent.
    - However we are only update the `Balance` of account. The lock is unneeded.
    - change: `FOR UPDATE` -> `FOR NO KEY UPDATE`
- Refactor 
    - `getAccountForUpdate`+`UpdateAccount` = `AddAccountBalance`


### 8. Avoid DeadLock
- We will encounter deadlock when 2 transactions: `acc1` -> `acc2` and `acc2` -> `acc1` are running concurrently.
```sql
-- gorutine 1: transfer from id=1 to id=2
BEGIN;
UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *; 
COMMIT; 

-- gorutine 2: transfer from id=2 to id=1
BEGIN;
UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *; 
UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
COMMIT; 
```

- However if we switch the order so that **transactions always acquire locks in a consistent order**
```golang
if arg.FromAccountID < arg.ToAccountID {
	result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
} else {
	result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
}
```

the deadlock will not happen.
- we can test with `TestTransferTxDeadlock`

### 9. Isolation levels & read phenomena in MySQL & PostgreSQL 

#### MySQL
```sql
SELECT @@transaction_isolation --isolation level of current session
SELECT @@global.transaction_isolation --isolation level of global session
```

#### PostgreSQL
- only has 3 isolation level since `read uncommitted` is actually `read committed`
- transaction isolation can only be set in one transaction.
```sql
show transaction isolation level;
```
#### repeatable read is different in MySQL and PostgreSQL
- in isolation level `repeatable read`

```
Table accounts
 id  | owner  | balance | currency |          
-----+--------+---------+----------+
   1 | tom    |     100 | USD      |
   2 | mary   |     100 | USD      |
```

```
Steps:
1. process A, select * from accounts: tom's balance = 100
2. process B, select * from accounts where id=1: tom's balance = 100
3. process A, update accounts set balance = balance - 10 where id=1 returning *; tom's balance = 90
4. process B, select * from accounts where balance >=100: tom will appear since tom's balance = 100, (repeatable read)
5. process A, commit; 
6. process B, update accounts set balance = balance - 10 where id=1 returning *; 
7. process B, commit;
```

- running those steps in MySQL:
    - after step 7 we will get tom's balance: 80
    - however it does not make sense since we expect tom's balance to be 100
- running those steps in MySQL:
    - after step 6 we will get `ERROR:  could not serialize access due to concurrent update`

#### MySQL v.s. PostgreSQL
|                         | MySQL         | PostgreSQL           |
|-------------------------|---------------|----------------------|
| isolation levels        | 4             | 3                    |
| mechanism               | lock          | dependency detection |
| default isolation level | read commited | repeatable read      |

#### Reference
- [MySQL - 15.7.2.1 Transaction Isolation Levels](https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html)
- [PostgreSQL - 13.2. Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)

### 10. Github Action For unit test
- [.github/workflows/test.yml](./.github/workflows/test.yml)

### 11. Implement RESTful HTTP API in Go using Gin
- install [gin-gonic/gin](https://github.com/gin-gonic/gin) for http router and implement REST API
- in [api](./api)
    - define routing and handlers
    - use [Model binding and validation](https://github.com/gin-gonic/gin#model-binding-and-validation) to do param validation.

### 12. Load config from file & environment variables in Golang with Viper
- https://github.com/spf13/viper


### 13. Mock DB for testing HTTP API in Go and achieve 100% coverage
- We need a mockDB
#### Currently Server is using `store` of real db
```golang
type Server struct {
	store  *db.Store
    ...
}
```
- we can generate `Querier` interface by setting`emit_interface: true`
    - `var _ Querier = (*Queries)(nil)` means `Queries` struct must implement `Querier` interface
- we can make a higher level abstraction `Store` interface
    - we can embed `Querier` interface in to `Store` to ensure it has all 
    - `SQLStore` must implements `Store`
- `make mock`
    - package name: `mockdb`
    - destination `db/mock/store.go` 
    - Interface to mock: `Store`
- write tests in [account_test.go](./api/account_test.go) using `store *mockdb.MockStore`
    - use [httptest.NewRecorder](https://pkg.go.dev/net/http/httptest#NewRecorder) and `router.ServeHTTP` to write test and record them.
    - achieve 100%

### 14.  Implement transfer money API with a custom params validator
- implement `POST /transfers` API
    - instead of hard code the available currency in param binding, register a custom currency validator to gin
    - use `ShouldBindJSON` in `*gin.Context` to parse JSON
    - validate requested transfer currency should match two accounts

### 15. Add users table with unique & foreign key constraints in PostgreSQL
- now we still missing user authentication and authorization, we need user Tables.
- modify [dbdiagram](./db/dbdiagram) to add user Table, some db constraints:
    - email should be `unique`
    - a user can have multiple account: `accounts.owner` should be the foreign key of username.
    - an owner cannot own two account with same currency: 
        - `(owner, currency) [unique]`
        - `ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");`
- generate postgreSQL script from https://dbdiagram.io/
- `migrate create -ext sql -dir db/migration -seq add_users` to add a new migration script, [000002_add_users.up.sql](./db/migration/000002_add_users.up.sql)
- `make migratedown`

### 16. How to handle DB errors in Golang correctly
- `make sqlc` to generate updated model and `*.sql.go` 
- write unit test [user_test.go](./db/sqlc/user_test.go)
- handle db constraints violation in `CreateAccount`
    - `foreign_key_violation`: an account should only be created **when owner is an existing user**
    - `unique_violation`: an `owner-currency` index should be unique
    - both of them should return `http.StatusForbidden` (403)