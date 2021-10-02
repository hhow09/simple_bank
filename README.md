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
        4. TODO update from_account
- Write [store_test.go](./db/sqlc/store_test.go)
    - create 5 goroutine to test transaction
    - get the err and result with go [channel](https://tour.golang.org/concurrency/2)
    