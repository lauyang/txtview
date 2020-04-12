CREATE TABLE "txtview" (
  "id" integer,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "del" boolean NOT NULL DEFAULT false,
  "unlocktime" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);