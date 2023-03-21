-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency,
  acc_type
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 AND acc_type = 'bank'
LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 
 AND acc_type = 'bank'
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetExtAccount :one
SELECT * FROM accounts
WHERE owner = $1 AND currency = $2 AND acc_type = 'external'
LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = $1 
AND acc_type = 'bank'
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance+ sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;