ALTER TABLE "accounts" DROP CONSTRAINT "uq_owner_currency_type";

ALTER TABLE "accounts" DROP COLUMN acc_type

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

DROP TYPE ACCOUNT_TYPE;