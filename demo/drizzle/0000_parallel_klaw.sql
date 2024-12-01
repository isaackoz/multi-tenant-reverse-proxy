CREATE TABLE IF NOT EXISTS "tenants" (
	"id" integer PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"hostname" text NOT NULL,
	CONSTRAINT "tenants_name_unique" UNIQUE("name"),
	CONSTRAINT "tenants_hostname_unique" UNIQUE("hostname")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "user" (
	"id" serial PRIMARY KEY NOT NULL,
	"age" integer
);
