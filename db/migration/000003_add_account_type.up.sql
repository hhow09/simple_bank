CREATE TYPE ACCOUNT_TYPE AS ENUM ('bank', 'external', 'credit_card');

ALTER TABLE "accounts" DROP CONSTRAINT "owner_currency_key";

ALTER TABLE "accounts" ADD acc_type ACCOUNT_TYPE NOT NULL DEFAULT('bank');

ALTER TABLE "accounts" ADD CONSTRAINT "uq_owner_currency_type" UNIQUE ("owner", "currency", "acc_type");