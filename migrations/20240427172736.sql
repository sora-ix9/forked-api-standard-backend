-- Create "roles" table
CREATE TABLE "roles" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "roles" ("deleted_at");
-- Create index "idx_roles_name" to table: "roles"
CREATE UNIQUE INDEX "idx_roles_name" ON "roles" ("name");
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "username" text NOT NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "role_id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
