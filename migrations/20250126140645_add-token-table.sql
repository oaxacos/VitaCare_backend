-- migrate:up
CREATE TABLE "tokens" (
  "id" uuid PRIMARY KEY,
  "token" text NOT NULL UNIQUE,
  "user_id" uuid,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "expired_at" timestamptz
);

ALTER TABLE IF EXISTS "tokens" ADD CONSTRAINT "fk_user_token_id" 
FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- migrate:down

DROP TABLE "tokens";
ALTER TABLE IF EXISTS "tokens" DROP CONSTRAINT "fk_user_token_id";
